/*
Questions:
- How much slower is it to check for overflow for each int operation?
- I want to impose that all my int operations are done using some methods that
define. For that I can use a struct with some methods defined for it,
where the struct only contains an int64.
Is using such a struct slower than using int64 directly?
- Is using a trivial function instead of the built-in operator slower (e.g.
'add(a, b)' vs. a+b)?

Conclusions:
- Checking for overflow makes int operations 3 to 6.5 times slower.
- Using a struct with an int instead of an int has no effect.
- Using a trivial function instead of the built-in operator has no effect.

Typical output:
--------------
main.addition
result: 3000000004
elapsed: 0.266202
--------------
main.additionFunc
result: 3000000004
elapsed: 0.266429
--------------
main.additionWithStruct
elapsed: 0.265107
--------------
main.additionWithOverflowChecking
result: 3000000004
elapsed: 1.713044
--------------
main.additionWithStructAndOverflowChecking
result: 3000000004
elapsed: 1.742316
--------------
main.complex
result: 833333458322002
elapsed: 1.648578
--------------
main.complexFunc
result: 833333458322002
elapsed: 1.653952
--------------
main.complexWithStruct
result: 833333458322002
elapsed: 1.619030
--------------
main.complexWithOverflowChecking
result: 833333458322002
elapsed: 4.667535
--------------
main.complexWithStructAndOverflowChecking
result: 833333458322002
elapsed: 4.862702


checked add is 6.545080 slower than basic add
checked complex is 2.949634 slower than basic complex
*/

package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func main() {
	basicAdd := measure(addition)
	measure(additionFunc)
	measure(additionWithStruct)
	measure(additionWithOverflowChecking)
	checkedAdd := measure(additionWithStructAndOverflowChecking)

	basicComplex := measure(complex)
	measure(complexFunc)
	measure(complexWithStruct)
	measure(complexWithOverflowChecking)
	checkedComplex := measure(complexWithStructAndOverflowChecking)

	fmt.Printf("\n\n\nchecked add is %f slower than basic add", checkedAdd.Seconds()/basicAdd.Seconds())
	fmt.Printf("\nchecked complex is %f slower than basic complex", checkedComplex.Seconds()/basicComplex.Seconds())
}

const loopSizeAdd = 1000000000
const loopSizeComplex = 100000000

// ---------------------------------SIMPLE-------------------------------------
func addition() int64 {
	var a, b int64
	a = 3
	b = 4
	c := a + b
	for i := 1; i < loopSizeAdd; i++ {
		c = c + a
	}
	return c
}

func complex() int64 {
	var x, y, t int64
	x = 1000
	y = 2134
	t = 0
	i := int64(10000)
	for ; i < loopSizeComplex; i++ {
		oldLen := x*x + y*y
		x = x * i / oldLen
		y = y * i / oldLen
		t = t + x + y
	}
	return t
}

// --------------------------TRIVIAL-FUNCTION-----------------------------------
func additionFunc() int64 {
	var a, b int64
	a = 3
	b = 4
	c := addTrivial(a, b)
	for i := 1; i < loopSizeAdd; i++ {
		c = addTrivial(c, a)
	}
	return c
}

func complexFunc() int64 {
	var x, y, t int64
	x = 1000
	y = 2134
	t = 0
	i := int64(10000)
	for ; i < loopSizeComplex; i++ {
		oldLen := addTrivial(mulTrivial(x, x), mulTrivial(y, y))
		x = divTrivial(mulTrivial(x, i), oldLen)
		y = divTrivial(mulTrivial(y, i), oldLen)
		t = addTrivial(addTrivial(t, x), y)
	}
	return t
}

func addTrivial(a, b int64) int64 {
	return a + b
}

func mulTrivial(a, b int64) int64 {
	return a * b
}

func divTrivial(a, b int64) int64 {
	return a / b
}

// ---------------------------------STRUCT-------------------------------------
func additionWithStruct() int64 {
	var a, b mov
	a = mov{3}
	b = mov{4}
	c := mov{a.val + b.val}
	for i := 1; i < loopSizeAdd; i++ {
		c.val = c.val + a.val
	}
	return c.val
}

