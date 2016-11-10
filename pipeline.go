package parallel

// Pipeline instance, which executes jobs by serial
type Pipeline struct {
	handlers []*Handler
}

// NewPipeline creates a new Pipeline instance
func NewPipeline() *Pipeline {
	res := new(Pipeline)
	return res
}

// Register add a new function to pipeline
func (p *Pipeline) Register(f interface{}, args ...interface{}) *Handler {
	h := NewHandler(f, args...)
	p.Add(h)
	return h
}

// Add add new handlers to pipeline
func (p *Pipeline) Add(hs ...*Handler) *Pipeline {
	p.handlers = append(p.handlers, hs...)
	return p
}

// Do calls all handlers as the sequence they are added into pipeline.
func (p *Pipeline) Do() {
	for _, h := range p.handlers {
		h.Do()
	}
}
