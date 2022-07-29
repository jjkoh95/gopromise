package main

import (
	"errors"
	"log"
	"time"

	"github.com/jjkoh95/gopromise"
)

func main() {
	sleep3Promise := gopromise.NewPromise(func(resolve func(string), reject func(error)) {
		time.Sleep(3 * time.Second)
		resolve("done 3 seconds")
	})
	log.Println("start sleep3")

	sleep10Promise := gopromise.NewPromise(func(resolve func(string), reject func(error)) {
		time.Sleep(10 * time.Second)
		resolve("done 10 seconds")
	})
	log.Println("start sleep10")

	sleep5PromiseWithErr := gopromise.NewPromise(func(resolve func(string), reject func(error)) {
		time.Sleep(5 * time.Second)
		reject(errors.New("some error happened"))
	})
	log.Println("start sleep5")

	sleep3Res := gopromise.Await(sleep3Promise)
	log.Printf("%-10s -- Result: %-15v Error: %v\n", "sleep3Res", sleep3Res.Res, sleep3Res.Err)

	sleep5Res := gopromise.Await(sleep5PromiseWithErr)
	log.Printf("%-10s -- Result: %-15v Error: %v\n", "sleep5Res", sleep5Res.Res, sleep5Res.Err)

	sleep10Res := gopromise.Await(sleep10Promise)
	log.Printf("%-10s -- Result: %-15v Error: %v\n", "sleep10Res", sleep10Res.Res, sleep10Res.Err)

	// output
	// 2022/07/30 04:21:57 start sleep3
	// 2022/07/30 04:21:57 start sleep10
	// 2022/07/30 04:21:57 start sleep5
	// 2022/07/30 04:22:00 sleep3Res  -- Result: done 3 seconds  Error: <nil>
	// 2022/07/30 04:22:02 sleep5Res  -- Result:                 Error: some error happened
	// 2022/07/30 04:22:07 sleep10Res -- Result: done 10 seconds Error: <nil>
}
