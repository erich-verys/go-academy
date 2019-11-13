package main

import "fmt"

func main() {
	// make an instance of Handler with a RealCaller instance
	real := &Handler{
		ExternalCaller: &RealCaller{},
		name:           "real",
	}
	// process real
	real.process()
	// make an instance of Handler with a FakeCaller instance
	fake := &Handler{
		ExternalCaller: &FakeCaller{},
		name:           "fake",
	}
	// process fake
	fake.process()
}

type Handler struct {
	_ struct{}
	ExternalCaller
	name string
}

type RealCaller struct{}
type FakeCaller struct{}

type ExternalCaller interface {
	call()
}

func (h *Handler) process() {
	h.call()
}

func (rc *RealCaller) call() {
	// a real api call would happen here
	fmt.Println("real")
}

func (fc *FakeCaller) call() {
	// our fake api call would happen here
	fmt.Println("fake")
}
