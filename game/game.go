package game

import (
	"fmt"

	"github.com/barrettj12/500/card"
	c "github.com/barrettj12/collections"
)

type PlayInfo struct {
	Player int
	Card   card.Card
}

// Returns the 500 deck
func GetDeck() *c.List[card.Card] {
	return c.AsList([]card.Card{
		{4, card.Diamonds}, {4, card.Hearts},
		{5, card.Spades}, {5, card.Clubs}, {5, card.Diamonds}, {5, card.Hearts},
		{6, card.Spades}, {6, card.Clubs}, {6, card.Diamonds}, {6, card.Hearts},
		{7, card.Spades}, {7, card.Clubs}, {7, card.Diamonds}, {7, card.Hearts},
		{8, card.Spades}, {8, card.Clubs}, {8, card.Diamonds}, {8, card.Hearts},
		{9, card.Spades}, {9, card.Clubs}, {9, card.Diamonds}, {9, card.Hearts},
		{10, card.Spades}, {10, card.Clubs}, {10, card.Diamonds}, {10, card.Hearts},
		{card.Jack, card.Spades}, {card.Jack, card.Clubs}, {card.Jack, card.Diamonds}, {card.Jack, card.Hearts},
		{card.Queen, card.Spades}, {card.Queen, card.Clubs}, {card.Queen, card.Diamonds}, {card.Queen, card.Hearts},
		{card.King, card.Spades}, {card.King, card.Clubs}, {card.King, card.Diamonds}, {card.King, card.Hearts},
		{card.Ace, card.Spades}, {card.Ace, card.Clubs}, {card.Ace, card.Diamonds}, {card.Ace, card.Hearts},
		card.JokerCard,
	})
}

// HandResult represents the outcome of a hand.
type HandResult interface {
	Info() string
}

// Redeal is a HandResult representing all players passing during bidding.
type Redeal struct{}

func (r Redeal) Info() string {
	return "Re-deal due to all players passing"
}

// BidWon says that the contractor won their bid.
type BidWon struct {
	Bid    Bid
	Tricks int
}

func (r BidWon) Info() string {
	return fmt.Sprintf("Contractors won their bid of %s with %d tricks",
		r.Bid, r.Tricks)
}

// BidLost says that the contractors lost their bid.
type BidLost struct {
	Bid    Bid
	Tricks int
}

func (r BidLost) Info() string {
	return fmt.Sprintf("Contractors lost their bid of %s with %d tricks",
		r.Bid, r.Tricks)
}
