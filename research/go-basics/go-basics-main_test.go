package bakoko

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"

	//. "playful-patterns.com/bakoko/ints"
	"testing"
)

func TestModOverflow(t *testing.T) {
	a := int64(math.MinInt64)
	//a := int64(-10)
	b := int64(-3)
	fmt.Println(a % b)
	assert.True(t, true)
}

type A struct {
	val int
}

func do(x *[]int) {
	*x = append(*x, 13)

	//x[2] = 13
}

func TestSlices(t *testing.T) {
	var x []int
	x = append(x, 1)
	x = append(x, 3)
	x = append(x, 5)
	x = append(x, 7)
	do(&x)
	fmt.Println(x)
	assert.True(t, true)
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}

// from fib_test.go
func BenchmarkFib10(b *testing.B) {
	// world-run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Fib(10)
	}
}

type Inter1 interface {
	Some()
}

type Abel struct {
	a int
}

func (p *Abel) Some() {
	fmt.Println(p.a)
}

func TestInterfaces(t *testing.T) {
	var x Inter1
	var y Abel
	y.a = 17
	x = &y
	x.Some()

	assert.True(t, true)
}
