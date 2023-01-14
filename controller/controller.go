package controller

import (
	"os"
	"time"

	"github.com/barrettj12/500/card"
	"github.com/barrettj12/500/game"
	"github.com/barrettj12/500/player"
	"github.com/barrettj12/500/util"

	c "github.com/barrettj12/collections"
	"github.com/kr/pretty"
)

// Controller handles the gameplay of a 500 game.
// It keeps track of the game state and the hands, transmits events to players,
// and contacts players to make plays, checking these plays are valid.
type Controller struct {
	Players [4]player.Player

	hands      [4]*c.List[card.Card]
	kitty      *c.List[card.Card]
	bidHistory []bidInfo
	bid        game.Bid
	contractor int

	leader       int
	trickHistory [10]trickInfo
}

// Play plays this game of 500.
func (ct *Controller) Play() {
	for i := 0; i < 4; i++ {
		ct.Players[i].NotifyPlayerNum(i)
	}

	// Shuffle cards
	deck := game.GetDeck()
	deck.Shuffle()

	// Deal cards and notify each player of their hand
	for i := 0; i < 4; i++ {
		ct.hands[i] = util.E(deck.CopyPart(i*10, i*10+10))
		game.NoTrumpsBid{}.SortHand(ct.hands[i])
		ct.Players[i].NotifyHand(ct.hands[i])
	}
	ct.kitty = util.E(deck.CopyPart(40, 43))

	// Bidding
	hasPassed := c.AsList([]bool{false, false, false, false})
	winningBidder := -1
	var winningBid game.Bid
	// Player 0 starts the bidding
	bidder := 0

	for {
		ct.writeGamestate()

		numPasses := hasPassed.Count(func(i int, b bool) bool { return b })
		if numPasses == 4 {
			// All players passed - re-deal
			for i := 0; i < 4; i++ {
				ct.Players[i].NotifyHandResult(game.Redeal{})
			}
			return
		}
		if winningBid != nil && numPasses == 3 {
			// All other players have passed - bid is won
			break
		}

		if util.E(hasPassed.Get(bidder)) {
			bidder = (bidder + 1) % 4
			continue
		}

		newBid := retryTillValid(func() (game.Bid, bool) {
			b := ct.Players[bidder].Bid()
			if (b == game.Pass{}) {
				hasPassed.Set(bidder, true)
				return b, true
			}
			if winningBid == nil || b.Value() > winningBid.Value() {
				winningBid = b
				winningBidder = bidder
				return b, true
			}
			// Otherwise, bid not valid, so we will ask again
			return nil, false
		})
		ct.bidHistory = append(ct.bidHistory, bidInfo{bidder, newBid})

		// Notify other players of bid
		for i := 0; i < 4; i++ {
			ct.Players[i].NotifyBid(bidder, newBid)
		}

		// Bidding passes to next player
		bidder = (bidder + 1) % 4
	}

	ct.bid = winningBid
	ct.contractor = winningBidder
	for i := 0; i < 4; i++ {
		ct.Players[i].NotifyBidWinner(ct.contractor, ct.bid)
		// Sort hand according to bid
		ct.bid.SortHand(ct.hands[i])
		ct.Players[i].NotifyHand(ct.hands[i])
	}

	// Kitty
	ct.hands[ct.contractor].Append(*ct.kitty...)
	ct.bid.SortHand(ct.hands[ct.contractor])
	ct.Players[ct.contractor].NotifyHand(ct.hands[ct.contractor])
	ct.writeGamestate()

	// Ask contractor to drop 3 cards from hand
	toDrop := retryTillValid(func() (*c.Set[int], bool) {
		toDrop := ct.Players[ct.contractor].Drop3()
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
	ct.hands[ct.contractor] = ct.hands[ct.contractor].Filter(func(i int, _ card.Card) bool { return !toDrop.Contains(i) })
	ct.Players[ct.contractor].NotifyHand(ct.hands[ct.contractor])
	ct.writeGamestate()

	// Play game
	ct.leader = ct.contractor
	for trickNum := 0; trickNum < 10; trickNum++ {
		ct.trickHistory[trickNum] = newTrickInfo(ct.leader)

		for i := 0; i < 4; i++ {
			playerNum := (i + ct.leader) % 4
			if _, ok := ct.bid.(game.MisereBid); ok {
				// Contractor's partner doesn't play in Misere
				if playerNum == (ct.contractor+2)%4 {
					continue
				}
			}

			validPlays := ct.bid.ValidPlays(ct.trickHistory[trickNum].plays, ct.hands[playerNum])
			var cardNum int
			if validPlays.Size() == 1 {
				time.Sleep(player.SLEEP)
				cardNum = util.E(validPlays.Get(0))
			} else {
				cardNum = retryTillValid(func() (int, bool) {
					cardNum := ct.Players[playerNum].Play(
						ct.trickHistory[trickNum].plays, // trick so far
						validPlays,
					)
					return cardNum, validPlays.Contains(cardNum)
				})
			}

			cd := util.E(ct.hands[playerNum].Remove(cardNum))
			ct.trickHistory[trickNum].AddPlay(playerNum, cd)

			// Handle Joker lead in no trumps / misere
			// TODO: doesn't seem to be working
			if cd == card.JokerCard && playerNum == ct.leader {
				switch b := ct.bid.(type) {
				case *game.NoTrumpsBid:
					jokerSuit := ct.Players[playerNum].JokerSuit()
					b.JokerSuit = jokerSuit
				case *game.MisereBid:
					jokerSuit := ct.Players[playerNum].JokerSuit()
					b.NoTrumpsBid.JokerSuit = jokerSuit
				}
			}

			// Notify players of played card
			for i := 0; i < 4; i++ {
				ct.Players[i].NotifyPlay(playerNum, cd)
			}
			ct.Players[playerNum].NotifyHand(ct.hands[playerNum])
			ct.writeGamestate()
		}

		// Determine winner
		winner := ct.trickHistory[trickNum].Winner(ct.bid)
		ct.leader = winner

		for i := 0; i < 4; i++ {
			ct.Players[i].NotifyTrickWinner(winner)
		}
		ct.writeGamestate()
	}

	// Determine hand result
	teamTricks := 0
	for _, t := range ct.trickHistory {
		if t.winner == ct.contractor || t.winner == (ct.contractor+2)%4 {
			teamTricks++
		}
	}

	var res game.HandResult
	if ct.bid.Won(teamTricks) {
		res = game.BidWon{ct.bid, teamTricks}
	} else {
		res = game.BidLost{ct.bid, teamTricks}
	}

	for i := 0; i < 4; i++ {
		ct.Players[i].NotifyHandResult(res)
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

func (ct *Controller) writeGamestate() {
	os.WriteFile(".gamestate.log", []byte(pretty.Sprint(ct)), os.ModePerm)
}

// bidInfo holds information about a bid.
type bidInfo struct {
	player int
	bid    game.Bid
}

// trickInfo holds information about a trick.
type trickInfo struct {
	leader int
	plays  *c.List[game.PlayInfo]
	winner int
}

func newTrickInfo(leader int) trickInfo {
	return trickInfo{
		leader: leader,
		plays:  c.NewList[game.PlayInfo](4),
	}
}

func (t *trickInfo) AddPlay(player int, card card.Card) {
	t.plays.Append(game.PlayInfo{Player: player, Card: card})
}

func (t *trickInfo) Winner(bid game.Bid) int {
	leadCard := util.E(t.plays.Get(0)).Card
	order := bid.CardOrder(leadCard)
	var winner int

	for _, card := range *order {
		p, err := t.plays.Filter(
			func(i int, p game.PlayInfo) bool { return p.Card == card }).Get(0)
		if err != nil {
			// not found
			continue
		}

		winner = p.Player
		break
	}

	t.winner = winner
	return t.winner
}
