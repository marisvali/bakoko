package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	. "playful-patterns.com/bakoko"
	"time"
)
import . "playful-patterns.com/bakoko/ints"

var input Input
var w World

func init() {
	w.Player1 = Player{
		Bounds: Circle{
			Pos:      UPt(60, 260),
			Diameter: U(50)},
		NBalls:   I(3),
		BallType: I(1),
		Health:   I(3)}
	w.Player2 = Player{
		Bounds: Circle{
			Pos:      UPt(180, 60),
			Diameter: U(50)},
		NBalls:   I(3),
		BallType: I(2),
		Health:   I(3)}
	w.Balls = []Ball{
		{
			Bounds: Circle{
				Pos:      UPt(120, 70),
				Diameter: U(30)},
			Type: I(1),
		},
		{
			Bounds: Circle{
				Pos:      UPt(90, 210),
				Diameter: U(30)},
			Type: I(1),
		},
		{
			Bounds: Circle{
				Pos:      UPt(190, 140),
				Diameter: U(30)},
			Type: I(2),
		}}
	w.Obstacles.Init(I(10), I(10))
	w.ObstacleSize = U(30)
	for y := I(0); y.Lt(w.Obstacles.NCols()); y.Inc() {
		for x := I(0); x.Lt(w.Obstacles.NCols()); x.Inc() {
			if y.Plus(x).Mod(I(2)).Eq(I(0)) {
				w.Obstacles.Set(x, y, I(0))
			} else {
				w.Obstacles.Set(x, y, I(1))
			}
		}
	}
}
func main5() {
	originalWorld := w
	w.SerializeToFile("world.bin")
	inputs := DeserializeInputs("recorded-inputs")
	start := time.Now()
	bigIdx := 0
	for ; bigIdx < 1000000; bigIdx++ {
		frameIdx := 0
		w := originalWorld
		for ; !input.Player1Input.Quit; frameIdx++ {
			input = Input{}
			if frameIdx < len(inputs) {
				input = inputs[frameIdx]
			}

			w.Step(&input, frameIdx)
		}
	}
	fmt.Println(time.Since(start).Seconds())
	fmt.Printf("%.12f\n", time.Since(start).Seconds()/float64(bigIdx))
	w.SerializeToFile("world.bin")
}
func main4() {
	w.SerializeToFile("world.bin")
	inputs := DeserializeInputs("recorded-inputs")

	frameIdx := 0
	for ; !input.Player1Input.Quit; frameIdx++ {
		//frameStart := time.Now()
		WaitForFile("input-ready")
		//input.DeserializeFromFile("input.bin")
		input = Input{}
		if frameIdx < len(inputs) {
			input = inputs[frameIdx]
		}

		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")

		//frameDuration := time.Since(frameStart)
		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
	}
}
func main3() {
	w.SerializeToFile("world.bin")

	frameIdx := 0
	for ; !input.Player1Input.Quit; frameIdx++ {
		//frameStart := time.Now()
		WaitForFile("input-ready")
		//input.DeserializeFromFile("input.bin")
		input = Input{}

		input.Player1Input.MoveLeft = rand.Int()%6 == 0
		input.Player1Input.MoveRight = rand.Int()%5 == 0
		input.Player1Input.MoveUp = rand.Int()%6 == 0
		input.Player1Input.MoveDown = rand.Int()%6 == 0
		input.Player1Input.Quit = frameIdx == 3600

		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")

		//frameDuration := time.Since(frameStart)
		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
	}
}

var currentInputs []Input

func GetNewRecordingFile() string {
	date := time.Now()
	for i := 0; i < 1000000; i++ {
		filename := fmt.Sprintf("recordings/recorded-inputs-%04d-%02d-%02d-%06d",
			date.Year(), date.Month(), date.Day(), i)
		if !FileExists(filename) {
			return filename
		}
	}
	panic("Cannot record, no available filename found.")
}

type interfacePeer struct {
	endpoint string
	conn     net.Conn
}

func (p *interfacePeer) getInput() PlayerInput {
	// Keep trying to get an input from a peer.
	for {
		// If we don't have a peer, wait until we get one.
		if p.conn == nil {
			// Listen for incoming connections
			listener, err := net.Listen("tcp", p.endpoint)
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
func (p *interfacePeer) sendWorld(w *World) {
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

func main() {
	frameIdx := 0
	player1 := interfacePeer{}
	player1.endpoint = "localhost:56901"
	//player2 := interfacePeer{}
	//player2.endpoint = "localhost:56902"
	for w.Over.Eq(I(0)) {
		var input Input
		input.Player1Input = player1.getInput()
		//input.Player2Input = player2.getInput()
		w.Step(&input, frameIdx)
		player1.sendWorld(&w)
		//player2.sendWorld(&w)

		if input.Player1Input.Quit || input.Player2Input.Quit {
			break
		}
	}
}

func main77() {
	//if utils.FileExists("world.bin") {
	//	w.DeserializeFromFile("world.bin")
	//}

	playbackFile := ""
	// Change the file name to choose a specific playthrough to play back or
	// comment this line to disable playback.
	//playbackFile = "recordings/recorded-inputs-2024-01-02-000004"

	var playbackInputs []Input
	if playbackFile != "" {
		playbackInputs = DeserializeInputs(playbackFile)
	}

	DeleteFile("world-ready")
	// Race condition here: the GUI might already have started reading
	// world.bin, before we got a chance to tell it to stop.
	// So wait until it is done with that, if it is ever done with that.
	time.Sleep(50 * time.Millisecond)
	w.SerializeToFile("world.bin")
	TouchFile("world-ready")
	recordingFile := GetNewRecordingFile()
	frameIdx := 0
	for ; !input.Player1Input.Quit &&
		!input.Player2Input.Quit &&
		(playbackFile == "" || frameIdx < len(playbackInputs)); frameIdx++ {
		WaitForFile("input-ready")
		if playbackFile != "" {
			input = playbackInputs[frameIdx]
		} else {
			input.DeserializeFromFile("input.bin")
		}
		currentInputs = append(currentInputs, input)
		SerializeInputs(currentInputs, recordingFile)

		w.Step(&input, frameIdx)
		w.SerializeToFile("world.bin")
		DeleteFile("input-ready")
		TouchFile("world-ready")
	}
}
