package ints

import "fmt"

type Int struct {
	val int64
}

func (a Int) Lt(b Int) bool {
	return a.val < b.val
}

func (a Int) Leq(b Int) bool {
	return a.val <= b.val
}

func (a Int) Eq(b Int) bool {
	return a.val == b.val
}

func (a Int) Neq(b Int) bool {
	return a.val != b.val
}

func (a Int) Gt(b Int) bool {
	return a.val > b.val
}

func (a Int) Geq(b Int) bool {
	return a.val >= b.val
}

func (a Int) Plus(b Int) Int {
	c := Int{a.val + b.val}
	if (c.val > a.val) == (b.val > 0) {
		return c
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
}

func (a Int) Minus(b Int) Int {
	c := Int{a.val - b.val}
	if (c.val < a.val) == (b.val > 0) {
		return c
	}
	panic(fmt.Errorf("subtraction overflow: %d %d", a, b))
}

func (a Int) Times(b Int) Int {
	if a.val == 0 || b.val == 0 {
		return Int{0}
	}

	c := Int{a.val * b.val}
	if (c.val < 0) == ((a.val < 0) != (b.val < 0)) {
		if c.val/b.val == a.val {
			return c
		}
	}
	panic(fmt.Errorf("multiplicatin overflow: %d %d", a, b))
}

func (a Int) DivBy(b Int) Int {
	if b.val == 0 {
		panic(fmt.Errorf("division by zero: %d %d", a, b))
	}

	c := Int{a.val / b.val}
	if (c.val < 0) != ((a.val < 0) != (b.val < 0)) {
		panic(fmt.Errorf("division overflow: %d %d", a, b))
	}
	return c
}

func (a Int) Sqr() Int {
	return Int{a.val * a.val}
}
