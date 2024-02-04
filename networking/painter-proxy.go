package networking

import (
	"bytes"
	"net"
	. "playful-patterns.com/bakoko/world"
)

// This is an object that represents an entity that wants to send some
// graphics to be displayed (a painter).
// If someone wants to talk to the painter, they talk to this object
// and this object passes on information to the painter.
// The communication with the painter so far is this:
// - give me what you want to paint
// This is meant to be used by the gui which draws graphics sent by the world
// and the AI.
// This is a server that waits for a painter to connect to it.
type PainterProxy struct {
	Endpoint string
	conn     net.Conn
}

func (p *PainterProxy) GetPaintData() (info DebugInfo) {
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
		info.Deserialize(bytes.NewBuffer(data))
		return
	}
}
