package main

import (
	"fmt"

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

	leader    int
	tricks    [10]*c.List[Card]
	tricksWon [4]int
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

		numPasses := hasPassed.Count(func(i int, b bool) bool { return b == true })
		if numPasses == 4 {
			// All players passed - re-deal
			for i := 0; i < 4; i++ {
				ct.players[i].NotifyHandResult(redeal{})
			}
			return
		}
		if winningBid != nil && numPasses == 3 {
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

	// Play game
	ct.leader = ct.contractor
	for trickNum := 0; trickNum < 10; trickNum++ {
		ct.tricks[trickNum] = c.NewList[Card](4)

		for i := 0; i < 4; i++ {
			playerNum := (i + ct.leader) % 4

			validPlays := ct.bid.ValidPlays(ct.tricks[trickNum], ct.hands[playerNum])
			cardNum := retryTillValid(func() (int, bool) {
				cardNum := ct.players[playerNum].Play(
					ct.tricks[trickNum], // trick so far
					validPlays,
				)
				return cardNum, validPlays.Contains(cardNum)
			})

			card := E(ct.hands[playerNum].Remove(cardNum))
			ct.tricks[trickNum].Append(card)

			// Notify players of played card
			for i := 0; i < 4; i++ {
				ct.players[i].NotifyPlay(playerNum, card)
			}
			ct.players[playerNum].NotifyHand(ct.hands[playerNum])
		}

		// Determine winner
		winner := ct.whoWins(trickNum)
		ct.tricksWon[winner]++
		ct.leader = winner

		for i := 0; i < 4; i++ {
			ct.players[i].NotifyTrickWinner(winner)
		}
	}

	// Determine hand result
	teamTricks := ct.tricksWon[ct.contractor] + ct.tricksWon[(ct.contractor+2)%4]
	var res HandResult
	if ct.bid.Won(teamTricks) {
		res = bidWon{ct.bid, teamTricks}
	} else {
		res = bidWon{ct.bid, teamTricks}
	}

	for i := 0; i < 4; i++ {
		ct.players[i].NotifyHandResult(res)
	}
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

func (ct *Controller) whoWins(trickNum int) int {
	leadCard := E(ct.tricks[trickNum].Get(0))
	order := ct.bid.CardOrder(leadCard)
	winner := ct.leader // whoever lead wins by default

	for _, card := range *order {
		i, err := ct.tricks[trickNum].Find(card)
		if err != nil {
			// not found
			continue
		}

		winner = i
		break
	}

	return winner
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
	NotifyPlay(player int, card Card)
	NotifyTrickWinner(player int)
	NotifyHandResult(res HandResult)

	// Requests
	Bid() Bid
	Drop3() *c.Set[int]
	// Play asks the player to play a card on the given trick.
	// The returned response must be an element of validPlays.
	Play(trick *c.List[Card], validPlays *c.List[int]) int
}

// HandResult represents the outcome of a hand.
type HandResult interface {
	Info() string
}

// redeal is a HandResult representing all players passing during bidding.
type redeal struct{}

func (r redeal) Info() string {
	return "Re-deal due to all players passing"
}

// bidWon says that the contractor won their bid.
type bidWon struct {
	bid    Bid
	tricks int
}

func (r bidWon) Info() string {
	return fmt.Sprintf("Contractors won their bid of %s with %d tricks",
		r.bid, r.tricks)
}

// bidLost says that the contractors lost their bid.
type bidLost struct {
	bid    Bid
	tricks int
}

func (r bidLost) Info() string {
	return fmt.Sprintf("Contractors lost their bid of %s with %d tricks",
		r.bid, r.tricks)
}
