package main

import (
	"testing"

	c "github.com/barrettj12/collections"
	"github.com/stretchr/testify/assert"
)

func TestTrickInfoWinnerMisere(t *testing.T) {
	tr := trickInfo{
		leader: 3,
		plays: c.AsList([]playInfo{
			{player: 3, card: Card{1, Diamonds}},
			{player: 0, card: Card{13, Spades}},
			{player: 1, card: Card{5, Diamonds}},
		}),
	}
	assert.Equal(t, tr.Winner(MisereBid{}), 3)
}
