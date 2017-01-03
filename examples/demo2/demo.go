package main

import (
	"github.com/buptmiao/parallel"
)

func testJobA() string {
	return "job"
}

func testJobB(x, y int) int {
	return x + y
}

func testJobC(x int) int {
	return -x
}

func main() {
	var s string
	var x, y int

	p := parallel.NewParallel()

	p.Register(testJobA).SetReceivers(&s)
	p.Register(testJobB, 1, 2).SetReceivers(&x)
	p.Register(testJobC, 3).SetReceivers(&y)
	// block here
	p.Run()

	if s != "job" || x != 3 || y != -3 {
		panic("unexpected result")
	}
}
