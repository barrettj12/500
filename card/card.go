package card

import (
	"fmt"

	"github.com/barrettj12/500/util"
)

// Represents a playing card
type Card struct {
	Rank Rank
	Suit Suit
}

var (
	JokerCard = Card{Joker, NoSuit}
)

func (c Card) String() string {
	return c.Rank.String() + c.Suit.Symbol(true)
}

func (c Card) PrintGrey() string {
	return util.Grey(c.Rank.String() + c.Suit.Symbol(false))
}

type Rank int

// Ranks 2-10 are just the corresponding int
const (
	Ace   Rank = 1
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Joker Rank = 14
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Joker:
		return "JOK"
	default:
		return fmt.Sprint(int(r))
	}
}

type Suit string

const (
	Spades   Suit = "Spades"
	Clubs    Suit = "Clubs"
	Diamonds Suit = "Diamonds"
	Hearts   Suit = "Hearts"
	NoSuit   Suit = ""
)

func (s Suit) Symbol(colour bool) string {
	switch s {
	case Spades:
		return "♠"
	case Clubs:
		return "♣"
	case Diamonds:
		if colour {
			return util.Red("♦")
		}
		return "♦"
	case Hearts:
		if colour {
			return util.Red("♥")
		}
		return "♥"
	default:
		return string(s)
	}
}
