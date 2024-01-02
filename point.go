package bakoko

import (
	. "playful-patterns.com/bakoko/ints"
)

type Pt struct {
	X, Y Int
}

func IPt(x, y int64) Pt {
	return Pt{I(x), I(y)}
}

func (p Pt) SquaredDistTo(other Pt) Int {
	return p.To(other).SquaredLen()
}

func (p *Pt) Add(other Pt) {
	p.X = p.X.Plus(other.X)
	p.Y = p.Y.Plus(other.Y)
}

func (p *Pt) Scale(multiply Int, divide Int) {
	p.X = p.X.Times(multiply).DivBy(divide)
	p.Y = p.Y.Times(multiply).DivBy(divide)
}

func (p Pt) SquaredLen() Int {
	return p.X.Sqr().Plus(p.Y.Sqr())
}

func (p Pt) Len() Int {
	return p.SquaredLen().Sqrt()
}

func (p Pt) To(other Pt) Pt {
	return Pt{other.X.Minus(p.X), other.Y.Minus(p.Y)}
}

func (p *Pt) SetLen(newLen Int) {
	oldLen := p.Len()
	if oldLen.Eq(I(0)) {
		return
	}
	p.Scale(newLen, oldLen)
}

func (p *Pt) AddLen(extraLen Int) {
	oldLen := p.Len()
	newLen := oldLen.Plus(extraLen)
	if newLen.Lt(I(0)) {
		newLen = I(0)
	}
	p.Scale(newLen, oldLen)
}
