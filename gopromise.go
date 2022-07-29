package gopromise

import (
	"fmt"
	"sync"
)

// Result is the main struct that holds promise result
type Result[T any] struct {
	Res T
	Err error
}

// Promise is a promise wrapper to Result
type Promise[T any] struct {
	result    Result[T]
	ch        chan Result[T]
	resolveFn func(T)
	rejectFn  func(error)
	// make sure we only await channel once
	awaitOnce sync.Once
	// make sure we only resolve/reject once
	resolveRejectOnce sync.Once
}

// NewPromise wraps a function and executes it in the background
func NewPromise[T any](fn func(
	resolve func(T),
	reject func(error),
)) *Promise[T] {
	p := &Promise[T]{
		ch: make(chan Result[T]),
	}
	p.resolveFn = func(t T) {
		p.resolveRejectOnce.Do(func() {
			p.ch <- Result[T]{Res: t}
			close(p.ch)
		})
	}
	p.rejectFn = func(err error) {
		p.resolveRejectOnce.Do(func() {
			p.ch <- Result[T]{Err: err}
			close(p.ch)
		})
	}

	go func() {
		defer func() {
			// handlePanic
			if r := recover(); r != nil {
				if err, isErr := r.(error); isErr {
					p.rejectFn(err)
					return
				}
				p.rejectFn(fmt.Errorf("error: %v", r))
			}

			// make sure all promises are resolved
			Await(p)
		}()

		// execute function
		fn(p.resolveFn, p.rejectFn)
	}()

	return p
}

// Await blocks until promise is either resolved or rejected
func Await[T any](p *Promise[T]) Result[T] {
	p.awaitOnce.Do(func() {
		p.result = <-p.ch
	})
	return p.result
}

// AwaitAll blocks until all promises are resolved or rejected
// it returns results in the order of promises
func AwaitAll[T any](promises ...*Promise[T]) []Result[T] {
	res := make([]Result[T], len(promises))
	for i, p := range promises {
		res[i] = Await(p)
	}
	return res
}

// Any blocks until whichever promise is resolved/rejected
// and returns its result
func Any[T any](promises ...*Promise[T]) Result[T] {
	agg := make(chan Result[T])
	for _, promise := range promises {
		go func(p *Promise[T]) {
			val := Await(p)
			agg <- val
		}(promise)
	}

	return <-agg
}
