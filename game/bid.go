package game

import (
	"fmt"

	"github.com/barrettj12/500/card"
	"github.com/barrettj12/500/util"

	c "github.com/barrettj12/collections"
)

type Bid interface {
	Value() int
	Suit(card.Card) card.Suit
	CardOrder(leadCard card.Card) *c.List[card.Card]
	ValidPlays(trick *c.List[PlayInfo], hand *c.List[card.Card]) *c.List[int]
	SortHand(*c.List[card.Card])
	Won(tricksWon int) bool
}

type SuitBid struct {
	Tricks    int
	TrumpSuit card.Suit
}

func (s SuitBid) String() string {
	return fmt.Sprintf("%d%s", s.Tricks, s.TrumpSuit.Symbol(true))
}

var suitBidValues = map[card.Suit]int{
	card.Spades: 40, card.Clubs: 60, card.Diamonds: 80, card.Hearts: 100,
}

func (s SuitBid) Value() int {
	return suitBidValues[s.TrumpSuit] + 100*(s.Tricks-6)
}

// Which suit has same colour?
var sameColour = map[card.Suit]card.Suit{
	card.Spades: card.Clubs, card.Clubs: card.Spades,
	card.Diamonds: card.Hearts, card.Hearts: card.Diamonds,
}

func (s SuitBid) lowBower() card.Card {
	return card.Card{card.Jack, sameColour[s.TrumpSuit]}
}

