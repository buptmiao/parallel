package main

import (
	"fmt"
	"github.com/buptmiao/parallel"
)

/*
jobA      jobC
  |        /
  |       /
jobB     /
  |     /
  |    /
  result
*/

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
	var res string
	var x, y int

	p := parallel.NewParallel()
	pipe := p.NewPipeline()
	pipe.Register(testJobA).SetReceivers(&res)
	pipe.Register(testJobB, 1, 2).SetReceivers(&x)
	p.Register(testJobC, 3).SetReceivers(&y)
	// block here
	p.Run()
	fmt.Println(res, x, y) //job 3 -3
}
