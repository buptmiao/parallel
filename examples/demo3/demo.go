package main

import (
	"fmt"
	"github.com/buptmiao/parallel"
)

func testJobB(x ...int) int {
	res := 0
	for _, v := range x {
		res += v
	}
	return res
}

func main() {
	var x int

	p := parallel.NewParallel()

	p.Register(testJobB, 1, 2).SetReceivers(&x)
	// block here
	p.Run()
	fmt.Println(x)

	if x != 3 {
		panic("unexpected result")
	}
}
