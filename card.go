package main

import "fmt"

// Represents a playing card
type Card struct {
	rank Rank
	suit Suit
}

var (
	JokerCard Card = Card{Joker, NoSuit}
)

func (c Card) String() string {
	return c.rank.String() + c.suit.Symbol(true)
}

func (c Card) PrintGrey() string {
	return grey(c.rank.String() + c.suit.Symbol(false))
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
			return red("♦")
		}
		return "♦"
	case Hearts:
		if colour {
			return red("♥")
		}
		return "♥"
	default:
		return string(s)
	}
}

// Coloured printing

func red(s string) string {
	return fmt.Sprintf("\u001b[31m%s\u001b[0m", s)
}

func grey(s string) string {
	return fmt.Sprintf("\033[38;2;200;200;200m%s\033[0m", s)
}