func complexWithStruct() int64 {
	var x, y, t mov
	x = mov{1000}
	y = mov{2134}
	t = mov{0}
	i := mov{10000}
	for ; i.val < loopSizeComplex; i.val++ {
		oldLen := mov{mov{x.val * x.val}.val + mov{y.val * y.val}.val}
		x = mov{mov{x.val * i.val}.val / oldLen.val}
		y = mov{mov{y.val * i.val}.val / oldLen.val}
		t = mov{mov{t.val + x.val}.val + y.val}
	}
	return t.val
}

type mov struct {
	val int64
}

// ----------------------------OVERFLOW-CHECKING-------------------------------
func additionWithOverflowChecking() int64 {
	var a, b int64
	a = 3
	b = 4
	c := add(a, b)
	for i := 1; i < loopSizeAdd; i++ {
		c = add(c, a)
	}
	return c
}

func complexWithOverflowChecking() int64 {
	var x, y, t int64
	x = 1000
	y = 2134
	t = 0
	i := int64(10000)
	for ; i < loopSizeComplex; i++ {
		oldLen := add(mul(x, x), mul(y, y))
		x = div(mul(x, i), oldLen)
		y = div(mul(y, i), oldLen)
		t = add(add(t, x), y)
	}
	return t
}

func add(a, b int64) int64 {
	c := a + b
	if (c > a) == (b > 0) {
		return c
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
	return 0
}

func mul(a int64, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}

	c := a * b
	if (c < 0) == ((a < 0) != (b < 0)) {
		if c/b == a {
			return c
		}
	}
	panic(fmt.Errorf("multiplicatin overflow: %d %d", a, b))
	return 0 // should never be reached
}

func div(a, b int64) int64 {
	if b == 0 {
		panic(fmt.Errorf("division by zero: %d %d", a, b))
	}

	c := a / b
	if (c < 0) != ((a < 0) != (b < 0)) {
		panic(fmt.Errorf("division overflow: %d %d", a, b))
	}
	return c
}

// ---------------------STRUCT-AND-OVERFLOW-CHECKING---------------------------
func additionWithStructAndOverflowChecking() int64 {
	var a, b mov
	a = mov{3}
	b = mov{4}
	c := addMov(a, b)
	for i := 1; i < loopSizeAdd; i++ {
		c = addMov(c, a)
	}
	return c.val
}

func complexWithStructAndOverflowChecking() int64 {
	var x, y, t mov
	x = mov{1000}
	y = mov{2134}
	t = mov{0}
	i := mov{10000}
	for ; i.val < loopSizeComplex; i.val++ {
		oldLen := addMov(mulMov(x, x), mulMov(y, y))
		x = divMov(mulMov(x, i), oldLen)
		y = divMov(mulMov(y, i), oldLen)
		t = addMov(addMov(t, x), y)
	}
	return t.val
}

func addMov(a, b mov) mov {
	c := mov{a.val + b.val}
	if (c.val > a.val) == (b.val > 0) {
		return c
	}
	panic(fmt.Errorf("addition overflow: %d %d", a.val, b.val))
	return mov{0}
}

func mulMov(a mov, b mov) mov {
	if a.val == 0 || b.val == 0 {
		return mov{0}
	}

	c := mov{a.val * b.val}
	if (c.val < 0) == ((a.val < 0) != (b.val < 0)) {
		if c.val/b.val == a.val {
			return c
		}
	}
	panic(fmt.Errorf("multiplicatin overflow: %d %d", a, b))
	return mov{0} // should never be reached
}

func divMov(a, b mov) mov {
	if b.val == 0 {
		panic(fmt.Errorf("division by zero: %d %d", a, b))
	}

	c := mov{a.val / b.val}
	if (c.val < 0) != ((a.val < 0) != (b.val < 0)) {
		panic(fmt.Errorf("division overflow: %d %d", a, b))
	}
	return c
}

// ---------------------------UTILITIES----------------------------------------
func measure(f func() int64) time.Duration {
	start := time.Now()
	result := f()
	elapsed := time.Since(start)
	fmt.Printf("\n--------------\n%s\nresult: %d\nelapsed: %f", funcName(f), result, elapsed.Seconds())
	return elapsed
}

func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}
