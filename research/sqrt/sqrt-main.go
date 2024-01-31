/*
Questions:
- I found a fast square root algorithm for integers on the internet.
Is it correct?
- Is it faster or slower than converting back and forth from float64?
- Does it need checking for overflow?
- I found another algorithm, is it better?
- How much slower would it be if I just checked for overflow?

Conclusions:
- The first algorithm is correct.
- It's 5 times slower than using float operations, without checking for overflow.
- I cannot determine that it doesn't overflow during its operation just by
reading the algorithm. I don't have the necessary math skills and/or I am not
currently willing to invest the time in this. However, I've done several tests
and they all show that it doesn't overflow. Plus the source of the algorithm
seems to be the kind of source that would consider this kind of thing.
However, I don't know 100% that it doesn't overflow.
- The second algorithm is correct and is 9 times slower than using float
operations. However, it's easy for me to tell just by reading it that it never
overflows.
- Algorithm 1 with checks for overflow is 135 times slower than using float
operations.

Typical output:
results for values between 100000 and 10000000 are correct
int sqrt (1) is 4.871541 times slower than float sqrt
float sqrt: 2147483648.000000000000
float sqrt (truncated to int): 2147483648
int sqrt (1): 2147483647
results for top 100000000 int64 values are correct
results for 100000000 other large int64 values are correct
for values between 1000000 and 10000000 the biggest intermediary value generated during the execution of the algorithm was the smallest power of two bigger than the input number
results for values between 100000 and 10000000 are correct
int sqrt (2) is 8.966365 times slower than float sqrt
results for values between 100000 and 10000000 are correct
int sqrt (2) is 135.199295 times slower than float sqrt
*/

package main

import (
	"fmt"
	"math"
	. "playful-patterns.com/bakoko/ints"
	"reflect"
	"runtime"
	"time"
)

func main() {
	intSqrtTime, intSqrtRes := measure(intSqrt)
	floatSqrtTime, floatSqrtRes := measure(floatSqrt)
	for i := 0; i < len(intSqrtRes); i++ {
		if intSqrtRes[i] != floatSqrtRes[i] {
			panic(fmt.Errorf("found difference at %d: %d %d", i,
				intSqrtRes[i], floatSqrtRes[i]))
		}
	}
	fmt.Printf("\nresults for values between %d and %d are correct",
		loopStart, loopEnd)
	fmt.Printf("\nint sqrt (1) is %f times slower than float sqrt",
		intSqrtTime/floatSqrtTime)

	// This is a very interesting case. The Windows calculator gives me the
	// result: 2,147,483,647.9999999997671693563461
	// Which shows me that this integer method is actually more precise than
	// the floating point sqrt, in some sense.
	fmt.Printf("\nfloat sqrt: %.12f", math.Sqrt(float64(math.MaxInt64/2)))                        // 2147483648.000000000000
	fmt.Printf("\nfloat sqrt (truncated to int): %d", int64(math.Sqrt(float64(math.MaxInt64/2)))) // 2147483648
	fmt.Printf("\nint sqrt (1): %d", sqrt(math.MaxInt64/2))                                       // 2147483647

	// check correctness for top 100000000 int64 values
	for i := 0; i < 100000000; i++ {
		val := int64(math.MaxInt64 - i)
		res1 := sqrt(uint64(val))
		res2 := uint32(math.Sqrt(float64(val)))
		if res1 != res2 {
			panic(fmt.Errorf("found difference at %d: %d %d", val, res1, res2))
		}
	}
	fmt.Printf("\nresults for top 100000000 int64 values are correct")

	// check correctness for some large "random" values
	for i := 0; i < 100000000; i++ {
		val := int64(math.MaxInt64 - 100000000 - i*1234)
		res1 := sqrt(uint64(val))
		res2 := uint32(math.Sqrt(float64(val)))
		if math.Abs(float64(res1)-float64(res2)) > 1 { // have to account for
			// the fact that float sqrt has rounding errors and int sqrt doesn't
			panic(fmt.Errorf("found difference at %d: %d %d", val, res1, res2))
		}
	}
	fmt.Printf("\nresults for 100000000 other large int64 values are correct")

	for i := uint64(1000000); i < 10000000; i++ {
		arr := sqrtVerbose(i)
		for _, val := range arr {
			if val > i {
				if val != smallestPowerOfTwoBiggerThan(i) {
					panic(fmt.Errorf("found one: %d %d %d\n", i, val,
						smallestPowerOfTwoBiggerThan(i)))
				}
			}
		}
	}
	fmt.Printf("\nfor values between 1000000 and 10000000 the biggest " +
		"intermediary value generated during the execution of the algorithm was" +
		" the smallest power of two bigger than the input number")

	intSqrt2Time, intSqrt2Res := measure(intSqrt2)
	for i := 0; i < len(intSqrt2Res); i++ {
		if intSqrt2Res[i] != floatSqrtRes[i] {
			panic(fmt.Errorf("found difference at %d: %d %d", i,
				intSqrt2Res[i], floatSqrtRes[i]))
		}
	}
	fmt.Printf("\nresults for values between %d and %d are correct",
		loopStart, loopEnd)
	fmt.Printf("\nint sqrt (2) is %f times slower than float sqrt",
		intSqrt2Time/floatSqrtTime)

	intSqrtOverflowTime, intSqrtOverflowRes := measure(intSqrtOverflow)
	for i := 0; i < len(intSqrtOverflowRes); i++ {
		if intSqrtOverflowRes[i] != floatSqrtRes[i] {
			panic(fmt.Errorf("found difference at %d: %d %d", i,
				intSqrt2Res[i], floatSqrtRes[i]))
		}
	}
	fmt.Printf("\nresults for values between %d and %d are correct",
		loopStart, loopEnd)
	fmt.Printf("\nint sqrt (2) is %f times slower than float sqrt",
		intSqrtOverflowTime/floatSqrtTime)
}

