package main

import (
	"testing"
)

func TestNewID(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(NewNarrowRandom())
	}
}
