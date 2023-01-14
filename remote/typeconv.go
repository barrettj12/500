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

// decodeHand converts a *Hand to a *c.List[main.Card].
func decodeHand(hand *Hand) *c.List[main.Card] {
	list := c.NewList[main.Card](len(hand.Hand))
	for _, c := range hand.Hand {
		list.Append(decodeCard(c))
	}
	return list
}

// encodeCard converts a main.Card to a *Card.
func encodeCard(c main.Card) *Card {
	return &Card{
		Rank: encodeRank(c.Rank),
		Suit: encodeSuit(c.Suit),
	}
}

// decodeCard converts a *Card to a main.Card.
func decodeCard(c *Card) main.Card {
	return main.Card{
		Rank: decodeRank(c.Rank),
		Suit: decodeSuit(c.Suit),
	}
}

// encodeRank converts a main.Rank to a Rank.
func encodeRank(r main.Rank) Rank {
	return Rank(r)
}

// decodeRank converts a Rank to a main.Rank.
func decodeRank(r Rank) main.Rank {
	return main.Rank(r)
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

// decodeSuit converts a Suit to a main.Suit.
func decodeSuit(s Suit) main.Suit {
	switch s {
	case Suit_NO_SUIT:
		return main.NoSuit
	case Suit_SPADES:
		return main.Spades
	case Suit_CLUBS:
		return main.Clubs
	case Suit_DIAMONDS:
		return main.Diamonds
	case Suit_HEARTS:
		return main.Hearts
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

// decodeTrick converts a []*PlayInfo to a *c.List[main.PlayInfo].
func decodeTrick(trick []*PlayInfo) *c.List[main.PlayInfo] {
	list := c.NewList[main.PlayInfo](len(trick))
	for _, pi := range trick {
		list.Append(decodePlayInfo(pi))
	}
	return list
}

// encodePlayInfo converts a main.PlayInfo to a *PlayInfo.
func encodePlayInfo(pi main.PlayInfo) *PlayInfo {
	return &PlayInfo{
		Player: int32(pi.Player),
		Card:   encodeCard(pi.Card),
	}
}

// decodePlayInfo converts a *PlayInfo to a main.PlayInfo.
func decodePlayInfo(pi *PlayInfo) main.PlayInfo {
	return main.PlayInfo{
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
