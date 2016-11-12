package main

import (
	"github.com/buptmiao/parallel"
	"fmt"
)

func exceptionHandler(topic string, e interface{}) {
	fmt.Println(topic, e)
}

func exceptionJob() {
	var a map[string]int
	//assignment to entry in nil map
	a["123"] = 1
}

func main() {
	p := parallel.NewParallel()
	p.Register(exceptionJob)
	p.Except(exceptionHandler, "topic1")
	p.Run()
}
