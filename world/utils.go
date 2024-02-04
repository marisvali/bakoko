package world

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	name = "e:/" + name
	err := os.WriteFile(name, data, 0644)
	Check(err)
}

func ReadFile(name string) []byte {
	name = "e:/" + name
	data, err := os.ReadFile(name)
	Check(err)
	return data
}

func TouchFile(name string) {
	name = "e:/" + name
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	Check(err)
	err = file.Close()
	Check(err)
}

func FileExists(name string) bool {
	name = "e:/" + name
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}

func WaitForFile(name string) {
	name = "e:/" + name
	for {
		if _, err := os.Stat(name); err == nil {
			for {
				// Attempt to remove the file until the attempt succeeds.
				err = os.Remove(name)
				if err == nil {
					return
				}
			}
		} else if errors.Is(err, os.ErrNotExist) {
			// name does not exist
		} else {
			Check(err)
		}
	}
}

func DeleteFile(name string) {
	name = "e:/" + name
	err := os.Remove(name)
	if !errors.Is(err, os.ErrNotExist) {
		Check(err)
	}
}

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
