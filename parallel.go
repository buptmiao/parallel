package parallel

import (
	"sync"
	"time"
)

// Parallel instance, which executes pipelines by parallel
type Parallel struct {
	wg        *sync.WaitGroup
	pipes     []*Pipeline
	wgChild   *sync.WaitGroup
	children  []*Parallel
	exception *Handler
}

// NewParallel creates a new Parallel instance
func NewParallel() *Parallel {
	res := new(Parallel)
	res.wg = new(sync.WaitGroup)
	res.wgChild = new(sync.WaitGroup)
	res.pipes = make([]*Pipeline, 0, 10)
	return res
}

// Except set the exception handling routine, when unexpected panic occur
// this routine will be executed.
func (p *Parallel) Except(f interface{}, args ...interface{}) *Handler {
	h := NewHandler(f, args...)
	p.exception = h
	return h
}

// Register add a new pipeline with a single handler info parallel
func (p *Parallel) Register(f interface{}, args ...interface{}) *Handler {
	return p.NewPipeline().Register(f, args...)
}

// NewPipeline create a new pipeline of parallel
func (p *Parallel) NewPipeline() *Pipeline {
	pipe := NewPipeline()
	p.Add(pipe)
	return pipe
}

// Add add new pipelines to parallel
func (p *Parallel) Add(pipes ...*Pipeline) *Parallel {
	p.wg.Add(len(pipes))
	p.pipes = append(p.pipes, pipes...)
	return p
}

// NewChild create a new child of p
func (p *Parallel) NewChild() *Parallel {
	child := NewParallel()
	child.exception = p.exception
	p.AddChildren(child)
	return child
}

// AddChildren add children to parallel to handle dependency
func (p *Parallel) AddChildren(children ...*Parallel) *Parallel {
	p.wgChild.Add(len(children))
	p.children = append(p.children, children...)
	return p
}

// Run start up all the jobs
func (p *Parallel) Run() {
	for _, child := range p.children {
		// this func will never panic
		go func(ch *Parallel) {
			ch.Run()
			p.wgChild.Done()
		}(child)
	}
	p.wgChild.Wait()
	p.do()
	p.wg.Wait()
}

// Do just do it
func (p *Parallel) do() {
	// if only one pipeline no need go routines
	if len(p.pipes) == 1 {
		p.secure(p.pipes[0])
		return
	}
	for _, pipe := range p.pipes {
		go p.secure(pipe)
	}
}

// exec pipeline safely
func (p *Parallel) secure(pipe *Pipeline) {
	defer func() {
		err := recover()
		if err != nil {
			if err == ErrArgNotFunction || err == ErrInArgLenNotMatch || err == ErrOutArgLenNotMatch || err == ErrRecvArgTypeNotPtr || err == ErrRecvArgNil {
				panic(err)
			}
			if p.exception != nil {
				p.exception.OnExcept(err)
			}
		}
		p.wg.Done()
	}()
	pipe.Do()
}

// RunWithTimeOut start up all the jobs, and time out after d duration
func (p *Parallel) RunWithTimeOut(d time.Duration) {
	success := make(chan struct{}, 1)
	go func() {
		p.Run()
		success <- struct{}{}
	}()
	select {
	case <-success:
	case <-time.After(d):
	}
}
