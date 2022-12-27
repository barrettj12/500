package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"text/template"
	"time"

	c "github.com/barrettj12/collections"
	"github.com/kr/pretty"
)

const SLEEP = 500 * time.Millisecond

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	g := new500Game()
	g.redrawBoard()

	g.bid = getBid()
	fmt.Println("bid: ", g.bid)
	pressToContinue()

	// Kitty
	// yourHand := g.Players[0]
	// yourHand.Append(*g.kitty...)
	// sortHand(yourHand)

	// g.redrawBoard()
	// prompt("Cards to dump:", func(s string) bool {
	// 	// nums := strings.Split(s, ",")
	// 	return true
	// })

	for g.Players[0].Size() > 0 {
		g.clearTable()
		g.redrawBoard()
		time.Sleep(SLEEP)

		for i := 0; i < 4; i++ {
			nextPlayer := (i + g.leader) % 4
			hand := g.Players[nextPlayer]

			// Pick next card to play
			var j int
			if nextPlayer == 0 {
				// Your turn
				g.valid = g.validPlays(0)

				if g.valid.Size() == 1 {
					j = E(g.valid.Get(0))
				} else {
					// Show valid cards
					g.redrawBoard()

					in := prompt("play card: ",
						func(s string) bool {
							j, err := strconv.Atoi(s)
							return err == nil && g.valid.Contains(j)
						})
					j = E(strconv.Atoi(in))
				}

			} else {
				valid := g.validPlays(nextPlayer)
				r := rand.Intn(valid.Size())
				j = E(valid.Get(r))
			}

			card, err := hand.Remove(j)
			if err != nil {
				panic(err)
			}

			g.Table.Set(nextPlayer, card)
			g.valid = nil
			g.redrawBoard()
			time.Sleep(SLEEP)
		}

		// Determine winner
		k := g.bid.WhoWins(g.Table, g.leader)
		winner := (g.leader + k) % 4
		g.leader = winner
		if winner == 0 || winner == 2 {
			g.tricksWon++
		}

		fmt.Println("winner: ", g.PlayerNames[winner])
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
	Players [4]*c.List[Card]
	Table   *c.List[Card]

	PlayerNames map[int]string

	kitty     *c.List[Card]
	bid       Bid
	tricksWon int
	leader    int

	valid *c.List[int] // valid plays in current hand
}

