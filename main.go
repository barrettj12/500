package main

import (
	"math/rand"
	"time"

	c "github.com/barrettj12/collections"
)

const SLEEP = 500 * time.Millisecond

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ct := Controller{
		players: [4]Player{
			&HumanPlayer{},
			&RandomPlayer{delay: SLEEP},
			&RandomPlayer{delay: SLEEP},
			&RandomPlayer{delay: SLEEP},
		},
	}
	ct.Play()
}

// Returns the 500 deck
func getDeck() *c.List[Card] {
	return c.AsList([]Card{
		{4, Diamonds}, {4, Hearts},
		{5, Spades}, {5, Clubs}, {5, Diamonds}, {5, Hearts},
		{6, Spades}, {6, Clubs}, {6, Diamonds}, {6, Hearts},
		{7, Spades}, {7, Clubs}, {7, Diamonds}, {7, Hearts},
		{8, Spades}, {8, Clubs}, {8, Diamonds}, {8, Hearts},
		{9, Spades}, {9, Clubs}, {9, Diamonds}, {9, Hearts},
		{10, Spades}, {10, Clubs}, {10, Diamonds}, {10, Hearts},
		{Jack, Spades}, {Jack, Clubs}, {Jack, Diamonds}, {Jack, Hearts},
		{Queen, Spades}, {Queen, Clubs}, {Queen, Diamonds}, {Queen, Hearts},
		{King, Spades}, {King, Clubs}, {King, Diamonds}, {King, Hearts},
		{Ace, Spades}, {Ace, Clubs}, {Ace, Diamonds}, {Ace, Hearts},
		JokerCard,
	})
}

// Utility functions

// E evaluates the given function and panics if it returns a non-nil error.
func E0(err error) {
	if err != nil {
		panic(err)
	}
}

func E[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
