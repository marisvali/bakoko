package networking

import (
	"bytes"
	"net"
	. "playful-patterns.com/bakoko/world"
)

type PlayerProxy interface {
	GetInput() PlayerInput
	SendWorld(w *World)
}

type PlayerProxyRegular struct {
	worldData  []byte
	WorldProxy *WorldProxyRegular
}

func (p *PlayerProxyRegular) GetInput() PlayerInput {
	// I have to block here until input is actually available.

	// Get the input from the world proxy.
	return p.WorldProxy.input
}

func (p *PlayerProxyRegular) SendWorld(w *World) {
	// Store the world.
	p.worldData = w.Serialize()
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

func (p *PlayerProxyTcpIp) GetInput() PlayerInput {
	// Keep trying to get an input from a peer.
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

		// Try to get data from our peer.
		data, err := ReadData(p.conn)
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
		return input
	}
}

// Doesn't matter if this fails.
func (p *PlayerProxyTcpIp) SendWorld(w *World) {
	// Don't do anything if we don't have a peer.
	// The communication between us and the peer is always that:
	// - the peer connects
	// - we get input from the peer
	// - we send an ouput to the peer
	// If the peer disconnects in middle of that, we start from the beginning,
	// we don't accept a connection then continue with sending the output.
	if p.conn == nil {
		return
	}

	data := w.Serialize()

	err := WriteData(p.conn, data)
	if err != nil {
		p.conn = nil
	}
}