func new500Game() gameState {
	deck := getDeck()
	deck.Shuffle()

	// Teams are (0, 2), (1, 3)
	players := [4]*c.List[Card]{
		E(deck.CopyPart(0, 10)),
		E(deck.CopyPart(10, 20)),
		E(deck.CopyPart(20, 30)),
		E(deck.CopyPart(30, 40)),
	}
	sortHand(players[0])

	kitty := E(deck.CopyPart(40, 43))

	return gameState{
		Players: players,
		Table:   c.AsList(make([]Card, 4)),
		PlayerNames: map[int]string{
			0: "You",
			1: "Op1",
			2: "Partner",
			3: "Op2",
		},
		kitty:     kitty,
		tricksWon: 0,
		leader:    0,
	}
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

func sortHand(hand *c.List[Card]) {
	hand.Sort(func(c, d Card) bool {
		if c.rank == Joker {
			return false
		}
		if d.rank == Joker {
			return true
		}

		suitOrder := map[Suit]int{Spades: 1, Clubs: 2, Diamonds: 3, Hearts: 4}
		suit1 := suitOrder[c.suit]
		suit2 := suitOrder[d.suit]
		if suit1 < suit2 {
			return true
		}
		if suit1 > suit2 {
			return false
		}

		// Same suit - aces high
		if c.rank == Ace {
			return false
		}
		if d.rank == Ace {
			return true
		}
		return c.rank < d.rank
	})
}

func (g *gameState) redrawBoard() {
	tmpl := E(template.New("test").Parse("\033[H\033[2J" + // clear screen
		`
      {{index .PlayerNames 2}}
        {{.FmtTable 2}}
  {{index .PlayerNames 1}}         {{index .PlayerNames 3}}
  {{.FmtTable 1}}         {{.FmtTable 3}}
        {{index .PlayerNames 0}}
        {{.FmtTable 0}}

{{.PrintHand}}

`[1:]))

	// Print to buffer first - less flickering?
	buf := &bytes.Buffer{}
	E0(tmpl.Execute(buf, g))
	io.Copy(os.Stdout, buf)

	// Write gamestate to file
	os.WriteFile(".gamestate.log", []byte(pretty.Sprint(g)), os.ModePerm)
}

// Returns player's card suitable for printing.
// Always has 3 characters.
func FmtCard(card Card) string {
	if (card == Card{}) {
		return "[_]"
	}
	if (card == JokerCard) || card.rank == 10 {
		return card.String()
	}
	return card.String() + " "
}

func (g *gameState) FmtTable(player int) string {
	card := E(g.Table.Get(player))
	return FmtCard(card)
}

func (g *gameState) PrintHand() string {
	// TODO: use g.valid, grey out invalid cards
	str := ""
	hand := g.Players[0]

	for i := 0; i < hand.Size(); i++ {
		num := fmt.Sprintf("%-4d", i)
		if g.valid != nil && !g.valid.Contains(i) {
			num = grey(num)
		}
		str += num
	}
	str += "\n"
	for i, card := range *hand {
		c := FmtCard(card)
		if g.valid != nil && !g.valid.Contains(i) {
			c = grey(c)
		}
		str += c + " "
	}

	return str
}

// Ask user for bid
func getBid() Bid {
	suitStr := prompt("Enter bid [s/c/d/h]: ", func(s string) bool {
		switch s {
		case "s", "c", "d", "h": //, "n", "m":
			return true
		}
		return false
	})
	// if suitStr == "m" {
	// 	// Misere
	// 	openStr := prompt("", func(s string) bool {
	// 		_, err := strconv.ParseBool(s)
	// 		return err == nil
	// 	})
	// 	bid = MisereBid{E(strconv.ParseBool(s))}
	// } else {
	// }
	tricksStr := prompt("Tricks [6-10]: ", func(s string) bool {
		i, err := strconv.Atoi(s)
		return err == nil && i >= 6 && i <= 10
	})
	return SuitBid{
		E(strconv.Atoi(tricksStr)),
		parseBid(suitStr),
	}
}

func prompt(pr string, validate func(string) bool) string {
	s := bufio.NewScanner(os.Stdin)
	var input string

	for {
		fmt.Print(pr)
		s.Scan()
		if err := s.Err(); err != nil {
			panic(err)
		}

		input = s.Text()
		if validate(input) {
			break
		}

		// Invalid input
		fmt.Println(red("INVALID"))
	}

	return input
}

func pressToContinue() {
	fmt.Println("[press enter to continue]")
	prompt("", func(s string) bool { return true })
}

func parseBid(s string) Suit {
	switch s {
	case "s":
		return Spades
	case "c":
		return Clubs
	case "d":
		return Diamonds
	case "h":
		return Hearts
	default:
		panic(fmt.Sprintf("unknown suit %q", s))
	}
}

func (g *gameState) clearTable() {
	for i := 0; i < 4; i++ {
		g.Table.Set(i, Card{})
	}
}

// Returns the indices of valid cards to play.
func (g *gameState) validPlays(player int) *c.List[int] {
	valids := c.NewList[int](g.Players[player].Size())
	hand := g.Players[player]

	for i, card := range *hand {
		if player == g.leader {
			// Can lead with any card
			valids.Append(i)
			continue
		}

		// We have to follow suit if we can
		leadCard := E(g.Table.Get(g.leader))
		leadSuit := g.bid.Suit(leadCard)
		if g.bid.Suit(card) == leadSuit {
			valids.Append(i)
			continue
		}

		// Check if we can't follow suit: then we can play anything
		numOfLeadSuit := hand.Count(func(_ int, c Card) bool { return g.bid.Suit(c) == leadSuit })

		if numOfLeadSuit == 0 {
			valids.Append(i)
			continue
		}
	}

	return valids
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
