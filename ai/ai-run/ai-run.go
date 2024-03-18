package ai_run

import (
	. "playful-patterns.com/bakoko/ai"
	. "playful-patterns.com/bakoko/networking"
	. "playful-patterns.com/bakoko/world"
	"time"
)

func RunAi(guiProxy GuiProxy, worldProxy WorldProxy) {
	var w *World
	var ai PlayerAI
	ai.PauseBetweenShots = 1500 * time.Millisecond
	ai.LastShot = time.Now()

	for {
		input := ai.Step(w)

		// This should block as the AI doesn't make sense if it doesn't
		// synchronize with the simulation.
		for {
			if err := worldProxy.Connect(); err != nil {
				continue // Retry from the beginning.
			}

			if err := worldProxy.SendInput(&input); err != nil {
				continue // Retry from the beginning.
			}

			var err error
			if w, err = worldProxy.GetWorld(); err != nil {
				continue // Retry from the beginning.
			}

			break
		}

		// This may or may not block, who cares?
		//guiProxy.SendPaintData(&ai.DebugInfo)
	}
}
