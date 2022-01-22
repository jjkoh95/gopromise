# gopromise

## import
```bash
go get -u github.com/jjkoh95/gopromise
```

## TLDR
```go
// non-blocking
doSomethingPromise := gopromise.NewPromise(func(resolve func(interface{}), reject func(error)) {
    // doSomething or call API
    res, err := doSomething()
    if err != nil {
        reject(err)
        return
    }
    resolve(res)
})

// non-blocking
doSomething2Promise := gopromise.NewPromise(func(resolve func(interface{}), reject func(error)) {
    // doSomething or call API
    res, err := doSomething()
    if err != nil {
        reject(err)
        return
    }
    resolve(res)
})

// blocking
doSomethingRes := gopromise.Await(doSomethingPromise)

// blocking
doSomething2Res := gopromise.Await(doSomething2Promise)

// alternatively you can wait for multiple promises
// blocking
res := gopromise.AwaitAll(doSomethingPromise, doSomething2Promise)
```
