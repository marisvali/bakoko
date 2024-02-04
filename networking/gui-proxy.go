package networking

import (
	"log"
	"net"
	. "playful-patterns.com/bakoko/world"
	"time"
)

// This is an object that represents the gui, for a module that wants to draw
// some graphics on it.
// If someone wants to talk to the gui, they talk to this object
// and this object passes on information to the gui.
// The communication with the gui so far is this:
// - here are some graphics to draw
// This is meant to be used by the world and the AI module which acts in the
// world.
// This is a client that connects to a server.
type GuiProxy struct {
	Endpoint string
	conn     net.Conn
}

// Try to send an input to the peer, but don't block.
func (p *GuiProxy) SendPaintData(debugInfo *DebugInfo) {
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
	data := debugInfo.Serialize()

	err := WriteData(p.conn, data)
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn = nil
		log.Println("lost connection (2)")
	}
}
