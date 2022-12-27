package main

import (
	"fmt"

	c "github.com/barrettj12/collections"
)

type Bid interface {
	Value() int
	Suit(Card) Suit
	WhoWins(trick *c.List[Card], leader int) (winner int)
	Won(tricksWon int) bool
}

type SuitBid struct {
	tricks    int
	trumpSuit Suit
}

var suitBidValues = map[Suit]int{
	Spades: 40, Clubs: 60, Diamonds: 80, Hearts: 100, NoSuit: 120,
}

func (s SuitBid) Value() int {
	return suitBidValues[s.trumpSuit] + 100*(s.tricks-6)
}

func (s SuitBid) WhoWins(trick *c.List[Card], leader int) int {
	leadCard := E(trick.Get(leader))
	leadSuit := leadCard.suit
	trumpSuit := s.trumpSuit
	order := cardOrder(leadSuit, trumpSuit)
	winner := leader // whoever lead wins by default

	for _, card := range *order {
		i, err := trick.Find(card)
		if err != nil {
			// not found
			continue
		}

		winner = i
		break
	}

	return winner
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
func cardOrder(lead, trumps Suit) *c.List[Card] {
	order := c.NewList[Card](12)
	order.Append(
		JokerCard,
		Card{Jack, trumps},
		Card{Jack, sameColour[trumps]},
		Card{Ace, trumps},
		Card{King, trumps},
		Card{Queen, trumps},
		Card{10, trumps},
		Card{9, trumps},
		Card{8, trumps},
		Card{7, trumps},
		Card{6, trumps},
		Card{5, trumps},
	)

	trump4 := Card{4, trumps}
	if getDeck().Contains(trump4) {
		order.Append(trump4)
	}

	if lead != trumps {
		// Append cards of lead suit
		order.Append(
			Card{Ace, lead},
			Card{King, lead},
			Card{Queen, lead},
			Card{Jack, lead},
			Card{10, lead},
			Card{9, lead},
			Card{8, lead},
			Card{7, lead},
			Card{6, lead},
			Card{5, lead},
		)

		lead4 := Card{4, lead}
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
	return fmt.Sprintf("%d%s", s.tricks, s.trumpSuit.Symbol())
}

type MisereBid struct {
	open bool
}

func (b MisereBid) Value() int {
	if b.open {
		return 500
	}
	return 250
}
