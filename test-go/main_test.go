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
