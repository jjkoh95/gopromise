package main

import (
	"log"
	"time"

	"github.com/jjkoh95/gopromise"
)

func main() {
	sleep1Promise := gopromise.NewPromise(func(resolve func(interface{}), reject func(error)) {
		time.Sleep(1 * time.Second)
		resolve("done waiting 1 second")
	})
	log.Println("start sleep1")

	sleep10Promise := gopromise.NewPromise(func(resolve func(interface{}), reject func(error)) {
		time.Sleep(10 * time.Second)
		resolve("done waiting 10 seconds")
	})
	log.Println("start sleep10")

	sleep1Res := gopromise.Await(sleep1Promise)
	log.Println("sleep1Res", sleep1Res)

	sleep10Res := gopromise.Await(sleep10Promise)
	log.Println("sleep10Res", sleep10Res)

	// output
	// 2022/01/22 18:26:23 start sleep1
	// 2022/01/22 18:26:23 start sleep10
	// 2022/01/22 18:26:24 sleep1Res {done waiting 1 second <nil>}
	// 2022/01/22 18:26:33 sleep10Res {done waiting 10 seconds <nil>}
}
