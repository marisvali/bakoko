package ints

import (
	"fmt"
	"math"
)

type Int struct {
	// This is made public for the sake of serializing and deserializing
	// using the encoding/binary package.
	// Don't access it otherwise.
	Val int64
}

func I(val int64) Int {
	return Int{val}
}

func (a Int) ToInt64() int64 {
	return a.Val
}

func (a Int) ToFloat64() float64 {
	return float64(a.Val)
}

func (a Int) Lt(b Int) bool {
	return a.Val < b.Val
}

func (a Int) Leq(b Int) bool {
	return a.Val <= b.Val
}

func (a Int) Eq(b Int) bool {
	return a.Val == b.Val
}

func (a Int) Neq(b Int) bool {
	return a.Val != b.Val
}

func (a Int) Gt(b Int) bool {
	return a.Val > b.Val
}

func (a Int) Geq(b Int) bool {
	return a.Val >= b.Val
}

func (a *Int) Inc() {
	if a.Val == math.MaxInt64 {
		panic(fmt.Errorf("increment overflow: %d", a))
	}
	a.Val++
}

func (a *Int) Dec() {
	if a.Val == math.MinInt64 {
		panic(fmt.Errorf("decrement overflow: %d", a))
	}
	a.Val--
}

func (a Int) Plus(b Int) Int {
	c := Int{a.Val + b.Val}
	if (c.Val > a.Val) == (b.Val > 0) {
		return c
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
}

func (a *Int) Add(b Int) {
	c := a.Val + b.Val
	if (c > a.Val) == (b.Val > 0) {
		a.Val = c
		return
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
}

func (a Int) Minus(b Int) Int {
	c := Int{a.Val - b.Val}
	if (c.Val < a.Val) == (b.Val > 0) {
		return c
	}
	panic(fmt.Errorf("subtraction overflow: %d %d", a, b))
}

func (a *Int) Subtract(b Int) {
	c := a.Val - b.Val
	if (c < a.Val) == (b.Val > 0) {
		a.Val = c
		return
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
}

func (a Int) Times(b Int) Int {
	if a.Val == 0 || b.Val == 0 {
		return Int{0}
	}

	c := Int{a.Val * b.Val}
	if (c.Val < 0) == ((a.Val < 0) != (b.Val < 0)) {
		if c.Val/b.Val == a.Val {
			return c
		}
	}
	panic(fmt.Errorf("multiplicatin overflow: %d %d", a, b))
}

func (a Int) DivBy(b Int) Int {
	if b.Val == 0 {
		panic(fmt.Errorf("division by zero: %d %d", a, b))
	}
	if a.Val == math.MinInt64 && b.Val == -1 {
		panic(fmt.Errorf("division overflow: %d %d", a, b))
	}
	return Int{a.Val / b.Val}
}

func (a Int) Sqr() Int {
	return Int{a.Val * a.Val}
}

/**
 * \brief    Fast Square root algorithm
 *
 * Fractional parts of the answer are discarded. That is:
 *      - SquareRoot(3) --> 1
 *      - SquareRoot(4) --> 2
 *      - SquareRoot(5) --> 2
 *      - SquareRoot(8) --> 2
 *      - SquareRoot(9) --> 3
 *
 * \param[in] a_nInput - unsigned integer for which to find the square root
 *
 * \return Integer square root of the input value.
 */
//func (a Int) Sqrt() Int {
//	op := a.Val
//	res := int64(0)
//	// The second-to-top bit is set: use 1 << 14 for int16; use 1 << 30 for
//	// int32.
//	one := int64(1) << 62
//
//	// "one" starts at the highest power of four <= than the argument.
//	for one > op {
//		one >>= 2
//	}
//
//	for one != 0 {
//		if op >= res+one {
//			op = op - (res + one)
//			res = res + 2*one
//		}
//		res >>= 1
//		one >>= 2
//	}
//	return Int{res}
//}
// TODO: come back and find an integer-only sqrt algorithm
// the commented method above seems to be fine but can reach numbers above
// the number you give it as input, so I'm not sure when and if it overflows
// might not overflow ever since math.MaxInt64 is a power of two and the
// reaching above is only with the power of two above the input number from
// what I've seen in tests, but who knows if that holds true every time
func (a Int) Sqrt() Int {
	return Int{int64(math.Sqrt(float64(a.Val)))}
}
