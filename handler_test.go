package parallel_test

import (
	"fmt"
	"github.com/buptmiao/parallel"
	"testing"
)

func testFunc1(x, y int) int {
	return x + y
}

func EatPanic(s chan struct{}) {
	if r := recover(); r != nil {
		s <- struct{}{}
	}
}

func TestNewHandler(t *testing.T) {
	h := parallel.NewHandler(testFunc1, 1, 1)
	result := 0
	h.SetReceivers(&result)
	h.Do()
	if result != 2 {
		panic(fmt.Errorf("expect 2, but %s", result))
	}
}

func testFunc10(a *int) {

}

func TestCheckPanic(t *testing.T) {
	// test function type
	s := make(chan struct{}, 1)
	func() {
		defer EatPanic(s)
		h := parallel.NewHandler(0, 1, 1)
		h.Do()
	}()
	<-s

	// test argument length
	func() {
		defer EatPanic(s)
		h := parallel.NewHandler(testFunc1, 1)
		h.Do()
	}()
	<-s

	// test return value length
	func() {
		defer EatPanic(s)
		h := parallel.NewHandler(testFunc1, 1, 1)
		h.Do()
	}()
	<-s

	// test return value type
	result := 0
	func() {
		defer EatPanic(s)
		h := parallel.NewHandler(testFunc1, 1, 1).SetReceivers(result)
		h.Do()
	}()
	<-s

	// test return value nil
	var res *int
	func() {
		defer EatPanic(s)
		h := parallel.NewHandler(testFunc1, 1, 1).SetReceivers(res)
		h.Do()
	}()
	<-s

	// test arguments nil
	var res1 interface{}
	func() {
		h := parallel.NewHandler(testFunc10, res1)
		h.Do()
	}()
}
