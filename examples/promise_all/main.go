package main

import (
	"errors"
	"log"
	"time"

	"github.com/jjkoh95/gopromise"
)

func main() {
	sleep10Promise := gopromise.NewPromise(func(resolve func(string), reject func(error)) {
		time.Sleep(10 * time.Second)
		resolve("done 10 seconds")
	})
	log.Println("start sleep10")

	sleep3Promise := gopromise.NewPromise(func(resolve func(string), reject func(error)) {
		time.Sleep(3 * time.Second)
		resolve("done 3 seconds")
	})
	log.Println("start sleep3")

	sleep5PromiseWithErr := gopromise.NewPromise(func(resolve func(string), reject func(error)) {
		time.Sleep(5 * time.Second)
		reject(errors.New("some error happened"))
	})
	log.Println("start sleep5")

	promises := gopromise.AwaitAll(sleep10Promise, sleep3Promise, sleep5PromiseWithErr)
	log.Println("await all promises done")
	for i, promise := range promises {
		log.Printf("%3d -- Result: %15v Error: %v\n", i, promise.Res, promise.Err)
	}

	// output
	// 2022/07/30 05:43:04 start sleep10
	// 2022/07/30 05:43:04 start sleep3
	// 2022/07/30 05:43:04 start sleep5
	// 2022/07/30 05:43:14 await all promises done
	// 2022/07/30 05:43:14   0 -- Result: done 10 seconds Error: <nil>
	// 2022/07/30 05:43:14   1 -- Result:  done 3 seconds Error: <nil>
	// 2022/07/30 05:43:14   2 -- Result:                 Error: some error happened
}
