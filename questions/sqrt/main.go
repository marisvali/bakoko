/*
Questions:
- I found a fast square root algorithm for integers on the internet.
Is it correct?
- Is it faster or slower than converting back and forth from float64?
- Does it need checking for overflow?

Conclusions:
- The algorithm is correct.
- It's 7 times slower than using float operations.

Typical output:
results are correct
int sqrt is 6.925699 times slower than float sqrt
*/

package main

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"time"
)

func smallestPowerOfTwoBiggerThan(x int64) int64 {
	one := int64(1)
	for one <= x {
		one <<= 1
	}
	return one
}
func main() {
	for i := int64(10000000); i < 100000000; i++ {
		arr := sqrtVerbose(i)
		for _, val := range arr {
			if val > i {
				if val != smallestPowerOfTwoBiggerThan(i) {
					fmt.Printf("found one: %d %d %d\n", i, val,
						smallestPowerOfTwoBiggerThan(i))
				}
			}
		}
	}
	fmt.Println("\ndone")
	//intSqrtTime, intSqrtRes := measure(intSqrt)
	//floatSqrtTime, floatSqrtRes := measure(floatSqrt)
	//for i := 0; i < len(intSqrtRes); i++ {
	//	if intSqrtRes[i] != floatSqrtRes[i] {
	//		panic(fmt.Errorf("found difference at %d: %d %d", i,
	//			intSqrtRes[i], floatSqrtRes[i]))
	//	}
	//}
	//fmt.Printf("\nresults are correct")
	//fmt.Printf("\nint sqrt is %f times slower than float sqrt",
	//	intSqrtTime/floatSqrtTime)
}

//const loopStart = 100000000
//const loopEnd = 1000000000

const loopStart = 0
const loopEnd = 1000

func intSqrt() []int64 {
	res := make([]int64, loopEnd-loopStart)
	x := int64(loopStart)
	for ; x < loopEnd; x++ {
		res[x-loopStart] = sqrt(x)
	}
	return res
}

func floatSqrt() []int64 {
	res := make([]int64, loopEnd-loopStart)
	x := int64(loopStart)
	for ; x < loopEnd; x++ {
		res[x-loopStart] = int64(math.Sqrt(float64(x)))
	}
	return res
}

func sqrt(a int64) int64 {
	op := a
	res := int64(0)
	// The second-to-top bit is set: use 1 << 14 for int16; use 1 << 30 for
	// int32.
	one := int64(1) << 62

	// "one" starts at the highest power of four <= than the argument.
	for one > op {
		one >>= 2
	}

	for one != 0 {
		if op >= res+one {
			op = op - (res + one)
			res = res + 2*one
		}
		res >>= 1
		one >>= 2
	}
	return res
}

func sqrtVerbose(a int64) []int64 {
	var all []int64
	op := a
	res := int64(0)
	all = append(all, res)
	// The second-to-top bit is set: use 1 << 14 for int16; use 1 << 30 for
	// int32.
	one := int64(1) << 62 // for sure not overflowing
	//all = append(all, one)

	// "one" starts at the highest power of four <= than the argument.
	for one > op {
		one >>= 2 // for sure not overflowing
		//all = append(all, one)
	}

	for one != 0 {
		if op >= res+one {
			all = append(all, res+one)
			op = op - (res + one)
			all = append(all, op)
			res = res + 2*one
			all = append(all, res)
		}
		res >>= 1
		all = append(all, res)
		one >>= 2
		all = append(all, res)
	}
	return all
}

// ---------------------------UTILITIES----------------------------------------
func measure(f func() []int64) (float64, []int64) {
	start := time.Now()
	//result := f().val
	result := f()
	elapsed := time.Since(start).Seconds()
	//fmt.Printf("\n--------------\n%s\nresult: %d\nelapsed: %f", funcName(f),
	//	result, elapsed.Seconds())
	return elapsed, result
}

func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}
