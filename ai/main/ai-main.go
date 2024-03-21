package main

import (
	"os"
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"time"
)

func main() {
	var worldProxy WorldProxyTcpIp
	var guiProxy GuiProxyTcpIp

	worldProxy.Endpoint = os.Args[1] // localhost:56901 or localhost:56902
	worldProxy.Timeout = 0 * time.Millisecond
	guiProxy.Endpoint = os.Args[2]

	RunAiSplitPlay(&guiProxy, &worldProxy)
}

func RunAiSplitPlay(guiProxy GuiProxy, worldProxy WorldProxy) {
	var ai PlayerAI
	ai.Initialize()
	for {
		w := getWorld(worldProxy)
		input := ai.Step(w)
		sendInput(worldProxy, &input)

		// This may or may not block, who cares?
		//guiProxy.SendPaintData(&Ai.DebugInfo)
	}
}

func getWorld(worldProxy WorldProxy) *World {
	// This should block as the AI doesn't make sense if it doesn't
	// synchronize with the simulation.
	for {
		if err := worldProxy.Connect(); err != nil {
			continue // Retry from the beginning.
		}
		var err error
		var w *World
		if w, err = worldProxy.GetWorld(); err != nil {
			continue // Retry from the beginning.
		}
		return w
	}
}

func sendInput(worldProxy WorldProxy, input *PlayerInput) {
	// This should block as the AI doesn't make sense if it doesn't
	// synchronize with the simulation.
	for {
		if err := worldProxy.Connect(); err != nil {
			continue // Retry from the beginning.
		}
		if err := worldProxy.SendInput(input); err != nil {
			continue // Retry from the beginning.
		}
		break
	}
}
