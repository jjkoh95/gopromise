package main

import (
	"log"
	"time"

	"github.com/jjkoh95/gopromise"
)

type CustomType struct {
	ID          string
	Name        string
	TimeElapsed time.Duration
}

func main() {
	sleep3Promise := gopromise.NewPromise(func(resolve func(CustomType), reject func(error)) {
		res := CustomType{
			ID:   "ID-3",
			Name: "Name-3",
		}
		before := time.Now()
		time.Sleep(3 * time.Second)
		res.TimeElapsed = time.Since(before)
		resolve(res)
	})
	log.Println("start sleep3")

	sleep5Promise := gopromise.NewPromise(func(resolve func(CustomType), reject func(error)) {
		res := CustomType{
			ID:   "ID-5",
			Name: "Name-5",
		}
		before := time.Now()
		time.Sleep(5 * time.Second)
		res.TimeElapsed = time.Since(before)
		resolve(res)
	})
	log.Println("start sleep5")

	sleep3Res := gopromise.Await(sleep3Promise)
	log.Printf("%-10s -- Result: %5v Error: %v\n", "sleep3Res", sleep3Res.Res, sleep3Res.Err)

	sleep5Res := gopromise.Await(sleep5Promise)
	log.Printf("%-10s -- Result: %5v Error: %v\n", "sleep5Res", sleep5Res.Res, sleep5Res.Err)

	// Output
	// 2022/07/30 04:26:31 start sleep3
	// 2022/07/30 04:26:31 start sleep5
	// 2022/07/30 04:26:34 sleep3Res  -- Result: { ID-3 Name-3 3.002143009s} Error: <nil>
	// 2022/07/30 04:26:36 sleep3Res  -- Result: { ID-5 Name-5 5.000642377s} Error: <nil>
}
