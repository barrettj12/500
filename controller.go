package main

import (
	c "github.com/barrettj12/collections"
)

// Controller handles the gameplay of a 500 game.
// It keeps track of the game state and the hands, transmits events to players,
// and contacts players to make plays, checking these plays are valid.
type Controller struct {
	players [4]Player

	hands      [4]*c.List[Card]
	kitty      *c.List[Card]
	bid        Bid
	contractor int
}

// Play plays this game of 500.
func (ct *Controller) Play() {
	// Shuffle cards
	deck := getDeck()
	deck.Shuffle()

	// Deal cards and notify each player of their hand
	for i := 0; i < 4; i++ {
		ct.hands[i] = E(deck.CopyPart(i*10, i*10+10))
		ct.players[i].NotifyHand(ct.hands[i])
	}
	ct.kitty = E(deck.CopyPart(40, 43))

	// Bidding
	hasPassed := c.AsList([]bool{false, false, false, false})
	winningBidder := -1
	var winningBid Bid
	// Player 0 starts the bidding
	bidder := 0

	for {
		if E(hasPassed.Get(bidder)) {
			bidder = (bidder + 1) % 4
			continue
		}
		if hasPassed.Count(func(i int, b bool) bool { return b == true }) == 3 {
			// All other players have passed - bid is won
			break
		}

		newBid := retryTillValid(func() (Bid, bool) {
			bid := ct.players[bidder].Bid()
			if (bid == Pass{}) {
				hasPassed.Set(bidder, true)
				return bid, true
			}
			if bid.Value() > winningBid.Value() {
				winningBid = bid
				winningBidder = bidder
				return bid, true
			}
			// Otherwise, bid not valid, so we will ask again
			return nil, false
		})

		// Notify other players of bid
		for i := 0; i < 4; i++ {
			ct.players[i].NotifyBid(bidder, newBid)
		}

		// Bidding passes to next player
		bidder = (bidder + 1) % 4
	}

	if winningBidder == -1 {
		// Everyone passed - re-deal
		return
	}
	ct.bid = winningBid
	ct.contractor = winningBidder

	// Kitty
	ct.hands[ct.contractor].Append(*ct.kitty...)
	ct.players[ct.contractor].NotifyHand(ct.hands[ct.contractor])

	// Ask contractor to drop 3 cards from hand
	toDrop := retryTillValid(func() (*c.Set[int], bool) {
		toDrop := ct.players[ct.contractor].Drop3()
		if toDrop.Size() != 3 {
			return nil, false
		}
		for n := range *toDrop {
			if n < 0 || n > 12 {
				return nil, false
			}
		}
		return toDrop, true
	})
	ct.hands[ct.contractor] = ct.hands[ct.contractor].Filter(func(i int, _ Card) bool { return !toDrop.Contains(i) })

	// - ask player for play
	// - check play is valid
	// - notify player if play is valid or not
	// - notify all players of play
	// - ...
	// - notify players of trick winner
}

// retryTillValid repeatedly calls the given function until it returns a true
// response, then returns the function's other output.
func retryTillValid[T any](f func() (T, bool)) T {
	for {
		t, valid := f()
		if valid {
			return t
		}
	}
}

// Player represents a 500 player. The controller interacts with the player
// in two distinct ways:
//   - Events: the controller informs a player of something that has happened
//     (e.g. another player making a bid or play).
//   - Requests: the controller asks for input from a player (e.g. what card
//     they would like to play).
type Player interface {
	// Events
	NotifyHand(*c.List[Card])
	NotifyBid(player int, bid Bid)

	// Requests
	Bid() Bid
	Drop3() *c.Set[int]
}
