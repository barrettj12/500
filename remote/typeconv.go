package remote

import (
	"fmt"

	"github.com/barrettj12/500/card"
	"github.com/barrettj12/500/game"

	c "github.com/barrettj12/collections"
)

// encodeHand converts a *c.List[main.Card] to a *Hand.
func encodeHand(list *c.List[card.Card]) *Hand {
	cards := make([]*Card, 0, list.Size())
	for _, c := range *list {
		cards = append(cards, encodeCard(c))
	}
	return &Hand{
		Hand: cards,
	}
}

// decodeHand converts a *Hand to a *c.List[main.Card].
func decodeHand(hand *Hand) *c.List[card.Card] {
	list := c.NewList[card.Card](len(hand.Hand))
	for _, c := range hand.Hand {
		list.Append(decodeCard(c))
	}
	return list
}

// encodeCard converts a main.Card to a *Card.
func encodeCard(c card.Card) *Card {
	return &Card{
		Rank: encodeRank(c.Rank),
		Suit: encodeSuit(c.Suit),
	}
}

// decodeCard converts a *Card to a main.Card.
func decodeCard(c *Card) card.Card {
	return card.Card{
		Rank: decodeRank(c.Rank),
		Suit: decodeSuit(c.Suit),
	}
}

// encodeRank converts a main.Rank to a Rank.
func encodeRank(r card.Rank) Rank {
	return Rank(r)
}

// decodeRank converts a Rank to a main.Rank.
func decodeRank(r Rank) card.Rank {
	return card.Rank(r)
}

// encodeSuit converts a main.Suit to a Suit.
func encodeSuit(s card.Suit) Suit {
	switch s {
	case card.Spades:
		return Suit_SPADES
	case card.Clubs:
		return Suit_CLUBS
	case card.Diamonds:
		return Suit_DIAMONDS
	case card.Hearts:
		return Suit_HEARTS
	case card.NoSuit:
		return Suit_NO_SUIT
	default:
		panic(fmt.Sprintf("unknown suit %q", s))
	}
}

// decodeSuit converts a Suit to a main.Suit.
func decodeSuit(s Suit) card.Suit {
	switch s {
	case Suit_NO_SUIT:
		return card.NoSuit
	case Suit_SPADES:
		return card.Spades
	case Suit_CLUBS:
		return card.Clubs
	case Suit_DIAMONDS:
		return card.Diamonds
	case Suit_HEARTS:
		return card.Hearts
	default:
		panic(fmt.Sprintf("unknown suit %q", s))
	}
}

// encodeTrick converts a *c.List[main.PlayInfo] to a []*PlayInfo.
func encodeTrick(list *c.List[game.PlayInfo]) []*PlayInfo {
	out := make([]*PlayInfo, 0, list.Size())
	for _, pi := range *list {
		out = append(out, encodePlayInfo(pi))
	}
	return out
}

// decodeTrick converts a []*PlayInfo to a *c.List[main.PlayInfo].
func decodeTrick(trick []*PlayInfo) *c.List[game.PlayInfo] {
	list := c.NewList[game.PlayInfo](len(trick))
	for _, pi := range trick {
		list.Append(decodePlayInfo(pi))
	}
	return list
}

// encodePlayInfo converts a main.PlayInfo to a *PlayInfo.
func encodePlayInfo(pi game.PlayInfo) *PlayInfo {
	return &PlayInfo{
		Player: int32(pi.Player),
		Card:   encodeCard(pi.Card),
	}
}

// decodePlayInfo converts a *PlayInfo to a main.PlayInfo.
func decodePlayInfo(pi *PlayInfo) game.PlayInfo {
	return game.PlayInfo{
		Player: int(pi.Player),
		Card:   decodeCard(pi.Card),
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

// decodeValidPlays converts an []int32 to a *c.List[int].
func decodeValidPlays(vp []int32) *c.List[int] {
	list := c.NewList[int](len(vp))
	for _, n := range vp {
		list.Append(int(n))
	}
	return list
}