const loopStart = 100000
const loopEnd = 10000000

//const loopStart = math.MaxInt64
//const loopEnd = math.MaxInt64

//const loopStart = 0
//const loopEnd = 1000

func intSqrt() []int64 {
	res := make([]int64, loopEnd-loopStart)
	x := int64(loopStart)
	for ; x < loopEnd; x++ {
		res[x-loopStart] = int64(sqrt(uint64(x)))
	}
	return res
}

func intSqrt2() []int64 {
	res := make([]int64, loopEnd-loopStart)
	x := int64(loopStart)
	for ; x < loopEnd; x++ {
		res[x-loopStart] = int64(sqrt2(uint64(x)))
	}
	return res
}

func intSqrtOverflow() []int64 {
	res := make([]int64, loopEnd-loopStart)
	x := int64(loopStart)
	for ; x < loopEnd; x++ {
		res[x-loopStart] = sqrtOverflow(I(x)).ToInt64()
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

func sqrt(a uint64) uint32 {
	op := a
	res := uint64(0)
	// The second-to-top bit is set: use 1 << 14 for uint16; use 1 << 30 for
	// uint32.
	one := uint64(1) << 62

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
	return uint32(res)
}

func sqrtOverflow(a Int) Int {
	op := a
	res := I(0)
	// The second-to-top bit is set: use 1 << 14 for int16; use 1 << 30 for
	// int32.
	one := I(int64(1) << 62)

	// "one" starts at the highest power of four <= than the argument.
	for one.Gt(op) {
		one = one.DivBy(I(4))
	}

	for one.Neq(I(0)) {
		if op.Geq(res.Plus(one)) {
			op = op.Minus(res.Plus(one))
			res = res.Plus(one.Times(I(2)))
		}
		res = res.DivBy(I(2))
		one = one.DivBy(I(4))
	}
	return res
}

func sqrtVerbose(a uint64) []uint64 {
	var all []uint64
	op := a
	res := uint64(0)
	all = append(all, res)
	// The second-to-top bit is set: use 1 << 14 for int16; use 1 << 30 for
	// int32.
	one := uint64(1) << 62 // for sure not overflowing
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

func sqrt2(x uint64) uint32 {
	res := uint32(0)
	add := uint32(0x80000000)

	for i := 0; i < 32; i++ {
		temp := res | add                 // can never overflow
		g2 := uint64(temp) * uint64(temp) // can never overflow
		if x >= g2 {
			res = temp
		}
		add >>= 1 // can never overflow
	}
	return res
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

func smallestPowerOfTwoBiggerThan(x uint64) uint64 {
	one := uint64(1)
	for one <= x {
		one <<= 1
	}
	return one
}
