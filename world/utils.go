package world

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

var CheckCrashes = true
var CheckFailed error

func Check(e error) {
	if e != nil {
		CheckFailed = e
		if CheckCrashes {
			panic(e)
		}
	}
}

func WriteFile(name string, data []byte) {
	err := os.WriteFile(name, data, 0644)
	Check(err)
}

func ReadFile(name string) []byte {
	data, err := os.ReadFile(name)
	Check(err)
	return data
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}

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

//
//func TouchFile(name string) {
//	name = "e:/" + name
//	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
//	Check(err)
//	err = file.Close()
//	Check(err)
//}
//
//func FileExists(name string) bool {
//	name = "e:/" + name
//	if _, err := os.Stat(name); err == nil {
//		return true
//	}
//	return false
//}
//
//func WaitForFile(name string) {
//	name = "e:/" + name
//	for {
//		if _, err := os.Stat(name); err == nil {
//			for {
//				// Attempt to remove the file until the attempt succeeds.
//				err = os.Remove(name)
//				if err == nil {
//					return
//				}
//			}
//		} else if errors.Is(err, os.ErrNotExist) {
//			// name does not exist
//		} else {
//			Check(err)
//		}
//	}
//}
//
//func DeleteFile(name string) {
//	name = "e:/" + name
//	err := os.Remove(name)
//	if !errors.Is(err, os.ErrNotExist) {
//		Check(err)
//	}
//}

func Serialize(w io.Writer, data any) {
	err := binary.Write(w, binary.LittleEndian, data)
	Check(err)
}

func Deserialize(r io.Reader, data any) {
	err := binary.Read(r, binary.LittleEndian, data)
	Check(err)
}

func SerializeSlice[T any](buf *bytes.Buffer, s []T) {
	Serialize(buf, int64(len(s)))
	Serialize(buf, s)
}

func DeserializeSlice[T any](buf *bytes.Buffer, s *[]T) {
	var lenSlice int64
	Deserialize(buf, &lenSlice)
	*s = make([]T, lenSlice)
	Deserialize(buf, *s)
}

type TimedFunction func()

func Duration(function TimedFunction) float64 {
	start := time.Now()
	function()
	return time.Since(start).Seconds()
}

func ReadAllText(filename string) string {
	file, err := os.Open(filename)
	Check(err)
	bytes, err := io.ReadAll(file)
	Check(err)
	return string(bytes)
}

func LoadJSON(filename string, v any) {
	file, err := os.Open(filename)
	Check(err)
	bytes, err := io.ReadAll(file)
	Check(err)
	err = json.Unmarshal(bytes, v)
	Check(err)
}

type FolderWatcher struct {
	Folder string
	times  []time.Time
}

func (f *FolderWatcher) FolderContentsChanged() bool {
	files, err := os.ReadDir(f.Folder)
	Check(err)
	if len(files) != len(f.times) {
		f.times = make([]time.Time, len(files))
	}
	changed := false
	for idx, file := range files {
		info, err := file.Info()
		Check(err)
		if f.times[idx] != info.ModTime() {
			changed = true
			f.times[idx] = info.ModTime()
		}
	}
	return changed
}
