package parallel_test

import (
	"fmt"
	"github.com/buptmiao/parallel"
	"testing"
)

/*

jobA  jobB   jobC
 \      \     /
  \      \   /
   \      middle
    \      /
     \    /
     final

*/

type middle struct {
	B int
	C int
}

type testResult struct {
	A string
	M middle
}

func testJobA() string {
	return fmt.Sprintf("job")
}

func testJobB(x, y int) int {
	return x + y
}

func testJobC(x int) int {
	return -x
}

func TestNewParallel(t *testing.T) {
	p := parallel.NewParallel()
	var res testResult
	p.Register(testJobA).SetReceivers(&res.A)
	child := p.NewChild()
	child.Register(testJobB, 1, 2).SetReceivers(&res.M.B)
	child.Register(testJobC, 2).SetReceivers(&res.M.C)
	p.Run()
	expect := testResult{
		"job",
		middle{
			3, -2,
		},
	}
	if res != expect {
		panic("unexpected result")
	}
}

func TestParallelPanic(t *testing.T) {
	p := parallel.NewParallel()
	p.Register(testJobA).SetReceivers()
	s := make(chan struct{}, 1)
	go func() {
		defer EatPanic(s)
		p.Run()
	}()
	<-s
}
