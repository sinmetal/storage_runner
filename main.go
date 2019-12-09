package main

import (
	"time"

	"github.com/sinmetal/storage_runner/redis"
)

func main() {
	rc, err := redis.NewClient()
	if err != nil {
		panic(err)
	}
	for {
		err = rc.Set()
		if err != nil {
			panic(err)
		}
		time.Sleep(3 * time.Second)
	}
}
