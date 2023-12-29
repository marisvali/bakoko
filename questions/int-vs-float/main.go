/*
Questions:
- Are float64 operations slower than int64 operations? If yes, by how much?

Conclusions:
- floats are slower than ints except for division, where floats are much faster
- addition: float is 4.1 times slower than int
- subtraction: float is 4.2 times slower than int
- multiplication: float is 2.6 times slower than int
- division: int is 6.3 times slower than float

Typical output:
--------------
main.addInts
result: 500000000450005000
elapsed: 0.276776
--------------
main.addFloats
result: 500000000017113984
elapsed: 1.141680
--------------
main.subtractInts
result: -499999999450005000
elapsed: 0.273866
--------------
main.subtractFloats
result: -499999999017113984
elapsed: 1.142923
--------------
main.multiplyInts
result: 6500000005850065000
elapsed: 0.453721
--------------
main.multiplyFloats
result: 6500000005368668160
elapsed: 1.204804
--------------
main.divideInts
result: 11090196448
elapsed: 7.211880
--------------
main.divideFloats
result: 11512975466
elapsed: 1.140739


addition: float is 4.124929 times slower than int
subtraction: float is 4.173300 times slower than int
multiplication: float is 2.655386 times slower than int
division: int is 6.322114 times slower than float
*/

package main

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

const loopSize = 1000000000

func main() {
	addIntsTime := measure(addInts)
	addFloatsTime := measure(addFloats)

	subtractIntsTime := measure(subtractInts)
	subtractFloatsTime := measure(subtractFloats)

	multiplyIntsTime := measure(multiplyInts)
	multiplyFloatsTime := measure(multiplyFloats)

	divideIntsTime := measure(divideInts)
	divideFloatsTime := measure(divideFloats)

	fmt.Printf("\n\n")
	fmt.Printf("\naddition: float is %f times slower than int", addFloatsTime.Seconds()/addIntsTime.Seconds())
	fmt.Printf("\nsubtraction: float is %f times slower than int", subtractFloatsTime.Seconds()/subtractIntsTime.Seconds())
	fmt.Printf("\nmultiplication: float is %f times slower than int", multiplyFloatsTime.Seconds()/multiplyIntsTime.Seconds())
	fmt.Printf("\ndivision: int is %f times slower than float", divideIntsTime.Seconds()/divideFloatsTime.Seconds())
}

func addInts() int64 {
	t := int64(0)
	for i := int64(10000); i <= loopSize; i++ {
		t += i
	}
	return t
}

func addFloats() int64 {
	t := float64(0)
	for i := float64(10000); i <= loopSize; i += 1 {
		t += i
	}
	return int64(t)
}

func subtractInts() int64 {
	t := int64(loopSize)
	for i := int64(10000); i <= loopSize; i++ {
		t -= i
	}
	return t
}

func subtractFloats() int64 {
	t := float64(loopSize)
	for i := float64(10000); i <= loopSize; i += 1 {
		t -= i
	}
	return int64(t)
}

func multiplyInts() int64 {
	t := int64(0)
	for i := int64(10000); i <= loopSize; i++ {
		t += i * 13
	}
	return t
}

func multiplyFloats() int64 {
	t := float64(0)
	for i := float64(10000); i <= loopSize; i += 1 {
		t += i * 13
	}
	return int64(t)
}

func divideInts() int64 {
	t := int64(0)
	for i := int64(10000); i <= loopSize; i++ {
		t += loopSize / i
	}
	return t
}

func divideFloats() int64 {
	t := float64(0)
	for i := float64(10000); i <= loopSize; i += 1 {
		t += loopSize / i
	}
	return int64(t)
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
