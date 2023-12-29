/*
Question:
- Is there a performance difference between regular functions and those with
value receivers? For example is add(a, b) faster or slower than a.add(b)?
- Is there a performance difference between value receivers and pointer
receivers for int64?

Conclusions:
- No statistically significant difference between regular functions and those
with value receivers.
- Using functions with pointer receivers which modify the int64 on which the
function is called is faster than functions with value receivers which don't
modify the object.

Typical output:
with value receiver is 0.999264 times slower than without receiver
with value receiver is 0.997096 times slower than without receiver
with value receiver is 1.008823 times slower than without receiver
with value receiver is 0.999684 times slower than without receiver
with value receiver is 1.002552 times slower than without receiver
with value receiver is 0.984321 times slower than without receiver
with value receiver is 0.989265 times slower than without receiver
with value receiver is 1.016460 times slower than without receiver
with value receiver is 1.001799 times slower than without receiver
with value receiver is 1.003736 times slower than without receiver
with value receiver is 1.034518 times slower than with pointer receiver
with value receiver is 1.040276 times slower than with pointer receiver
with value receiver is 1.041045 times slower than with pointer receiver
with value receiver is 1.041861 times slower than with pointer receiver
with value receiver is 1.042384 times slower than with pointer receiver
with value receiver is 1.042170 times slower than with pointer receiver
with value receiver is 1.035322 times slower than with pointer receiver
with value receiver is 1.010789 times slower than with pointer receiver
with value receiver is 1.030126 times slower than with pointer receiver
with value receiver is 1.028087 times slower than with pointer receiver
*/

package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

const loopSize = 100000000

func main() {
	for i := 1; i <= 10; i++ {
		withReceiverTime := measure(withReceiver)
		withoutReceiverTime := measure(withoutReceiver)

		fmt.Printf("\nwith value receiver is %f times slower than without"+
			" receiver", withReceiverTime.Seconds()/withoutReceiverTime.Seconds())
	}

	for i := 1; i <= 10; i++ {
		withReceiverTime := measure(withReceiver)
		withPointerReceiverTime := measure(withPointerReceiver)

		fmt.Printf("\nwith value receiver is %f times slower than with"+
			" pointer receiver", withReceiverTime.Seconds()/withPointerReceiverTime.Seconds())
	}
}

func withReceiver() Int {
	t := Int{0}
	i := Int{10000}
	for ; i.val <= loopSize; i.val++ {
		t = t.add2(i)
	}
	return t
}

func withoutReceiver() Int {
	t := Int{0}
	i := Int{10000}
	for ; i.val <= loopSize; i.val++ {
		t = add1(t, i)
	}
	return t
}

func withPointerReceiver() Int {
	t := Int{0}
	i := Int{10000}
	for ; i.val <= loopSize; i.val++ {
		t.add3(i)
	}
	return t
}

type Int struct {
	val int64
}

func add1(a, b Int) Int {
	c := Int{a.val + b.val}
	if (c.val > a.val) == (b.val > 0) {
		return c
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
	return Int{0}
}

func (a Int) add2(b Int) Int {
	c := Int{a.val + b.val}
	if (c.val > a.val) == (b.val > 0) {
		return c
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
	return Int{0}
}

func (a *Int) add3(b Int) {
	c := Int{a.val + b.val}
	if (c.val > a.val) == (b.val > 0) {
		a.val = c.val
		return
	}
	panic(fmt.Errorf("addition overflow: %d %d", a, b))
}

// ---------------------------UTILITIES----------------------------------------
func measure(f func() Int) time.Duration {
	start := time.Now()
	//result := f().val
	f()
	elapsed := time.Since(start)
	//fmt.Printf("\n--------------\n%s\nresult: %d\nelapsed: %f", funcName(f),
	//	result, elapsed.Seconds())
	return elapsed
}

func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}
