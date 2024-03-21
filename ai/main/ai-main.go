package main

import (
	"os"
	. "playful-patterns.com/bakoko/ai/ai-run"
	. "playful-patterns.com/bakoko/networking"
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
	var aiRunner AiRunner
	aiRunner.Initialize(worldProxy)
	for {
		aiRunner.Step()
	}
}
