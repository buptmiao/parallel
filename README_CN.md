### Parallel

[![Build Status](https://travis-ci.org/buptmiao/parallel.svg?branch=master)](https://travis-ci.org/buptmiao/parallel)
[![Coverage Status](https://coveralls.io/repos/github/buptmiao/parallel/badge.svg?branch=master)](https://coveralls.io/github/buptmiao/parallel?branch=master)

一个go语言并行程序库, 用于不改变现有接口声明前提下的业务聚合或者重构.

### 用法

#### 例子1
testjobA, testjobB, testjobC为三个现有的接口, 并行化它们只需要如下代码:
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


#### 例子2

来看一个稍微复杂的例子,并行执行3个任务jobA, jobB, jobC, final为任务的执行结果. 其中, B和C需要同步结果middle, middle和jobA需要同步最终结果.

```
jobA  jobB   jobC
 \      \     /
  \      \   /
   \      middle
    \      /
     \    /
     final
```

参考如下[demo](https://github.com/buptmiao/parallel/tree/master/example/demo.go):

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
