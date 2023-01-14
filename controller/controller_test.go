package controller

import (
	"testing"

	"github.com/barrettj12/500/card"
	"github.com/barrettj12/500/game"
	c "github.com/barrettj12/collections"
	"github.com/stretchr/testify/assert"
)

func TestTrickInfoWinnerMisere(t *testing.T) {
	tr := trickInfo{
		leader: 3,
		plays: c.AsList([]game.PlayInfo{
			{Player: 3, Card: card.Card{1, card.Diamonds}},
			{Player: 0, Card: card.Card{13, card.Spades}},
			{Player: 1, Card: card.Card{5, card.Diamonds}},
		}),
	}
	assert.Equal(t, tr.Winner(game.MisereBid{}), 3)
}
