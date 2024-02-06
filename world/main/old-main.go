package main

//var input Input
//func main77() {
//	//if utils.FileExists("world-main.bin") {
//	//	w.DeserializeFromFile("world-main.bin")
//	//}
//
//	playbackFile := ""
//	// Change the file name to choose a specific playthrough to play back or
//	// comment this line to disable playback.
//	//playbackFile = "recordings/recorded-inputs-2024-01-02-000004"
//
//	var playbackInputs []Input
//	if playbackFile != "" {
//		playbackInputs = DeserializeInputs(playbackFile)
//	}
//
//	DeleteFile("world-main-ready")
//	// Race condition here: the GUI might already have started reading
//	// world-main.bin, before we got a chance to tell it to stop.
//	// So wait until it is done with that, if it is ever done with that.
//	time.Sleep(50 * time.Millisecond)
//	w.SerializeToFile("world-main.bin")
//	TouchFile("world-main-ready")
//	recordingFile := GetNewRecordingFile()
//	frameIdx := 0
//	for ; !input.Player1Input.Quit &&
//		!input.Player2Input.Quit &&
//		(playbackFile == "" || frameIdx < len(playbackInputs)); frameIdx++ {
//		WaitForFile("input-ready")
//		if playbackFile != "" {
//			input = playbackInputs[frameIdx]
//		} else {
//			input.DeserializeFromFile("input.bin")
//		}
//		currentInputs = append(currentInputs, input)
//		SerializeInputs(currentInputs, recordingFile)
//
//		w.Step(&input, frameIdx)
//		w.SerializeToFile("world-main.bin")
//		DeleteFile("input-ready")
//		TouchFile("world-main-ready")
//	}
//}
//
//func level1() (m Matrix) {
//	m.Init(I(15), I(15))
//	for row := I(0); row.Lt(m.NRows()); row.Inc() {
//		m.Set(row, I(0), I(1))
//		m.Set(row, m.NCols().Minus(I(1)), I(1))
//	}
//	for col := I(0); col.Lt(m.NCols()); col.Inc() {
//		m.Set(I(0), col, I(1))
//		m.Set(m.NRows().Minus(I(1)), col, I(1))
//	}
//	m.Set(I(5), I(5), I(1))
//	m.Set(I(8), I(7), I(1))
//
//	m.Set(I(5), I(10), I(1))
//	m.Set(I(6), I(10), I(1))
//	m.Set(I(7), I(10), I(1))
//
//	m.Set(I(10), I(10), I(1))
//	m.Set(I(11), I(10), I(1))
//	m.Set(I(12), I(10), I(1))
//	return
//}
//
//func main5() {
//	originalWorld := w
//	w.SerializeToFile("world-main.bin")
//	inputs := DeserializeInputs("recorded-inputs")
//	start := time.Now()
//	bigIdx := 0
//	for ; bigIdx < 1000000; bigIdx++ {
//		frameIdx := 0
//		w := originalWorld
//		for ; !input.Player1Input.Quit; frameIdx++ {
//			input = Input{}
//			if frameIdx < len(inputs) {
//				input = inputs[frameIdx]
//			}
//
//			w.Step(&input, frameIdx)
//		}
//	}
//	fmt.Println(time.Since(start).Seconds())
//	fmt.Printf("%.12f\n", time.Since(start).Seconds()/float64(bigIdx))
//	w.SerializeToFile("world-main.bin")
//}
//func main4() {
//	w.SerializeToFile("world-main.bin")
//	inputs := DeserializeInputs("recorded-inputs")
//
//	frameIdx := 0
//	for ; !input.Player1Input.Quit; frameIdx++ {
//		//frameStart := time.Now()
//		WaitForFile("input-ready")
//		//input.DeserializeFromFile("input.bin")
//		input = Input{}
//		if frameIdx < len(inputs) {
//			input = inputs[frameIdx]
//		}
//
//		w.Step(&input, frameIdx)
//		w.SerializeToFile("world-main.bin")
//		DeleteFile("input-ready")
//		TouchFile("world-main-ready")
//
//		//frameDuration := time.Since(frameStart)
//		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
//	}
//}
//func main3() {
//	w.SerializeToFile("world-main.bin")
//
//	frameIdx := 0
//	for ; !input.Player1Input.Quit; frameIdx++ {
//		//frameStart := time.Now()
//		WaitForFile("input-ready")
//		//input.DeserializeFromFile("input.bin")
//		input = Input{}
//
//		input.Player1Input.MoveLeft = rand.Int()%6 == 0
//		input.Player1Input.MoveRight = rand.Int()%5 == 0
//		input.Player1Input.MoveUp = rand.Int()%6 == 0
//		input.Player1Input.MoveDown = rand.Int()%6 == 0
//		input.Player1Input.Quit = frameIdx == 3600
//
//		w.Step(&input, frameIdx)
//		w.SerializeToFile("world-main.bin")
//		DeleteFile("input-ready")
//		TouchFile("world-main-ready")
//
//		//frameDuration := time.Since(frameStart)
//		//fmt.Printf("frame idx: %d duration: %d\n", frameIdx, frameDuration.Milliseconds())
//	}
//}
//
//var currentInputs []Input
//
//func GetNewRecordingFile() string {
//	date := time.Now()
//	for i := 0; i < 1000000; i++ {
//		filename := fmt.Sprintf("recordings/recorded-inputs-%04d-%02d-%02d-%06d",
//			date.Year(), date.Month(), date.Day(), i)
//		if !FileExists(filename) {
//			return filename
//		}
//	}
//	panic("Cannot record, no available filename found.")
//}