// Returns suit for the given card in this bid.
// This is generally the same suit except for Joker and low bower.
func (s SuitBid) Suit(c card.Card) card.Suit {
	if c == card.JokerCard || c == s.lowBower() {
		return s.TrumpSuit
	}
	return c.Suit
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
func (b SuitBid) CardOrder(leadCard card.Card) *c.List[card.Card] {
	leadSuit := b.Suit(leadCard)

	order := c.NewList[card.Card](12)
	order.Append(
		card.JokerCard,
		card.Card{card.Jack, b.TrumpSuit},
		b.lowBower(),
		card.Card{card.Ace, b.TrumpSuit},
		card.Card{card.King, b.TrumpSuit},
		card.Card{card.Queen, b.TrumpSuit},
		card.Card{10, b.TrumpSuit},
		card.Card{9, b.TrumpSuit},
		card.Card{8, b.TrumpSuit},
		card.Card{7, b.TrumpSuit},
		card.Card{6, b.TrumpSuit},
		card.Card{5, b.TrumpSuit},
	)

	trump4 := card.Card{4, b.TrumpSuit}
	if GetDeck().Contains(trump4) {
		order.Append(trump4)
	}

	if leadSuit != b.TrumpSuit {
		// Append cards of lead suit
		order.Append(
			card.Card{card.Ace, leadSuit},
			card.Card{card.King, leadSuit},
			card.Card{card.Queen, leadSuit},
			card.Card{card.Jack, leadSuit},
			card.Card{10, leadSuit},
			card.Card{9, leadSuit},
			card.Card{8, leadSuit},
			card.Card{7, leadSuit},
			card.Card{6, leadSuit},
			card.Card{5, leadSuit},
		)

		lead4 := card.Card{4, leadSuit}
		if GetDeck().Contains(lead4) {
			order.Append(lead4)
		}
	}

	return order
}

// Returns indices of valid plays in hand.
func (b SuitBid) ValidPlays(trick *c.List[PlayInfo], hand *c.List[card.Card]) *c.List[int] {
	valids := c.NewList[int](hand.Size())

	for i, c := range *hand {
		if trick.Size() == 0 {
			// Can lead with any card
			valids.Append(i)
			continue
		}

		// We have to follow suit if we can
		leadCard := util.E(trick.Get(0)).Card
		leadSuit := b.Suit(leadCard)
		if b.Suit(c) == leadSuit {
			valids.Append(i)
			continue
		}

		// Check if we can't follow suit: then we can play anything
		numOfLeadSuit := hand.Count(func(_ int, c card.Card) bool {
			return b.Suit(c) == leadSuit
		})

		if numOfLeadSuit == 0 {
			valids.Append(i)
			continue
		}
	}

	return valids
}

// Sort hand as follows:
//
//	Off-suits (in bidding order): [4] 5 6 7 8 9 10 J Q K A
//	followed by trumps: [4] 5 6 7 8 9 10 Q K A LB J JOK
func (b SuitBid) SortHand(hand *c.List[card.Card]) {
	hand.Sort(func(c, d card.Card) bool {
		suit1 := b.suitOrder(c)
		suit2 := b.suitOrder(d)
		if suit1 < suit2 {
			return true
		}
		if suit1 > suit2 {
			return false
		}

		// Same suit - use card order
		cardOrder := b.CardOrder(c)
		i1 := util.E(cardOrder.Find(c))
		i2 := util.E(cardOrder.Find(d))
		return i1 > i2
	})
}

// Returns a number determining the order of suits in the hand.
func (b SuitBid) suitOrder(c card.Card) int {
	return map[card.Suit]map[card.Suit]int{
		card.Spades:   {card.Diamonds: 1, card.Clubs: 2, card.Hearts: 3, card.Spades: 4},
		card.Clubs:    {card.Diamonds: 1, card.Spades: 2, card.Hearts: 3, card.Clubs: 4},
		card.Diamonds: {card.Spades: 1, card.Hearts: 2, card.Clubs: 3, card.Diamonds: 4},
		card.Hearts:   {card.Spades: 1, card.Diamonds: 2, card.Clubs: 3, card.Hearts: 4},
	}[b.TrumpSuit][b.Suit(c)]
}

func (s SuitBid) Won(tricksWon int) bool {
	return tricksWon >= s.Tricks
}

type NoTrumpsBid struct {
	Tricks int
	// Keep track of the Joker suit (if led) so others follow suit
	JokerSuit card.Suit
}

func (b NoTrumpsBid) String() string {
	return fmt.Sprintf("%dNT", b.Tricks)
}

func (b NoTrumpsBid) Value() int {
	return 120 + 100*(b.Tricks-6)
}

func (b NoTrumpsBid) Suit(c card.Card) card.Suit {
	if b.JokerSuit != "" && c == card.JokerCard {
		return b.JokerSuit
	}
	return c.Suit
}

func (b NoTrumpsBid) CardOrder(leadCard card.Card) *c.List[card.Card] {
	leadSuit := b.Suit(leadCard)
	order := c.NewList[card.Card](12)
	order.Append(
		card.JokerCard,
		card.Card{card.Ace, leadSuit},
		card.Card{card.King, leadSuit},
		card.Card{card.Queen, leadSuit},
		card.Card{card.Jack, leadSuit},
		card.Card{10, leadSuit},
		card.Card{9, leadSuit},
		card.Card{8, leadSuit},
		card.Card{7, leadSuit},
		card.Card{6, leadSuit},
		card.Card{5, leadSuit},
	)

	lead4 := card.Card{4, leadSuit}
	if GetDeck().Contains(lead4) {
		order.Append(lead4)
	}

	return order
}

func (b NoTrumpsBid) ValidPlays(trick *c.List[PlayInfo], hand *c.List[card.Card]) *c.List[int] {
	valids := c.NewList[int](hand.Size())

	for i, c := range *hand {
		// TODO: fix Joker rules
		if c == card.JokerCard {
			valids.Append(i)
			continue
		}

		if trick.Size() == 0 {
			// Can lead with any card
			valids.Append(i)
			continue
		}

		// We have to follow suit if we can
		leadCard := util.E(trick.Get(0)).Card
		leadSuit := b.Suit(leadCard)
		if b.Suit(c) == leadSuit {
			valids.Append(i)
			continue
		}

		// Check if we can't follow suit: then we can play anything
		numOfLeadSuit := hand.Count(func(_ int, c card.Card) bool {
			return b.Suit(c) == leadSuit
		})

		if numOfLeadSuit == 0 {
			valids.Append(i)
			continue
		}
	}

	return valids
}

func (b NoTrumpsBid) SortHand(hand *c.List[card.Card]) {
	hand.Sort(func(c, d card.Card) bool {
		if c.Rank == card.Joker {
			return false
		}
		if d.Rank == card.Joker {
			return true
		}

		suitOrder := map[card.Suit]int{card.Spades: 1, card.Diamonds: 2, card.Clubs: 3, card.Hearts: 4}
		suit1 := suitOrder[c.Suit]
		suit2 := suitOrder[d.Suit]
		if suit1 < suit2 {
			return true
		}
		if suit1 > suit2 {
			return false
		}

		// Same suit - aces high
		if c.Rank == card.Ace {
			return false
		}
		if d.Rank == card.Ace {
			return true
		}
		return c.Rank < d.Rank
	})
}

func (b NoTrumpsBid) Won(tricksWon int) bool {
	return tricksWon >= b.Tricks
}

type MisereBid struct {
	NoTrumpsBid
	Open bool
}

func (b MisereBid) String() string {
	if b.Open {
		return "Open Misère"
	}
	return "Misère"
}

func (b MisereBid) Value() int {
	if b.Open {
		return 500
	}
	return 250
}

func (b MisereBid) Won(tricksWon int) bool {
	return tricksWon == 0
}

// Pass is a special Bid used by the controller in the bidding round
// to represent a Player passing.
type Pass struct{}

// None of these functions should be called - they just ensure that Pass is a
// Bid.
func (p Pass) Value() int                                      { panic("Pass.Value unimplemented") }
func (p Pass) Suit(card.Card) card.Suit                        { panic("Pass.Suit unimplemented") }
func (p Pass) CardOrder(leadCard card.Card) *c.List[card.Card] { panic("Pass.CardOrder unimplemented") }
func (p Pass) ValidPlays(trick *c.List[PlayInfo], hand *c.List[card.Card]) *c.List[int] {
	panic("Pass.ValidPlays unimplemented")
}
func (p Pass) SortHand(*c.List[card.Card]) { panic("Pass.SortHand unimplemented") }
func (p Pass) Won(tricksWon int) bool      { panic("Pass.Won unimplemented") }
