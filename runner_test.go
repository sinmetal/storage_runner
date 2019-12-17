package main

import (
	"testing"

	"github.com/sinmetal/storage_runner/redis"
)

func TestGoSetRedis(t *testing.T) {
	rc, err := redis.NewClient("127.0.0.1:6379")
	if err != nil {
		t.Fatal(err)
	}

	endCh := make(chan error, 10)

	GoSetRedis(rc, 1, endCh)
	GoGetRedis(rc, 1, endCh)

	err = <-endCh
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeadline(t *testing.T) {

}
