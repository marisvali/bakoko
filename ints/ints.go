/*
ints provides integer operations that check for overflow.

The general problem:
You want to do integer operations and you want to make sure you never overflow.
You plan to do too many such operations and you don't want to check each one
or analyze each one to be 100% sure you won't overflow.

The solution provided by this package:
- an Int type is defined, which is just a wrapper for an int64
- use Int in your code when you want to check for overflow
- arithmetic operations on Int can be done only through functions defined
in this package (e.g. c := a.Plus(b) instead of c := a + b)
- the functions do the arithmetic operation and check for overflow
- on overflow, the function panics

This means:
- your code will be full of function calls instead of nice math operators
- you only find out at runtime if you have an overflow
- you only get to crash on overflow

Go doesn't support operator overloading so this is the best solution I could
think of.

Performance hit:
Checking for overflow makes int operations 3 to 6.5 times slower.
Benchmarked on an Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz, 2208 Mhz).

The concrete problem that started this package:
I want to have games where all the world is simulated using only integers. This
means all operations need to be done on integers only. This ensures that my game
will be deterministic and act exactly the same on all processors. This means
recording a playthrough is cheap. I only need the input from the player, which
is always only a few bytes per frame (mouse and keyboard states).

Advantages of deterministic simulations:
- I can see exactly what the player experienced.
- I have perfect debug information. I can re-create the same code execution on
my machine, so I can recreate a bug and debug it.
- I can refactor my simulation code easily. I need to create some playthroughs
stored as a sequence of inputs, Then I run the playthroughs quickly without any
interface and check that the state of the world at the end is the same before
and after any refactoring work.
- I can have automated tests to protect against previous bugs.
- I can develop AI algorithms that play the game and have them play many games
without an interface. If I ever want to analyze a game, I can view it easily.

Other alternatives for recording playthroughs:
Alternative 1: record video.
Pros:
- I can see what the player experienced.
Cons:
- It takes a lot of processing power to record.
- It takes even more processing power to encode the video on the fly.
- Without encoding on the fly the raw video will be very large.
- Even the final encoded video will be large to upload from the test user to
my server and store on my server.
- I don't have any insight on what my variables looked like during the
playthrough, which would be nice for debugging or understanding things.
- I can only analyze the video visually and manually, I can't run algorithms to
extract metrics from it.

Alternative 2: record world state. Either at each frame or every X frames.
Pros:
- Much cheaper than video.
- Less information but more precise.
- I can use algorithms to analyze the information.
Cons:
- I can only see things like the position of every item and character at a
certain moment. I can't see the state of every algorithm (pathfinding, AI etc).
That would require a full memory dump at every frame, which would very
expensive.

Alternative 3: record events and metrics.
Pros:
- This is definitely the cheapest option.
Cons:
- Much less information.
- I have to decide what to measure before I see players play my game. But I
need to see how players play before I understand what might be interesting to
record. From past experience I know that without seeing them play it's very
hard to understand some of the metrics.
- You need many players before you can have trends which can then give you
insight. So now you are spying on a lot of people so that you get a little bit
of insight from each one, instead of deeply analyzing a few people and getting
most of your insight that way.
*/
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
