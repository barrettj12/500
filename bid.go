package main

import (
	"fmt"

	c "github.com/barrettj12/collections"
)

type Bid interface {
	Value() int
	Suit(Card) Suit
	CardOrder(leadCard Card) *c.List[Card]
	Won(tricksWon int) bool
}

type SuitBid struct {
	tricks    int
	trumpSuit Suit
}

var suitBidValues = map[Suit]int{
	Spades: 40, Clubs: 60, Diamonds: 80, Hearts: 100,
}

func (s SuitBid) Value() int {
	return suitBidValues[s.trumpSuit] + 100*(s.tricks-6)
}

// Returns the card order for the given lead suit and trump suit.
// If lead suit == trump suit (e.g. ♥):
//
//	JOK  J♥  J♦  A♥  K♥  Q♥  10♥  ...
//
// Else (e.g. lead = ♠, trump = ♥):
//
//	JOK  J♥  J♦  A♥  K♥  Q♥  10♥  ...
//	A♠  K♠  Q♠  J♠  10♠  ...
func (b SuitBid) CardOrder(leadCard Card) *c.List[Card] {
	leadSuit := b.Suit(leadCard)

	order := c.NewList[Card](12)
	order.Append(
		JokerCard,
		Card{Jack, b.trumpSuit},
		Card{Jack, sameColour[b.trumpSuit]},
		Card{Ace, b.trumpSuit},
		Card{King, b.trumpSuit},
		Card{Queen, b.trumpSuit},
		Card{10, b.trumpSuit},
		Card{9, b.trumpSuit},
		Card{8, b.trumpSuit},
		Card{7, b.trumpSuit},
		Card{6, b.trumpSuit},
		Card{5, b.trumpSuit},
	)

	trump4 := Card{4, b.trumpSuit}
	if getDeck().Contains(trump4) {
		order.Append(trump4)
	}

	if leadSuit != b.trumpSuit {
		// Append cards of lead suit
		order.Append(
			Card{Ace, leadSuit},
			Card{King, leadSuit},
			Card{Queen, leadSuit},
			Card{Jack, leadSuit},
			Card{10, leadSuit},
			Card{9, leadSuit},
			Card{8, leadSuit},
			Card{7, leadSuit},
			Card{6, leadSuit},
			Card{5, leadSuit},
		)

		lead4 := Card{4, leadSuit}
		if getDeck().Contains(lead4) {
			order.Append(lead4)
		}
	}

	return order
}

// Which suit has same colour?
var sameColour = map[Suit]Suit{
	Spades: Clubs, Clubs: Spades,
	Diamonds: Hearts, Hearts: Diamonds,
}

func (s SuitBid) lowBower() Card {
	return Card{Jack, sameColour[s.trumpSuit]}
}

// Returns suit for the given card in this bid.
// This is generally the same suit except for Joker and low bower.
func (s SuitBid) Suit(c Card) Suit {
	if c == JokerCard || c == s.lowBower() {
		return s.trumpSuit
	}
	return c.suit
}

func (s SuitBid) Won(tricksWon int) bool {
	return tricksWon >= s.tricks
}

func (s SuitBid) String() string {
	return fmt.Sprintf("%d%s", s.tricks, s.trumpSuit.Symbol(true))
}

type NoTrumpsBid struct {
	tricks int
}

func (b NoTrumpsBid) Value() int {
	return 120 + 100*(b.tricks-6)
}

func (b NoTrumpsBid) String() string {
	return fmt.Sprintf("%dNT", b.tricks)
}

type MisereBid struct {
	NoTrumpsBid
	open bool
}

func (b MisereBid) Value() int {
	if b.open {
		return 500
	}
	return 250
}

func (b MisereBid) String() string {
	if b.open {
		return "OpMis"
	}
	return "Mis"
}
