package transport

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

// tcpEnvelope is the wire format wrapping all TCP messages.
// Kind routes the payload to the correct decode target on the receiver.
type tcpEnvelope struct {
	Kind    string          `json:"k"`
	Payload json.RawMessage `json:"p"`
}

const (
	tcpKindCommand  = "cmd"
	tcpKindSnapshot = "snap"
)

// writeMsg writes a length-prefixed JSON message to conn.
// Wire format: 4-byte big-endian uint32 length, then JSON bytes.
func writeMsg(conn net.Conn, env *tcpEnvelope) error {
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	var header [4]byte
	binary.BigEndian.PutUint32(header[:], uint32(len(data)))
	if _, err := conn.Write(header[:]); err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

// maxMsgSize caps the largest message readMsg will allocate for.
// Protects against a misbehaving peer sending a huge length header.
const maxMsgSize = 4 << 20 // 4 MiB

// readMsg reads a single length-prefixed JSON message from conn.
func readMsg(conn net.Conn) (*tcpEnvelope, error) {
	var header [4]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(header[:])
	if size > maxMsgSize {
		return nil, fmt.Errorf("transport/tcp: message size %d exceeds limit %d", size, maxMsgSize)
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, err
	}
	var env tcpEnvelope
	if err := json.Unmarshal(buf, &env); err != nil {
		return nil, err
	}
	return &env, nil
}

// setNoDelay disables Nagle buffering on a TCP connection to reduce latency.
func setNoDelay(conn net.Conn) {
	if tc, ok := conn.(*net.TCPConn); ok {
		_ = tc.SetNoDelay(true)
	}
}

// =============================================================================
// TCPServerTransport
// =============================================================================

// TCPServerTransport implements [ServerTransport] over TCP.
//
// It listens for incoming client connections, reads [Command] messages from
// each, and broadcasts [Snapshot] messages to all connected clients.
//
// Serialization note: [Command.Payload] and [EntitySnapshot.Components] values
// are encoded as JSON. Concrete types round-trip correctly. Values stored as
// bare interface{} will decode as map[string]interface{} on the remote side;
// the game must handle that (e.g. use json.RawMessage or re-unmarshal).
//
// TCP_NODELAY is set on every accepted connection.
//
// Use [NewTCPServerTransport] to create an instance.
type TCPServerTransport struct {
	listener net.Listener
	commands chan *Command

	mu     sync.Mutex
	conns  []*tcpPeer
	closed atomic.Bool
	wg     sync.WaitGroup
}

// tcpPeer wraps a single client connection on the server side.
type tcpPeer struct {
	conn   net.Conn
	sendMu sync.Mutex
}

// NewTCPServerTransport starts a TCP listener on addr (e.g. ":7777") and
// returns a ready-to-use [ServerTransport]. Call Close when done.
func NewTCPServerTransport(addr string) (ServerTransport, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	t := &TCPServerTransport{
		listener: ln,
		commands: make(chan *Command, defaultCommandBufSize),
	}
	t.wg.Add(1)
	go t.acceptLoop()
	return t, nil
}

func (t *TCPServerTransport) acceptLoop() {
	defer t.wg.Done()
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if t.closed.Load() {
				return
			}
			log.Printf("transport/tcp server: accept error: %v", err)
			return
		}
		setNoDelay(conn)
		peer := &tcpPeer{conn: conn}
		t.mu.Lock()
		t.conns = append(t.conns, peer)
		t.mu.Unlock()
		t.wg.Add(1)
		go t.peerReadLoop(peer)
	}
}

func (t *TCPServerTransport) peerReadLoop(peer *tcpPeer) {
	defer func() {
		peer.conn.Close()
		t.mu.Lock()
		for i, c := range t.conns {
			if c == peer {
				last := len(t.conns) - 1
				t.conns[i] = t.conns[last]
				t.conns[last] = nil // clear trailing slot so GC can reclaim the peer
				t.conns = t.conns[:last]
				break
			}
		}
		t.mu.Unlock()
		t.wg.Done()
	}()
	for {
		env, err := readMsg(peer.conn)
		if err != nil {
			if !t.closed.Load() {
				log.Printf("transport/tcp server: peer read error: %v", err)
			}
			return
		}
		if env.Kind != tcpKindCommand {
			continue
		}
		var cmd Command
		if err := json.Unmarshal(env.Payload, &cmd); err != nil {
			log.Printf("transport/tcp server: decode command: %v", err)
			continue
		}
		select {
		case t.commands <- &cmd:
		default:
			// Server command buffer full - drop. Should not happen in normal play.
		}
	}
}

