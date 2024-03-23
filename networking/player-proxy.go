package networking

import (
	"bytes"
	"net"
	. "playful-patterns.com/bakoko/world"
)

type PlayerProxy interface {
	// Should block until both have happened for the same player:
	// - the world has been sent successfully
	// - the reaction has been received successfully
	SendWorldGetInput(w *World) *PlayerInput
}

// This is an object that represents a PlayerProxy.
// If someone wants to talk to the PlayerProxy, they talk to this object
// and this object passes on information to the PlayerProxy.
// The communication with the PlayerProxy so far is this:
// - give me an input
// - here's the world
// This is meant to be used by the world which talks to to players.
// This is a server that waits for a PlayerProxy to connect to it.
type PlayerProxyTcpIp struct {
	Endpoint string
	conn     net.Conn
}

// We want to:
// - connect to the player if we are not connected
// - send the world to the player
// - and get the reaction from the player
// If the player disconnects, we must start the process over from the beginning.
// We don't send the world to one connection and get the reaction from another
// connection.
func (p *PlayerProxyTcpIp) SendWorldGetInput(w *World) *PlayerInput {
	// Keep trying to perform the transaction.
	for {
		// If we don't have a peer, wait until we get one.
		if p.conn == nil {
			// Listen for incoming connections
			listener, err := net.Listen("tcp", p.Endpoint)
			Check(err)

			// Accept one incoming connection.
			p.conn, err = listener.Accept()
			Check(err)

			listener.Close()
		}

		// Try sending the world to our peer.
		data := w.Serialize()
		if err := WriteData(p.conn, data, 0); err != nil {
			// There was an error. Nevermind, close the connection and wait
			// for a new one.
			p.conn.Close()
			p.conn = nil
			continue // Wait for a peer again.
		}

		// Try to get data from our peer.
		data, err := ReadData(p.conn, 0)
		if err != nil {
			// There was an error. Nevermind, close the connection and wait
			// for a new one.
			p.conn.Close()
			p.conn = nil
			continue // Wait for a peer again.
		}

		// Finally, we can return the input.
		var input PlayerInput
		Deserialize(bytes.NewBuffer(data), &input)
		return &input
	}
}
