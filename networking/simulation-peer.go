package networking

import (
	"bytes"
	"log"
	"net"
	. "playful-patterns.com/bakoko/world"
	"time"
)

type SimulationPeer struct {
	Endpoint string
	conn     net.Conn
}

// Doesn't matter if this fails.
func (p *SimulationPeer) GetWorld(w *World) {
	// Don't do anything if we don't have a peer.
	// The communication between us and the peer is always that:
	// - we connect to the peer
	// - we send input to the peer
	// - we get an ouput from the peer
	// If the peer disconnects in middle of that, we start from the beginning,
	// we don't accept a connection then continue with getting the output.
	if p.conn == nil {
		return
	}

	data, err := ReadData(p.conn)
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn = nil
		log.Println("lost connection")
		return
	}

	w.Deserialize(bytes.NewBuffer(data))
}

// Try to send an input to the peer, but don't block.
func (p *SimulationPeer) SendInput(input *PlayerInput) {
	// If we don't have a peer, connect to one.
	if p.conn == nil {
		var err error
		p.conn, err = net.DialTimeout("tcp", p.Endpoint, 5*time.Millisecond)

		// If connection took too long or failed, screw it.
		// We'll try again later.
		if err != nil {
			//log.Println("could not connect!")
			return
		}
	}
	//log.Println("connection established!")

	// We have a connection, try to send our input.
	buf := new(bytes.Buffer)
	Serialize(buf, input)

	err := WriteData(p.conn, buf.Bytes())
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn = nil
		log.Println("lost connection")
	}
}
