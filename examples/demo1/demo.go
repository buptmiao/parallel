package main

import (
	"github.com/buptmiao/parallel"
)

type middle struct {
	B int
	C int
}

type testResult struct {
	A string
	M middle
}

func testJobA() string {
	return "job"
}

func testJobB(x, y int) int {
	return x + y
}

func testJobC(x int) int {
	return -x
}

func testFinal(s *string, m *middle) testResult {
	return testResult{
		*s, *m,
	}
}

func main() {
	var m middle
	var s string
	var res testResult

	p := parallel.NewParallel()

	// Create a child 1
	child1 := p.NewChild()
	child1.Register(testJobA).SetReceivers(&s)

	// Create another child 2
	child2 := p.NewChild()
	child2.Register(testJobB, 1, 2).SetReceivers(&m.B)
	child2.Register(testJobC, 2).SetReceivers(&m.C)

	p.Register(testFinal, &s, &m).SetReceivers(&res)
	// block here
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
