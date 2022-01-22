package gopromise

import (
	"fmt"
	"sync"
)

type Result struct {
	Res interface{}
	Err error
}

type Promise struct {
	result Result
	ch     chan Result
	// make sure we only await channel once
	awaitOnce sync.Once
	// make sure we only resolve/reject once
	resolveRejectOnce sync.Once
}

func NewPromise(fn func(
	resolve func(interface{}),
	reject func(error),
)) *Promise {
	p := &Promise{
		ch: make(chan Result),
	}

	go func() {
		// handlePanic
		defer func() {
			if r := recover(); r != nil {
				if err, isErr := r.(error); isErr {
					p.reject(err)
					return
				}
				p.reject(fmt.Errorf("error: %v", r))
			}
		}()

		fn(p.resolve, p.reject)
	}()
	return p
}

func (p *Promise) resolve(res interface{}) {
	p.resolveRejectOnce.Do(func() {
		p.ch <- Result{Res: res}
	})
}

func (p *Promise) reject(err error) {
	p.resolveRejectOnce.Do(func() {
		p.ch <- Result{Err: err}
	})
}

func Await(p *Promise) Result {
	p.awaitOnce.Do(func() {
		p.result = <-p.ch
		close(p.ch)
	})
	return p.result
}

func AwaitAll(promises ...*Promise) []Result {
	res := make([]Result, len(promises))
	for i, p := range promises {
		res[i] = Await(p)
	}
	return res
}
