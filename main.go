package main

import (
	"fmt"
	"math/rand"
	"time"

	c "github.com/barrettj12/collections"
)

const SLEEP = 500 * time.Millisecond

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	g := new500Game()
	g.redrawBoard()

	// TODO: allow other players to bid
	g.bidder = 0
	g.bid = g.Players[0].Bid()
	for _, p := range g.Players {
		p.SetBid(g.bid)
	}
	fmt.Println("bid: ", g.bid)
	pressToContinue()

	// Kitty
	g.Players[g.bidder].AwardKitty(g.kitty)

	for trickNum := 0; trickNum < 10; trickNum++ {
		g.clearTable()
		g.redrawBoard()
		time.Sleep(SLEEP)

		trick := c.NewList[Card](4)

		for i := 0; i < 4; i++ {
			playerNum := (i + g.leader) % 4
			player := g.Players[playerNum]

			card := player.Play(trick)
			trick.Append(card)
			g.Table.Set(playerNum, card)

			g.redrawBoard()
			time.Sleep(SLEEP)
		}

		// Determine winner
		winner := g.whoWins()
		g.leader = winner
		if winner == 0 || winner == 2 {
			g.tricksWon++
		}

		fmt.Println("winner: ", g.Players[winner].Name())
		time.Sleep(SLEEP)
		pressToContinue()
	}

	fmt.Printf("won %d tricks\n", g.tricksWon)
	if g.bid.Won(g.tricksWon) {
		fmt.Println("YOU WON!!!")
	} else {
		fmt.Println("You lost :(")
	}
}

// gameState represents the current state of a 500 game
type gameState struct {
	Players [4]Player
	Table   *c.List[Card]

	kitty     *c.List[Card]
	bid       Bid
	bidder    int
	tricksWon int
	leader    int
}

func new500Game() *gameState {
	deck := getDeck()
	deck.Shuffle()

	// Teams are (0, 2), (1, 3)
	players := [4]Player{
		nil,
		NewRandomPlayer("Op1", E(deck.CopyPart(10, 20))),
		NewRandomPlayer("Partner", E(deck.CopyPart(20, 30))),
		NewRandomPlayer("Op2", E(deck.CopyPart(30, 40))),
	}
	// sortHand(players[0])

	kitty := E(deck.CopyPart(40, 43))

	g := &gameState{
		Players:   players,
		Table:     c.AsList(make([]Card, 4)),
		kitty:     kitty,
		tricksWon: 0,
		leader:    0,
	}
	g.Players[0] = NewHumanPlayer("You", E(deck.CopyPart(0, 10)), g)
	g.redrawBoard()
	return g
}

// Returns the 500 deck
func getDeck() *c.List[Card] {
	return c.AsList([]Card{
		{4, Diamonds}, {4, Hearts},
		{5, Spades}, {5, Clubs}, {5, Diamonds}, {5, Hearts},
		{6, Spades}, {6, Clubs}, {6, Diamonds}, {6, Hearts},
		{7, Spades}, {7, Clubs}, {7, Diamonds}, {7, Hearts},
		{8, Spades}, {8, Clubs}, {8, Diamonds}, {8, Hearts},
		{9, Spades}, {9, Clubs}, {9, Diamonds}, {9, Hearts},
		{10, Spades}, {10, Clubs}, {10, Diamonds}, {10, Hearts},
		{Jack, Spades}, {Jack, Clubs}, {Jack, Diamonds}, {Jack, Hearts},
		{Queen, Spades}, {Queen, Clubs}, {Queen, Diamonds}, {Queen, Hearts},
		{King, Spades}, {King, Clubs}, {King, Diamonds}, {King, Hearts},
		{Ace, Spades}, {Ace, Clubs}, {Ace, Diamonds}, {Ace, Hearts},
		JokerCard,
	})
}

func pressToContinue() {
	fmt.Println("[press enter to continue]")
	prompt("", func(s string) (int, error) { return 0, nil })
}

// Utility functions

// E evaluates the given function and panics if it returns a non-nil error.
func E0(err error) {
	if err != nil {
		panic(err)
	}
}

func E[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
