package transport

// ServerTransport is held by the simulation goroutine.
// It reads player commands and sends world state snapshots.
type ServerTransport interface {
	// ReceiveCommands drains all pending commands from clients.
	// Returns an empty slice (not nil) when no commands are queued.
	// Non-blocking.
	ReceiveCommands() []*Command

	// SendSnapshot delivers a snapshot of authoritative world state to clients.
	// May drop the snapshot if the client buffer is full (jitter tolerance).
	SendSnapshot(snapshot *Snapshot)

	// Close shuts down the transport. Safe to call multiple times.
	Close()
}

// ClientTransport is held by the render/Ebitengine side.
// It sends player input commands and polls for the latest world snapshot.
type ClientTransport interface {
	// SendCommand queues a command for the server.
	// Non-blocking; drops the command if the server buffer is full.
	SendCommand(cmd *Command)

	// ReceiveSnapshot returns the most recent snapshot from the server,
	// or nil if no new snapshot has arrived since the last call.
	// Non-blocking.
	ReceiveSnapshot() *Snapshot

	// Close shuts down the transport. Safe to call multiple times.
	Close()
}
