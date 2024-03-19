package networking

import (
	"bytes"
	"errors"
	"log"
	"net"
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/ints"
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

type WorldRunner struct {
	w             World
	frameIdx      int
	watcher       FolderWatcher
	RecordingFile string
	currentInputs []PlayerInput
	initialized   bool
	player1       PlayerProxy
	player2       PlayerProxy
}

func (wr *WorldRunner) Initialize(player1 PlayerProxy, player2 PlayerProxy) {
	wr.frameIdx = 0
	wr.watcher.Folder = "world-data"
	wr.player1 = player1
	wr.player2 = player2
	wr.player1.SendWorld(&wr.w) // Should not block.
	wr.player2.SendWorld(&wr.w) // Should not block.
}

func (wr *WorldRunner) Step() {
	if wr.watcher.FolderContentsChanged() {
		LoadWorld(&wr.w)
	}

	var input Input
	input.Player1Input = wr.player1.GetInput() // Should block.
	input.Player2Input = wr.player2.GetInput() // Should block.

	if wr.RecordingFile != "" {
		wr.currentInputs = append(wr.currentInputs, input.Player1Input)
		SerializeInputs(wr.currentInputs, wr.RecordingFile)
	}

	if input.Player1Input.Reload || input.Player2Input.Reload {
		LoadWorld(&wr.w)
	}

	if !input.Player1Input.Pause && !input.Player2Input.Pause {
		wr.w.Step(&input, wr.frameIdx)
	}

	//guiProxy.SendPaintData(&world.DebugInfo) // Should not block.
	wr.player1.SendWorld(&wr.w) // Should not block.
	wr.player2.SendWorld(&wr.w) // Should not block.

	//if input.Player1Input.Quit || input.Player2Input.Quit {
	//	break
	//}
	wr.frameIdx++
	wr.w.JustReloaded = ZERO
}

type AiRunner struct {
	ai         PlayerAI
	worldProxy WorldProxy
}

func (ar *AiRunner) Initialize(worldProxy WorldProxy) {
	ar.worldProxy = worldProxy
	ar.ai.PauseBetweenShots = 1500 * time.Millisecond
	ar.ai.LastShot = time.Now()
}

func (ar *AiRunner) GetWorld() *World {
	// This should block as the AI doesn't make sense if it doesn't
	// synchronize with the simulation.
	for {
		if err := ar.worldProxy.Connect(); err != nil {
			continue // Retry from the beginning.
		}
		var err error
		var w *World
		if w, err = ar.worldProxy.GetWorld(); err != nil {
			continue // Retry from the beginning.
		}
		return w
	}
}

func (ar *AiRunner) SendInput(input *PlayerInput) {
	// This should block as the AI doesn't make sense if it doesn't
	// synchronize with the simulation.
	for {
		if err := ar.worldProxy.Connect(); err != nil {
			continue // Retry from the beginning.
		}
		if err := ar.worldProxy.SendInput(input); err != nil {
			continue // Retry from the beginning.
		}
		break
	}
}

func (ar *AiRunner) Step() {
	w := ar.GetWorld()
	input := ar.ai.Step(w)
	ar.SendInput(&input)

	// This may or may not block, who cares?
	//guiProxy.SendPaintData(&ai.DebugInfo)
}

// Regular
type WorldPlayerProxy struct {
	input *PlayerInput
	world *World
}

func (p *WorldPlayerProxy) Connect() error {
	return nil
}

func (p *WorldPlayerProxy) SendInput(input *PlayerInput) error {
	p.input = input
	return nil
}

func (p *WorldPlayerProxy) GetWorld() (world *World, err error) {
	return p.world, nil
}

func (p *WorldPlayerProxy) GetInput() PlayerInput {
	return *p.input
}

func (p *WorldPlayerProxy) SendWorld(world *World) {
	p.world = world
}

type WorldProxyPlayback struct {
	worldRunner  WorldRunner
	aiRunner     AiRunner
	worldAiProxy WorldPlayerProxy
	initialized  bool
	player1Input PlayerInput
	world        *World
}

func (p *WorldProxyPlayback) Connect() error {
	return nil
}

func (p *WorldProxyPlayback) SendInput(player1Input *PlayerInput) error {
	if !p.initialized {
		p.aiRunner.Initialize(&p.worldAiProxy)
		p.worldRunner.Initialize(p, &p.worldAiProxy)
		p.initialized = true
	}

	p.player1Input = *player1Input
	p.aiRunner.Step()
	p.worldRunner.Step()
	return nil
}

func (p *WorldProxyPlayback) GetInput() PlayerInput {
	return p.player1Input
}

func (p *WorldProxyPlayback) SendWorld(world *World) {
	p.world = world
}

func (p *WorldProxyPlayback) GetWorld() (w *World, err error) {
	return p.world, nil
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

	if p.Timeout.Milliseconds() > 0 {
		err = conn.SetDeadline(time.Now().Add(p.Timeout))
		if err != nil {
			return err
		}
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

	err := WriteData(p.conn, buf.Bytes())
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and move on.
	if err != nil {
		log.Println(err)
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

	data, err := ReadData(p.conn)
	// If there was an error, assume the peer is no longer available.
	// Invalidate the connection and try again later.
	if err != nil {
		p.conn = nil
		log.Println("lost connection (3)")
		return nil, err
	}

	w = &World{}
	w.Deserialize(bytes.NewBuffer(data))
	return w, nil
}
