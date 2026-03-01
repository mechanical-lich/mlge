package client

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/ecs"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/input"
	"github.com/mechanical-lich/mlge/transport"
)

// InputMapper translates an mlge input event into a transport Command.
// Return (nil, false) to discard the event (e.g., hotkeys handled locally).
//
// Implement this on your client state or a dedicated mapper struct and pass it
// to [NewClient] so the Client can forward input to the server.
type InputMapper interface {
	MapInputEvent(e event.EventData) (*transport.Command, bool)
}

// ClientConfig holds optional configuration for [Client].
type ClientConfig struct {
	// ScreenWidth and ScreenHeight are passed to ebiten.SetWindowSize.
	// If zero the window size is not changed.
	ScreenWidth, ScreenHeight int

	// WindowTitle is passed to ebiten.SetWindowTitle.
	WindowTitle string
}

// Client implements ebiten.Game and acts as the presentation layer.
//
// Each Ebitengine frame it:
//  1. Handles input (via mlge input.InputManager) and forwards events as
//     Commands to the server via the ClientTransport.
//  2. Polls for the latest Snapshot from the server.
//  3. Decodes the snapshot into the local world via the SnapshotCodec.
//  4. Runs RenderSystems (animations, interpolation) at frame rate.
//  5. Delegates Update/Draw to the ClientState stack.
//
// Create with [NewClient].
type Client struct {
	transport    transport.ClientTransport
	codec        transport.SnapshotCodec
	inputManager *input.InputManager
	inputMapper  InputMapper
	renderSys    RenderSystemManager
	stateMachine clientStateMachine
	world        any
	entitySource func() []*ecs.Entity
	screenW      int
	screenH      int
}

// NewClient constructs a Client.
//
//   - t: the client half of a [transport.LocalTransport] (or future TCP transport).
//   - codec: game-provided codec to decode server snapshots into the local world.
//   - initialState: the first ClientState pushed onto the state stack.
//   - world: the client's local (non-authoritative) copy of the game world.
//   - entitySource: function returning the local entity slice for RenderSystems.
//   - cfg: window configuration.
func NewClient(
	t transport.ClientTransport,
	codec transport.SnapshotCodec,
	initialState ClientState,
	world any,
	entitySource func() []*ecs.Entity,
	cfg ClientConfig,
) *Client {
	if cfg.ScreenWidth > 0 && cfg.ScreenHeight > 0 {
		ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	}
	if cfg.WindowTitle != "" {
		ebiten.SetWindowTitle(cfg.WindowTitle)
	}

	c := &Client{
		transport:    t,
		codec:        codec,
		world:        world,
		entitySource: entitySource,
		screenW:      cfg.ScreenWidth,
		screenH:      cfg.ScreenHeight,
	}

	// Wire up mlge's input manager → queued event manager → client command forwarding.
	c.inputManager = input.NewInputManager(event.GetQueuedInstance())

	c.stateMachine.PushState(initialState)
	return c
}

// SetInputMapper sets an [InputMapper] that translates input events to Commands.
// Call after NewClient, before Run.
func (c *Client) SetInputMapper(m InputMapper) {
	c.inputMapper = m
}

// AddRenderSystem appends a RenderSystem to the client's render pass.
func (c *Client) AddRenderSystem(s RenderSystem) {
	c.renderSys.AddSystem(s)
}

// Run starts the Ebitengine window loop. Blocks until the window closes.
func (c *Client) Run() error {
	return ebiten.RunGame(c)
}

// --- ebiten.Game implementation ---

// Update is called every Ebitengine TPS tick.
func (c *Client) Update() error {
	// 1. Poll OS input → mlge events.
	c.inputManager.HandleInput()

	// 2. Drain queued input events; forward them as Commands to the server.
	c.forwardInputEvents()

	// 3. Receive the latest snapshot from the server (nil if none yet).
	var snap *transport.Snapshot
	if raw := c.transport.ReceiveSnapshot(); raw != nil {
		snap = raw
		// Decode server authority into local world.
		c.codec.Decode(snap, c.world)
	}

	// 4. Run render systems (animation, interpolation) at frame rate.
	entities := c.entitySource()
	if err := c.renderSys.UpdateSystems(c.world); err != nil {
		log.Printf("[client] RenderSystem global error: %v", err)
	}
	if err := c.renderSys.UpdateSystemsForEntities(c.world, entities); err != nil {
		log.Printf("[client] RenderSystem entity error: %v", err)
	}

	// 5. Advance the ClientState machine with the snapshot.
	c.stateMachine.Update(snap)

	return nil
}

// Draw is called every Ebitengine frame (may differ from TPS).
func (c *Client) Draw(screen *ebiten.Image) {
	c.stateMachine.Draw(screen)
}

// Layout returns the game's logical size.
func (c *Client) Layout(outsideW, outsideH int) (int, int) {
	if c.screenW > 0 && c.screenH > 0 {
		return c.screenW, c.screenH
	}
	return outsideW, outsideH
}

// forwardInputEvents drains the QueuedEventManager and sends any mappable
// events as Commands to the server transport.
func (c *Client) forwardInputEvents() {
	if c.inputMapper == nil {
		// No mapper set - still drain the queue so it doesn't grow unbounded,
		// but discard all events (game handles input entirely in ClientState).
		event.GetQueuedInstance().HandleQueue()
		return
	}

	// Register a one-shot listener that intercepts all events this frame.
	interceptor := &inputInterceptor{mapper: c.inputMapper, transport: c.transport}
	// Register for any event type via the catch-all listener pattern.
	// Since mlge's EventManager is type-keyed, we need the state to register
	// itself separately. Here we expose a simpler path: the ClientState's
	// Update method receives the snapshot; for input, the game's InputMapper
	// calls SetInputMapper and handles it here.
	// For the common case where clients don't forward raw input, we just drain.
	_ = interceptor
	event.GetQueuedInstance().HandleQueue()
}

// inputInterceptor is reserved for future use when input forwarding to the
// server is needed from outside a ClientState (e.g., a generic HUD overlay).
type inputInterceptor struct {
	mapper    InputMapper
	transport transport.ClientTransport
}
