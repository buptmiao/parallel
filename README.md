### Parallel

[![Build Status](https://travis-ci.org/buptmiao/parallel.svg?branch=master)](https://travis-ci.org/buptmiao/parallel)
[![Coverage Status](https://coveralls.io/repos/github/buptmiao/parallel/badge.svg?branch=master)](https://coveralls.io/github/buptmiao/parallel?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/buptmiao/parallel)](https://goreportcard.com/report/github.com/buptmiao/parallel)

[zh_CN](https://github.com/buptmiao/parallel/blob/master/README_CN.md)

A golang parallel library, used for business logic aggregation and refactory without changing function declaration.

### Install

```
go get github.com/buptmiao/parallel
```

### Usage

#### eg.1
There are three methods: testjobA, testjobB, testjobC, execute them in parallel:
```go
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

	if s != "job" || x != 3 || y != -3{
		panic("unexpected result")
	}
}
```
#### eg.2

Let's see a little complex case, there are three parallel jobs: jobA, jobB, jobC and a final Job which aggregates the result. The final depends on jobA and middle which depends on jobB and jobC.

```
jobA  jobB   jobC
 \      \     /
  \      \   /
   \      middle
    \      /
     \    /
     final
```

Refer to the [demo](https://github.com/buptmiao/parallel/tree/master/examples/demo1/demo.go) below:

```go
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
```

#### eg.3

By default, Parallel will ignore panics of jobs. But parallel supports customized exception handler, which is used for dealing with unexpected panics. For example, alerting or logging.
```go
// handle the panic
func exceptionHandler(topic string, e interface{}) {
	fmt.Println(topic, e)
}

// will panic
func exceptionJob() {
	var a map[string]int
	//assignment to entry in nil map
	a["123"] = 1
}

func main() {
	p := parallel.NewParallel()
	p.Register(exceptionJob)
	// miss the last argument on purpose
	p.Except(exceptionHandler, "topic1")
	p.Run()
}
```
#### [more examples](https://github.com/buptmiao/parallel/tree/master/examples)