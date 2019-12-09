package main

import (
	"math/rand"

	"github.com/google/uuid"
)

func NewNarrowRandom() string {
	id := uuid.New().String()
	cards := []string{"hoge", "fuga", "moge", id}
	return cards[rand.Intn(len(cards))]
}
