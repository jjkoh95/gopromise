package gopromise_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jjkoh95/gopromise"
	"github.com/stretchr/testify/require"
)

func TestNewPromise(t *testing.T) {
	tests := []struct {
		name string
		fn   func(resolve func(interface{}), reject func(error))
	}{
		{
			name: "promise with resolve",
			fn: func(resolve func(interface{}), reject func(error)) {
				resolve(true)
			},
		},
		{
			name: "promise with reject",
			fn: func(resolve func(interface{}), reject func(error)) {
				reject(errors.New("invalid"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := gopromise.NewPromise(tc.fn)
			require.NotNil(t, p)
		})
	}
}

func TestAwait(t *testing.T) {
	var toPointer = func(x interface{}) interface{} {
		return &x
	}

	tests := []struct {
		name   string
		fn     func(resolve func(interface{}), reject func(error))
		expect gopromise.Result[any]
	}{
		{
			name: "resolve with true",
			fn: func(resolve func(interface{}), reject func(error)) {
				resolve(true)
			},
			expect: gopromise.Result[any]{Res: true},
		},
		{
			name: "resolve with nil",
			fn: func(resolve func(interface{}), reject func(error)) {
				resolve(nil)
			},
			expect: gopromise.Result[any]{Res: nil},
		},
		{
			name: "resolve with pointer",
			fn: func(resolve func(interface{}), reject func(error)) {
				resolve(toPointer("something"))
			},
			expect: gopromise.Result[any]{Res: toPointer("something")},
		},
		{
			name: "reject with error",
			fn: func(resolve func(interface{}), reject func(error)) {
				reject(errors.New("err"))
			},
			expect: gopromise.Result[any]{Err: errors.New("err")},
		},
		{
			name: "reject with nil",
			fn: func(resolve func(interface{}), reject func(error)) {
				reject(nil)
			},
			expect: gopromise.Result[any]{},
		},
		{
			name: "panic with err",
			fn: func(resolve func(interface{}), reject func(error)) {
				panic(errors.New("panic"))
			},
			expect: gopromise.Result[any]{Err: errors.New("panic")},
		},
		{
			name: "panic with non-error",
			fn: func(resolve func(interface{}), reject func(error)) {
				panic(100)
			},
			expect: gopromise.Result[any]{Err: errors.New("error: 100")},
		},
		{
			name: "multiple resolve",
			fn: func(resolve func(interface{}), reject func(error)) {
				resolve(1)
				resolve(2)
				resolve(3)
			},
			expect: gopromise.Result[any]{Res: 1},
		},
		{
			name: "multiple reject",
			fn: func(resolve func(interface{}), reject func(error)) {
				reject(errors.New("error1"))
				reject(errors.New("error2"))
				reject(errors.New("error3"))
			},
			expect: gopromise.Result[any]{Err: errors.New("error1")},
		},
		{
			name: "mix of resolve and reject, with resolve first",
			fn: func(resolve func(interface{}), reject func(error)) {
				resolve(1)
				reject(errors.New("error1"))
			},
			expect: gopromise.Result[any]{Res: 1},
		},
		{
			name: "mix of resolve and reject, with reject first",
			fn: func(resolve func(interface{}), reject func(error)) {
				reject(errors.New("error1"))
				resolve(1)
			},
			expect: gopromise.Result[any]{Err: errors.New("error1")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := gopromise.NewPromise(tc.fn)
			res := gopromise.Await(p)
			require.Equal(t, tc.expect, res)
		})
	}
}

func TestAwaitAll(t *testing.T) {
	tests := []struct {
		name   string
		inputs []*gopromise.Promise[any]
		expect []gopromise.Result[any]
	}{
		{
			name:   "empty",
			inputs: []*gopromise.Promise[any]{},
			expect: []gopromise.Result[any]{},
		},
		{
			name: "single resolve",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						resolve(true)
					},
				),
			},
			expect: []gopromise.Result[any]{{Res: true}},
		},
		{
			name: "single reject",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						reject(errors.New("err"))
					},
				),
			},
			expect: []gopromise.Result[any]{{Err: errors.New("err")}},
		},
		{
			name: "multiple resolve",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						resolve(true)
					},
				),
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						resolve(10)
					},
				),
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						resolve(nil)
					},
				),
			},
			expect: []gopromise.Result[any]{{Res: true}, {Res: 10}, {Res: nil}},
		},
		{
			name: "multiple reject",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						reject(errors.New("error1"))
					},
				),
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						reject(errors.New("error2"))
					},
				),
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						reject(errors.New("error3"))
					},
				),
			},
			expect: []gopromise.Result[any]{{Err: errors.New("error1")}, {Err: errors.New("error2")}, {Err: errors.New("error3")}},
		},
		{
			name: "mix of resolve and reject",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						reject(errors.New("error1"))
					},
				),
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						resolve(2)
					},
				),
				gopromise.NewPromise(
					func(resolve func(interface{}), reject func(error)) {
						reject(errors.New("error3"))
					},
				),
			},
			expect: []gopromise.Result[any]{{Err: errors.New("error1")}, {Res: 2}, {Err: errors.New("error3")}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := gopromise.AwaitAll(tc.inputs...)
			require.Equal(t, tc.expect, res)
		})
	}
}

func TestAwaitAny(t *testing.T) {
	tests := []struct {
		name   string
		inputs []*gopromise.Promise[any]
		expect gopromise.Result[any]
	}{
		{
			name: "fastest-resolve",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(func(resolve func(any), reject func(error)) {
					<-time.After(200 * time.Millisecond)
					resolve(200)
				}),
				gopromise.NewPromise(func(resolve func(any), reject func(error)) {
					<-time.After(50 * time.Millisecond)
					resolve(50)
				}),
				gopromise.NewPromise(func(resolve func(any), reject func(error)) {
					<-time.After(150 * time.Millisecond)
					resolve(150)
				}),
			},
			expect: gopromise.Result[any]{Res: 50},
		},
		{
			name: "fastest-reject",
			inputs: []*gopromise.Promise[any]{
				gopromise.NewPromise(func(resolve func(any), reject func(error)) {
					<-time.After(200 * time.Millisecond)
					resolve(200)
				}),
				gopromise.NewPromise(func(resolve func(any), reject func(error)) {
					<-time.After(50 * time.Millisecond)
					reject(errors.New("50"))
				}),
				gopromise.NewPromise(func(resolve func(any), reject func(error)) {
					<-time.After(150 * time.Millisecond)
					resolve(150)
				}),
			},
			expect: gopromise.Result[any]{Err: errors.New("50")},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := gopromise.Any(tc.inputs...)
			require.Equal(t, tc.expect, res)
		})
	}
}
