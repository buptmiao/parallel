package parallel

import (
	"errors"
	"reflect"
)

var (
	ErrArgNotFunction    = errors.New("argument type not function")            //argument type not function
	ErrInArgLenNotMatch  = errors.New("input arguments length not match")      //input arguments length not match
	ErrOutArgLenNotMatch = errors.New("output arguments length not match")     //output arguments length not match
	ErrRecvArgTypeNotPtr = errors.New("receiver argument type is not pointer") //receiver argument type is not pointer
	ErrRecvArgNil        = errors.New("receiver argument must not be nil")     //receiver argument must not be nil
)

// Handler instance
type Handler struct {
	// The type of f must be function
	f    interface{}
	args []interface{}
	// The type of every receiver must be ptr, to receive the return value of f call
	receivers []interface{}
}

// NewHandler create a new Handler which contains a single function call
func NewHandler(f interface{}, args ...interface{}) *Handler {
	res := new(Handler)
	res.f = f
	res.args = args
	return res
}

// SetReceivers sets the receivers of return values
func (h *Handler) SetReceivers(receivers ...interface{}) *Handler {
	h.receivers = receivers
	return h
}

// Do call the function and return values if exists
func (h *Handler) Do() {
	f := reflect.ValueOf(h.f)
	typ := f.Type()
	//check if f is a function
	if typ.Kind() != reflect.Func {
		panic(ErrArgNotFunction)
	}
	//check input length, only check '>' is to allow varargs.
	if typ.NumIn() > len(h.args) {
		panic(ErrInArgLenNotMatch)
	}
	//check output length
	if typ.NumOut() != len(h.receivers) {
		panic(ErrOutArgLenNotMatch)
	}
	//check if output args is ptr
	for _, v := range h.receivers {
		t := reflect.ValueOf(v)
		if t.Type().Kind() != reflect.Ptr {
			panic(ErrRecvArgTypeNotPtr)
		}
		if t.IsNil() {
			panic(ErrRecvArgNil)
		}
	}

	inputs := make([]reflect.Value, len(h.args))
	for i := 0; i < len(h.args); i++ {
		if h.args[i] == nil {
			inputs[i] = reflect.Zero(f.Type().In(i))
		} else {
			inputs[i] = reflect.ValueOf(h.args[i])
		}
	}
	out := f.Call(inputs)

	for i := 0; i < len(h.receivers); i++ {
		v := reflect.ValueOf(h.receivers[i])
		v.Elem().Set(out[i])
	}
}

// OnExcept will executed by parallel when application panic occur
// Note that the type of e is unknown.
func (h *Handler) OnExcept(e interface{}) {
	h.args = append(h.args, e)
	h.Do()
}
