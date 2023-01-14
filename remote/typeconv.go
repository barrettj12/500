package remote

import (
	"fmt"

	"github.com/barrettj12/500/card"
	"github.com/barrettj12/500/game"

	c "github.com/barrettj12/collections"
)

// encodeHand converts a *c.List[main.Card] to a *Hand.
func encodeHand(list *c.List[card.Card]) *Hand {
	return &Hand{
		Hand: encodeList(list, encodeCard),
	}
}

// decodeHand converts a *Hand to a *c.List[main.Card].
func decodeHand(hand *Hand) *c.List[card.Card] {
	return decodeList(hand.Hand, decodeCard)
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
	return encodeList(list, encodePlayInfo)
}

// decodeTrick converts a []*PlayInfo to a *c.List[main.PlayInfo].
func decodeTrick(trick []*PlayInfo) *c.List[game.PlayInfo] {
	return decodeList(trick, decodePlayInfo)
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
	return encodeList(list, func(n int) int32 { return int32(n) })
}

// decodeValidPlays converts an []int32 to a *c.List[int].
func decodeValidPlays(vp []int32) *c.List[int] {
	return decodeList(vp, func(n int32) int { return int(n) })
}

// encodeList converts a *c.List[T] to a []U, using the specified conversion
// function f : T -> U.
func encodeList[T comparable, U any](list *c.List[T], f func(T) U) []U {
	out := make([]U, 0, list.Size())
	for _, t := range *list {
		out = append(out, f(t))
	}
	return out
}

// decodeList converts a []U to a *c.List[T], using the specified conversion
// function f : U -> T.
func decodeList[T comparable, U any](arr []U, f func(U) T) *c.List[T] {
	list := c.NewList[T](len(arr))
	for _, u := range arr {
		list.Append(f(u))
	}
	return list
}
