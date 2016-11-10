package parallel_test

import (
	"fmt"
	"github.com/buptmiao/parallel"
	"testing"
)

// a function with a single return value
func testFunc2(x, y, z int) string {
	return fmt.Sprintf("%d:%d:%d", x, y, z)
}

// a function with no return value, but receiver the result by res argument
func testFunc3(x, y int, res *int) {
	*res = x + y
}

// a function with multiple return values
func testFunc4(x, y int) (int, int) {
	return y, x
}

func TestNewPipeline(t *testing.T) {
	pipe := parallel.NewPipeline()
	var res string
	pipe.Register(testFunc2, 1, 2, 3).SetReceivers(&res)
	var expect2 int
	pipe.Register(testFunc3, 1, 1, &expect2)
	var x, y int
	pipe.Register(testFunc4, 1, 2).SetReceivers(&x, &y)
	pipe.Do()
	if res != "1:2:3" || expect2 != 2 || x != 2 || y != 1 {
		panic("unexpected result")
	}
}
