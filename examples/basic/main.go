package main

import (
	"log"
	"time"

	"github.com/jjkoh95/gopromise"
)

func main() {
	sleep3Promise := gopromise.NewPromise(func(resolve func(interface{}), reject func(error)) {
		time.Sleep(3 * time.Second)
		resolve("done waiting 3 seconds")
	})
	log.Println("start sleep3")

	sleep10Promise := gopromise.NewPromise(func(resolve func(interface{}), reject func(error)) {
		time.Sleep(10 * time.Second)
		resolve("done waiting 10 seconds")
	})
	log.Println("start sleep10")

	sleep3Res := gopromise.Await(sleep3Promise)
	log.Println("sleep3Res", sleep3Res)

	sleep10Res := gopromise.Await(sleep10Promise)
	log.Println("sleep10Res", sleep10Res)

	// output
	// 2022/01/22 18:35:32 start sleep3
	// 2022/01/22 18:35:32 start sleep10
	// 2022/01/22 18:35:35 sleep3Res {done waiting 3 seconds <nil>}
	// 2022/01/22 18:35:42 sleep10Res {done waiting 10 seconds <nil>}
}
