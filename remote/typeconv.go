package remote

import (
	"fmt"

	main "github.com/barrettj12/500"

	c "github.com/barrettj12/collections"
)

// encodeHand converts a *c.List[main.Card] to a *Hand.
func encodeHand(list *c.List[main.Card]) *Hand {
	cards := make([]*Card, 0, list.Size())
	for _, c := range *list {
		cards = append(cards, encodeCard(c))
	}
	return &Hand{
		Hand: cards,
	}
}

// encodeCard converts a main.Card to a *Card.
func encodeCard(c main.Card) *Card {
	return &Card{
		Rank: encodeRank(c.Rank),
	}
}

// encodeRank converts a main.Rank to a Rank.
func encodeRank(r main.Rank) Rank {
	return Rank(r)
}

// encodeSuit converts a main.Suit to a Suit.
func encodeSuit(s main.Suit) Suit {
	switch s {
	case main.Spades:
		return Suit_SPADES
	case main.Clubs:
		return Suit_CLUBS
	case main.Diamonds:
		return Suit_DIAMONDS
	case main.Hearts:
		return Suit_HEARTS
	case main.NoSuit:
		return Suit_NO_SUIT
	default:
		panic(fmt.Sprintf("unknown suit %q", s))
	}
}

// encodeTrick converts a *c.List[main.PlayInfo] to a []*PlayInfo.
func encodeTrick(list *c.List[main.PlayInfo]) []*PlayInfo {
	out := make([]*PlayInfo, 0, list.Size())
	for _, pi := range *list {
		out = append(out, encodePlayInfo(pi))
	}
	return out
}

// encodePlayInfo converts a main.PlayInfo to a *PlayInfo.
func encodePlayInfo(pi main.PlayInfo) *PlayInfo {
	return &PlayInfo{
		Player: int32(pi.Player),
		Card:   encodeCard(pi.Card),
	}
}

// encodeValidPlays converts a *c.List[int] to a []int32.
func encodeValidPlays(list *c.List[int]) []int32 {
	out := make([]int32, 0, list.Size())
	for _, n := range *list {
		out = append(out, int32(n))
	}
	return out
}
