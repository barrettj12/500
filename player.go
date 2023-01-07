package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	c "github.com/barrettj12/collections"
)

// Player is an interface representing a 500 player (human or computer).
// Each Player is responsible for:
// - keeping track of its own hand
// - ensuring it has the correct number of cards in hand at each time
// - ensuring all its plays are valid
type Player interface {
	Name() string
	Bid() Bid
	SetBid(Bid)
	AwardKitty(*c.List[Card])
	Play(trick *c.List[Card]) Card
}

type HumanPlayer struct {
	name string
	hand *c.List[Card]
	bid  Bid

	gs    *gameState
	valid *c.List[int]
}

func NewHumanPlayer(name string, hand *c.List[Card], gs *gameState) *HumanPlayer {
	p := &HumanPlayer{
		name: name,
		hand: hand,
		gs:   gs,
	}
	NoTrumpsBid{}.SortHand(hand)
	return p
}

func (p *HumanPlayer) Name() string {
	return p.name
}

func (p *HumanPlayer) Bid() Bid {
	suit := prompt("Enter bid [s/c/d/h]: ", func(s string) (Suit, error) {
		switch s {
		case "s":
			return Spades, nil
		case "c":
			return Clubs, nil
		case "d":
			return Diamonds, nil
		case "h":
			return Hearts, nil
		// case "n":
		// case "m":
		default:
			return "", fmt.Errorf("unknown bid %q", s)
		}
	})

	tricks := prompt("Tricks [6-10]: ", func(s string) (int, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if i < 6 || i > 10 {
			return 0, fmt.Errorf("invalid # of tricks %q", s)
		}
		return i, nil
	})

	return SuitBid{
		tricks,
		suit,
	}
}

func (p *HumanPlayer) SetBid(bid Bid) {
	p.bid = bid
}

func (p *HumanPlayer) AwardKitty(kitty *c.List[Card]) {
	p.hand.Append(*kitty...)
	p.bid.SortHand(p.hand)

	p.redrawBoard()
	toDump := prompt("Cards to dump [x,y,z]: ", func(s string) (*c.Set[int], error) {
		nums := strings.Split(s, ",")
		if len(nums) != 3 {
			return nil, fmt.Errorf("expected 3 nums, received %d", len(nums))
		}

		ints := c.NewSet[int](3)
		for _, str := range nums {
			n, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			if n < 0 || n > 12 {
				return nil, fmt.Errorf("%d is out of range", n)
			}
			ints.Add(n)
		}

		// Ensure numbers are unique
		if ints.Size() != 3 {
			return nil, fmt.Errorf("repeated numbers in %s", s)
		}

		return ints, nil
	})
	p.hand = p.hand.Filter(func(i int, _ Card) bool { return !toDump.Contains(i) })
}

func (p *HumanPlayer) Play(trick *c.List[Card]) Card {
	p.valid = p.bid.ValidPlays(trick, p.hand)

	var j int
	if p.valid.Size() == 1 {
		j = E(p.valid.Get(0))
	} else {
		// Show valid cards
		p.redrawBoard()

		j = prompt("play card: ", func(s string) (int, error) {
			j, err := strconv.Atoi(s)
			if err != nil {
				return 0, err
			}
			if !p.valid.Contains(j) {
				return 0, fmt.Errorf("invalid play")
			}
			return j, nil
		})
	}

	card := E(p.hand.Remove(j))
	p.valid = nil
	return card
}

func (p *HumanPlayer) redrawBoard() {
	// TODO: move drawing logic into HumanPlayer
	p.gs.redrawBoard()
}

// Plays a random (valid) card each round.
type RandomPlayer struct {
	name string
	hand *c.List[Card]
	bid  Bid
}

func NewRandomPlayer(name string, hand *c.List[Card]) *RandomPlayer {
	return &RandomPlayer{
		name: name,
		hand: hand,
	}
}

func (p *RandomPlayer) Name() string {
	return p.name
}

func (p *RandomPlayer) Bid() Bid {
	return nil
}

func (p *RandomPlayer) SetBid(bid Bid) {
	p.bid = bid
}

func (p *RandomPlayer) AwardKitty(kitty *c.List[Card]) {
	// Ignore kitty
}

func (p *RandomPlayer) Play(trick *c.List[Card]) Card {
	valid := p.bid.ValidPlays(trick, p.hand)
	r := rand.Intn(valid.Size())
	j := E(valid.Get(r))

	card := E(p.hand.Remove(j))
	return card
}