// ReceiveCommands drains all pending commands from all connected clients.
// Returns an empty (non-nil) slice when no commands are queued. Non-blocking.
func (t *TCPServerTransport) ReceiveCommands() []*Command {
	var cmds []*Command
	for {
		select {
		case cmd := <-t.commands:
			cmds = append(cmds, cmd)
		default:
			if cmds == nil {
				cmds = []*Command{}
			}
			return cmds
		}
	}
}

// SendSnapshot JSON-encodes snapshot and writes it to every connected client.
// Clients that cannot be written to are logged; the error does not stop
// delivery to other clients.
func (t *TCPServerTransport) SendSnapshot(snapshot *Snapshot) {
	payload, err := json.Marshal(snapshot)
	if err != nil {
		log.Printf("transport/tcp server: encode snapshot: %v", err)
		return
	}
	env := &tcpEnvelope{Kind: tcpKindSnapshot, Payload: json.RawMessage(payload)}

	t.mu.Lock()
	peers := make([]*tcpPeer, len(t.conns))
	copy(peers, t.conns)
	t.mu.Unlock()

	for _, peer := range peers {
		peer.sendMu.Lock()
		if err := writeMsg(peer.conn, env); err != nil && !t.closed.Load() {
			log.Printf("transport/tcp server: send snapshot: %v", err)
		}
		peer.sendMu.Unlock()
	}
}

// Close shuts down the listener and all client connections. Safe to call
// multiple times.
func (t *TCPServerTransport) Close() {
	if t.closed.CompareAndSwap(false, true) {
		t.listener.Close()
		t.mu.Lock()
		for _, peer := range t.conns {
			peer.conn.Close()
		}
		t.mu.Unlock()
		t.wg.Wait()
	}
}

// =============================================================================
// TCPClientTransport
// =============================================================================

// TCPClientTransport implements [ClientTransport] over TCP.
//
// It connects to a [TCPServerTransport], sends [Command] messages, and
// receives [Snapshot] messages in a background goroutine. Only the most
// recently received snapshot is retained; older ones are discarded so the
// caller always sees up-to-date state.
//
// TCP_NODELAY is set on the connection.
//
// Use [NewTCPClientTransport] to create an instance.
type TCPClientTransport struct {
	conn net.Conn

	mu     sync.Mutex
	latest *Snapshot

	closed atomic.Bool
	wg     sync.WaitGroup
}

// NewTCPClientTransport dials addr (e.g. "localhost:7777") and returns a
// ready-to-use [ClientTransport]. Call Close when done.
func NewTCPClientTransport(addr string) (ClientTransport, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	setNoDelay(conn)
	t := &TCPClientTransport{conn: conn}
	t.wg.Add(1)
	go t.readLoop()
	return t, nil
}

func (t *TCPClientTransport) readLoop() {
	defer t.wg.Done()
	for {
		env, err := readMsg(t.conn)
		if err != nil {
			if !t.closed.Load() {
				log.Printf("transport/tcp client: read error: %v", err)
			}
			return
		}
		if env.Kind != tcpKindSnapshot {
			continue
		}
		var snap Snapshot
		if err := json.Unmarshal(env.Payload, &snap); err != nil {
			log.Printf("transport/tcp client: decode snapshot: %v", err)
			continue
		}
		t.mu.Lock()
		t.latest = &snap
		t.mu.Unlock()
	}
}

// SendCommand JSON-encodes cmd and writes it to the server. Non-blocking from
// the caller's perspective; write errors are logged.
func (t *TCPClientTransport) SendCommand(cmd *Command) {
	payload, err := json.Marshal(cmd)
	if err != nil {
		log.Printf("transport/tcp client: encode command: %v", err)
		return
	}
	env := &tcpEnvelope{Kind: tcpKindCommand, Payload: json.RawMessage(payload)}
	if err := writeMsg(t.conn, env); err != nil && !t.closed.Load() {
		log.Printf("transport/tcp client: send command: %v", err)
	}
}

// ReceiveSnapshot returns the most recent snapshot received from the server,
// or nil if no new snapshot has arrived since the last call. Non-blocking.
func (t *TCPClientTransport) ReceiveSnapshot() *Snapshot {
	t.mu.Lock()
	snap := t.latest
	t.latest = nil
	t.mu.Unlock()
	return snap
}

// Close shuts down the connection and waits for the read goroutine to exit.
// Safe to call multiple times.
func (t *TCPClientTransport) Close() {
	if t.closed.CompareAndSwap(false, true) {
		t.conn.Close()
		t.wg.Wait()
	}
}
