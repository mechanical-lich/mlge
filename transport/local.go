package transport

import "sync"

const (
	defaultCommandBufSize  = 64
	defaultSnapshotBufSize = 4 // client only needs the latest; small buffer avoids head-of-line blocking
)

// LocalTransport is an in-process, channel-backed implementation of
// [ServerTransport] and [ClientTransport].
//
// Use it for single-player or split-screen games where the server and client
// share a single executable. There is no serialization cost - commands and
// snapshots are passed as pointer values through buffered Go channels.
//
// Create an instance with [NewLocalTransport].
type LocalTransport struct {
	commands  chan *Command
	snapshots chan *Snapshot
	closeOnce sync.Once
}

// NewLocalTransport returns a joined ServerTransport/ClientTransport pair that
// communicate through shared in-process channels.
//
//	srvT, cliT := transport.NewLocalTransport()
//	go server.Run(srvT)
//	ebiten.RunGame(client.New(cliT))
func NewLocalTransport() (ServerTransport, ClientTransport) {
	lt := &LocalTransport{
		commands:  make(chan *Command, defaultCommandBufSize),
		snapshots: make(chan *Snapshot, defaultSnapshotBufSize),
	}
	return &localServerSide{lt}, &localClientSide{lt}
}

// --- server side ---

type localServerSide struct{ t *LocalTransport }

func (s *localServerSide) ReceiveCommands() []*Command {
	var cmds []*Command
	for {
		select {
		case cmd := <-s.t.commands:
			cmds = append(cmds, cmd)
		default:
			if cmds == nil {
				cmds = []*Command{} // always return non-nil
			}
			return cmds
		}
	}
}

func (s *localServerSide) SendSnapshot(snapshot *Snapshot) {
	select {
	case s.t.snapshots <- snapshot:
	default:
		// Buffer is full - drop oldest snapshot and try again.
		// This ensures the client always gets the most recent state, not a stale one.
		select {
		case <-s.t.snapshots:
		default:
		}
		select {
		case s.t.snapshots <- snapshot:
		default:
		}
	}
}

func (s *localServerSide) Close() {
	s.t.closeOnce.Do(func() {
		close(s.t.commands)
		close(s.t.snapshots)
	})
}

// --- client side ---

type localClientSide struct{ t *LocalTransport }

func (c *localClientSide) SendCommand(cmd *Command) {
	select {
	case c.t.commands <- cmd:
	default:
		// Server buffer full - drop command. Shouldn't happen at 60 TPS with
		// a 64-slot buffer unless the server is severely stalled.
	}
}

func (c *localClientSide) ReceiveSnapshot() *Snapshot {
	// Drain to latest: return the last available snapshot, discarding older ones.
	var latest *Snapshot
	for {
		select {
		case snap, ok := <-c.t.snapshots:
			if !ok {
				return latest // channel closed
			}
			latest = snap
		default:
			return latest
		}
	}
}

func (c *localClientSide) Close() {
	// Client close delegates to server side since they share the transport.
	// Calling Close from either side is safe.
	c.t.closeOnce.Do(func() {
		close(c.t.commands)
		close(c.t.snapshots)
	})
}
