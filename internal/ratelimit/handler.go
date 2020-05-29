package ratelimit

import (
	"context"
	"fmt"
	"time"
)
// Notes:

type endpoint interface {
	Call() ([]byte, error)
}
type result struct {
	response []byte
	err      error
}

type request struct {
	epoint endpoint
	respCh chan *result
	ctx    context.Context
}

type Handler struct {
	rate    time.Duration
	reqCh   chan *request
	timeout time.Duration
}

func NewHandler() *Handler {
	rate := time.Second
	//TODO: Look into the standard library to see if this is idiomatic. Unsure about making a channel of interfaces
	requests := make(chan *request, 10)
	reqTimeout := 5 * time.Second
	r := Handler{rate, requests, reqTimeout}

	go r.limiter()

	return &r
}

func (r *Handler) limiter() {
	throttle := time.NewTicker(r.rate)

	for req := range r.reqCh {
		<-throttle.C

		go r.process(req)
	}
}


func (r *Handler) process(req *request) {
	res, err := req.epoint.Call()
	select {
	case <-req.ctx.Done():
		return
	default:
		req.respCh <- &result{res, err}
	}
}

func (r *Handler) Get(req endpoint) ([]byte, error) {
	c := make(chan *result)
	ctx, cancel := context.WithCancel(context.Background())

	defer close(c)

	r.reqCh <- &request{req, c, ctx}
	select {
	case n := <-c:
		cancel()
		return n.response, nil
	case <-time.After(r.timeout):
		cancel()
		//TODO: Refactor error handling
		return nil, fmt.Errorf("timeout waiting for request")
	}
}
