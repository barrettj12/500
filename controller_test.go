package main

import (
	"testing"

	c "github.com/barrettj12/collections"
	"github.com/stretchr/testify/assert"
)

func TestTrickInfoWinnerMisere(t *testing.T) {
	tr := trickInfo{
		leader: 3,
		plays: c.AsList([]PlayInfo{
			{Player: 3, Card: Card{1, Diamonds}},
			{Player: 0, Card: Card{13, Spades}},
			{Player: 1, Card: Card{5, Diamonds}},
		}),
	}
	assert.Equal(t, tr.Winner(MisereBid{}), 3)
}
