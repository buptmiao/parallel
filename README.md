### Parallel

[![Build Status](https://travis-ci.org/buptmiao/parallel.svg?branch=master)](https://travis-ci.org/buptmiao/parallel)
[![Coverage Status](https://coveralls.io/repos/github/buptmiao/parallel/badge.svg?branch=master)](https://coveralls.io/github/buptmiao/parallel?branch=master)

[zh_CN](https://github.com/buptmiao/parallel/blob/master/README_CN.md)

A golang parallel library, used for buz aggregation and refactor without changing declaration of function.

### Usage

#### eg.1
There are three method: testjobA, testjobB, testjobC, execute them by parallel:
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

Refer to the [demo](https://github.com/buptmiao/parallel/tree/master/example/demo.go) below:

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
