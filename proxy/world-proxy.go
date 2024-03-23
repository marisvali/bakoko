package proxy

import (
	"bytes"
	"errors"
	"log"
	"net"
	. "playful-patterns.com/bakoko/world"
	"time"
)

// This is an object that represents the world, for a PlayerProxy module.
// If someone wants to talk to the world, they talk to this object
// and this object passes on information to the world.
// The communication with the world so far is this:
// - here's an input
// - give me the world
// This is meant to be used by PlayerProxy modules which act in the world.
// This is a client that connects to a server.
type WorldProxy interface {
	Connect() error
	SendInput(input *PlayerInput) error
	GetWorld() (w *World, err error)
}

// TCP IP
type WorldProxyTcpIp struct {
	Endpoint string
	Timeout  time.Duration
	conn     net.Conn
}

func (p *WorldProxyTcpIp) Connect() error {
	// If we already have a peer, move on.
	if p.conn != nil {
		return nil
	}

	// We don't have a peer, connect to one.
	conn, err := net.DialTimeout("tcp", p.Endpoint, p.Timeout)
	if err != nil {
		return err // Error, give up.
	}

	p.conn = conn
	return nil
}

// Try to send an input to the peer, but don't block.
func (p *WorldProxyTcpIp) SendInput(input *PlayerInput) error {
	if p.conn == nil {
		return errors.New("no connection")
	}

	// Try to send our input.
	buf := new(bytes.Buffer)
	Serialize(buf, input)

	err := WriteData(p.conn, buf.Bytes(), p.Timeout)
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and move on.
	if err != nil {
		log.Println(err)
		p.conn.Close()
		p.conn = nil
		log.Println("lost connection (1)")
		return err
	}
	return nil
}

// Try to get the world, but don't block if it doesn't work.
func (p *WorldProxyTcpIp) GetWorld() (w *World, err error) {
	if p.conn == nil {
		return nil, errors.New("no connection")
	}

	data, err := ReadData(p.conn, p.Timeout)
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn.Close()
		p.conn = nil
		log.Println("lost connection (3)")
		return nil, err
	}

	w = &World{}
	w.Deserialize(bytes.NewBuffer(data))
	return w, nil
}
