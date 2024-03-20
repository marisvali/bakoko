package ai_run

import (
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"time"
)

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

func RunAi(guiProxy GuiProxy, worldProxy WorldProxy) {
	var aiRunner AiRunner
	aiRunner.Initialize(worldProxy)
	for {
		aiRunner.Step()
	}
}
