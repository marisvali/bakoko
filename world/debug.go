package world

import (
	"bytes"
	"image/color"
	. "playful-patterns.com/bakoko/ints"
	"slices"
)

type DebugPoint struct {
	Pos  Pt
	Size Int
	Col  color.RGBA
}

type DebugLine struct {
	Line
	Col color.RGBA
}

type DebugCircle struct {
	Circle
	Col color.RGBA
}

type DebugSquare struct {
	Square
	Col color.RGBA
}
type DebugInfo struct {
	Points  []DebugPoint
	Lines   []DebugLine
	Circles []DebugCircle
	Squares []DebugSquare
}

func (d *DebugInfo) Serialize() []byte {
	buf := new(bytes.Buffer)
	SerializeSlice(buf, d.Points)
	SerializeSlice(buf, d.Lines)
	SerializeSlice(buf, d.Circles)
	SerializeSlice(buf, d.Squares)
	return buf.Bytes()
}

func (d *DebugInfo) Deserialize(buf *bytes.Buffer) {
	DeserializeSlice(buf, &d.Points)
	DeserializeSlice(buf, &d.Lines)
	DeserializeSlice(buf, &d.Circles)
	DeserializeSlice(buf, &d.Squares)
}

func (d *DebugInfo) Clone() (c DebugInfo) {
	c.Points = slices.Clone(d.Points)
	c.Lines = slices.Clone(d.Lines)
	c.Circles = slices.Clone(d.Circles)
	c.Squares = slices.Clone(d.Squares)
	return
}
