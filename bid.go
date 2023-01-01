package main

import (
	"fmt"

	c "github.com/barrettj12/collections"
)

type Bid interface {
	Value() int
	Suit(Card) Suit
	CardOrder(leadCard Card) *c.List[Card]
	SortHand(*c.List[Card])
	Won(tricksWon int) bool
}

type SuitBid struct {
	tricks    int
	trumpSuit Suit
}

func (s SuitBid) String() string {
	return fmt.Sprintf("%d%s", s.tricks, s.trumpSuit.Symbol(true))
}

var suitBidValues = map[Suit]int{
	Spades: 40, Clubs: 60, Diamonds: 80, Hearts: 100,
}

func (s SuitBid) Value() int {
	return suitBidValues[s.trumpSuit] + 100*(s.tricks-6)
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
		b.lowBower(),
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

// Sort hand as follows:
//
//	Off-suits (in bidding order): [4] 5 6 7 8 9 10 J Q K A
//	followed by trumps: [4] 5 6 7 8 9 10 Q K A LB J JOK
func (b SuitBid) SortHand(hand *c.List[Card]) {
	hand.Sort(func(c, d Card) bool {
		suitOrder := map[Suit]int{Spades: 1, Clubs: 2, Diamonds: 3, Hearts: 4}
		suitOrder[b.trumpSuit] = 5
		suit1 := suitOrder[b.Suit(c)]
		suit2 := suitOrder[b.Suit(d)]
		if suit1 < suit2 {
			return true
		}
		if suit1 > suit2 {
			return false
		}

		// Same suit - use card order
		cardOrder := b.CardOrder(c)
		i1 := E(cardOrder.Find(c))
		i2 := E(cardOrder.Find(d))
		return i1 > i2
	})
}

func (s SuitBid) Won(tricksWon int) bool {
	return tricksWon >= s.tricks
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
