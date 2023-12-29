package bakoko

import . "playful-patterns.com/bakoko/ints"

type Pt struct {
	X, Y Int
}

func (p *Pt) SquaredDistTo(other Pt) Int {
	dx := p.X.Minus(other.X)
	dy := p.Y.Minus(other.Y)
	return dx.Sqr().Plus(dy.Sqr())
}

func (p *Pt) Add(other Pt) {
	p.X = p.X.Plus(other.X)
	p.Y = p.Y.Plus(other.Y)
}

func (p *Pt) Scale(multiply Int, divide Int) Pt {
	p.X = p.X.Times(multiply).DivBy(divide)
	p.Y = p.Y.Times(multiply).DivBy(divide)
	return *p
}

func (p *Pt) AddLen(extraLen Int) {
	oldLenSquared := p.SquaredLen()
	newLenSquared := oldLenSquared.Minus(extraLen.Sqr())
	if newLenSquared.Lt(I(0)) {
		newLenSquared = I(0)
	}
	p.Scale(newLenSquared, oldLenSquared)
}

func (p *Pt) SquaredLen() Int {
	return p.X.Sqr().Plus(p.Y.Sqr())
}

func (p *Pt) Len() Int {
	return p.SquaredLen().Sqrt()
}

func (p *Pt) To(other Pt) Pt {
	return Pt{other.X.Minus(p.X), other.Y.Minus(p.Y)}
}

func (p *Pt) SetLen(newLen Int) Pt {
	oldLen := p.Len()
	p.X = p.X.Times(newLen).DivBy(oldLen)
	p.Y = p.Y.Times(newLen).DivBy(oldLen)
	return *p
}
