package parallel

import (
	"sync"
)

// Parallel instance, which executes pipelines by parallel
type Parallel struct {
	wg       *sync.WaitGroup
	pipes    []*Pipeline
	wgChild  *sync.WaitGroup
	children []*Parallel
}

// NewParallel creates a new Parallel instance
func NewParallel() *Parallel {
	res := new(Parallel)
	res.wg = new(sync.WaitGroup)
	res.wgChild = new(sync.WaitGroup)
	res.pipes = make([]*Pipeline, 0, 10)
	return res
}

// Register add a new pipeline with a single handler info parallel
func (p *Parallel) Register(f interface{}, args ...interface{}) *Handler {
	pipe := NewPipeline()
	h := pipe.Register(f, args...)
	p.Add(pipe)
	return h
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

// exec pipeline secure
func (p *Parallel) secure(pipe *Pipeline) {
	defer func() {
		err := recover()
		if err != nil && (err == ErrArgNotFunction || err == ErrInArgLenNotMatch || err == ErrOutArgLenNotMatch || err == ErrRecvArgTypeNotPtr || err == ErrRecvArgNil) {
			panic(err)
		}
		p.wg.Done()
	}()
	pipe.Do()
}
