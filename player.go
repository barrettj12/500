package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"

	c "github.com/barrettj12/collections"
	"github.com/barrettj12/screen"
)

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
	NotifyBidWinner(player int, bid Bid)
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

// HumanPlayer is a player controlled by the user.
// It controls printing of the table state to the terminal.
type HumanPlayer struct {
	Hand  *c.List[Card]
	Table [4]Card
	valid *c.List[int]

	bid    Bid
	bidder int
}

func (p *HumanPlayer) NotifyHand(hand *c.List[Card]) {
	p.Hand = hand
	p.redrawBoard()
}

func (p *HumanPlayer) NotifyBid(player int, bid Bid) {
	fmt.Printf("%s bid %s\n", playerNames[player], bid)
}

func (p *HumanPlayer) NotifyBidWinner(player int, bid Bid) {
	fmt.Printf("%s won the bidding with %s\n", playerNames[player], bid)
}

func (p *HumanPlayer) NotifyPlay(player int, card Card) {
	p.Table[player] = card
	p.redrawBoard()
	fmt.Printf("%s played %s\n", playerNames[player], card)
}

func (p *HumanPlayer) NotifyTrickWinner(player int) {
	fmt.Printf("%s won the trick\n", playerNames[player])
	p.clearTable()
}

func (p *HumanPlayer) clearTable() {
	for i := 0; i < 4; i++ {
		p.Table[i] = Card{}
	}
}

func (p *HumanPlayer) NotifyHandResult(res HandResult) {
	fmt.Println(res.Info())
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

func (p *HumanPlayer) Drop3() *c.Set[int] {
	return prompt("Cards to dump [x,y,z]: ", func(s string) (*c.Set[int], error) {
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
}

func (p *HumanPlayer) Play(trick *c.List[Card], validPlays *c.List[int]) int {
	if validPlays.Size() == 1 {
		return E(p.valid.Get(0))
	}

	// Show valid cards
	p.valid = validPlays
	defer func() { p.valid = nil }()
	p.redrawBoard()

	return prompt("play card: ", func(s string) (int, error) {
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

// Prompt the user for input.
// A function can be provided to validate and transform the given input.
func prompt[T any](pr string, f func(string) (T, error)) T {
	s := bufio.NewScanner(os.Stdin)
	var res T

	for {
		fmt.Print(pr)
		s.Scan()
		if err := s.Err(); err != nil {
			panic(err)
		}

		input := s.Text()
		var err error
		res, err = f(input)
		if err == nil {
			break
		}

		// Invalid input
		fmt.Println(red(fmt.Sprintf("INVALID: %s", err)))
	}

	return res
}

func (p *HumanPlayer) redrawBoard() {
	screen.Clear()

	tmpl := E(template.New("test").Parse(`
Bid: {{.PrintBid}}

        {{index playerNames 2}}
        {{.FmtTable 2}}
  {{index playerNames 1}}         {{index playerNames 3}}
  {{.FmtTable 1}}         {{.FmtTable 3}}
        {{index playerNames 0}}
        {{.FmtTable 0}}

{{.PrintHand}}

`[1:]))

	E0(tmpl.Execute(screen.Writer(), p))
	screen.Update()
}

var playerNames = []string{"You", "Op1", "Pnr", "Op2"}

func (p *HumanPlayer) PrintBid() string {
	if p.bid == nil {
		return "—"
	}
	return fmt.Sprintf("%s by %s", p.bid, playerNames[p.bidder])
}

// Returns player's card suitable for printing.
// Always has 3 characters.
func FmtCard(card Card, grey bool) string {
	if (card == Card{}) {
		return "[_]"
	}

	var str string
	if grey {
		str = card.PrintGrey()
	} else {
		str = card.String()
	}

	if (card == JokerCard) || card.rank == 10 {
		return str
	}
	return str + " "
}

func (p *HumanPlayer) FmtTable(player int) string {
	card := p.Table[player]
	return FmtCard(card, false)
}

func (p *HumanPlayer) PrintHand() string {
	str := ""

	for i := 0; i < p.Hand.Size(); i++ {
		num := fmt.Sprintf("%-4d", i)
		if p.valid != nil && !p.valid.Contains(i) {
			num = grey(num)
		}
		str += num
	}
	str += "\n"
	for i, card := range *p.Hand {
		grey := p.valid != nil && !p.valid.Contains(i)
		c := FmtCard(card, grey)
		str += c + " "
	}

	return str
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
