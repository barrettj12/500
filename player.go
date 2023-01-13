package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

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
	NotifyPlayerNum(int)
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
	Play(trick *c.List[playInfo], validPlays *c.List[int]) int
	// JokerSuit asks for a suit for the Joker when it is led in no trumps
	// or misere.
	JokerSuit() Suit
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

// HumanPlayer implements Player.
var _ Player = &HumanPlayer{}

func (p *HumanPlayer) NotifyPlayerNum(int) {}

func (p *HumanPlayer) NotifyHand(hand *c.List[Card]) {
	p.Hand = hand
	p.redrawBoard()
}

func (p *HumanPlayer) NotifyBid(player int, bid Bid) {
	if (bid == Pass{}) {
		fmt.Printf("%s passed\n", p.PlayerName(player))
	} else {
		fmt.Printf("%s bid %s\n", p.PlayerName(player), bid)
	}
}

func (p *HumanPlayer) NotifyBidWinner(player int, bid Bid) {
	p.bid = bid
	fmt.Printf("%s won the bidding with %s\n", p.PlayerName(player), bid)
	pressToContinue()
}

func (p *HumanPlayer) NotifyPlay(player int, card Card) {
	p.Table[player] = card
	p.redrawBoard()
	// fmt.Printf("%s played %s\n", p.PlayerName(player), card)
}

func (p *HumanPlayer) NotifyTrickWinner(player int) {
	fmt.Printf("%s won the trick\n", p.PlayerName(player))
	pressToContinue()
	p.clearTable()
	p.redrawBoard()
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
	promptTricks := func() int {
		return prompt("Tricks [6-10]: ", func(s string) (int, error) {
			i, err := strconv.Atoi(s)
			if err != nil {
				return 0, err
			}
			if i < 6 || i > 10 {
				return 0, fmt.Errorf("invalid # of tricks %q", s)
			}
			return i, nil
		})
	}
	promptOpenMis := func() bool {
		return prompt("Open [o] or closed [c]? ", func(s string) (bool, error) {
			switch s {
			case "o":
				return true, nil
			case "c":
				return false, nil
			default:
				return false, fmt.Errorf(`expected "o" or "c", received %q`, s)
			}
		})
	}

	return prompt("Enter bid [s/c/d/h/n/m/p]: ", func(s string) (Bid, error) {
		switch s {
		case "s":
			return SuitBid{trumpSuit: Spades, tricks: promptTricks()}, nil
		case "c":
			return SuitBid{trumpSuit: Clubs, tricks: promptTricks()}, nil
		case "d":
			return SuitBid{trumpSuit: Diamonds, tricks: promptTricks()}, nil
		case "h":
			return SuitBid{trumpSuit: Hearts, tricks: promptTricks()}, nil
		case "n":
			return NoTrumpsBid{tricks: promptTricks()}, nil
		case "m":
			return MisereBid{open: promptOpenMis()}, nil
		case "p":
			return Pass{}, nil
		default:
			return nil, fmt.Errorf("unknown bid %q", s)
		}
	})
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

func (p *HumanPlayer) Play(trick *c.List[playInfo], validPlays *c.List[int]) int {
	time.Sleep(SLEEP)
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

func (p *HumanPlayer) JokerSuit() Suit {
	return prompt("Choose suit for Joker [s/c/d/h]: ", func(s string) (Suit, error) {
		switch s {
		case "s":
			return Spades, nil
		case "c":
			return Clubs, nil
		case "d":
			return Diamonds, nil
		case "h":
			return Hearts, nil
		default:
			return "", fmt.Errorf("unknown suit %q", s)
		}
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

func pressToContinue() {
	fmt.Println("[press enter to continue]")
	prompt("", func(s string) (int, error) { return 0, nil })
}

func (p *HumanPlayer) redrawBoard() {
	screen.Clear()

	tmpl := E(template.New("test").Parse(`
Bid: {{.PrintBid}}

        {{.PlayerName 2}}
        {{.FmtTable 2}}
  {{.PlayerName 1}}         {{.PlayerName 3}}
  {{.FmtTable 1}}         {{.FmtTable 3}}
        {{.PlayerName 0}}
        {{.FmtTable 0}}

{{.PrintHand}}

`[1:]))

	E0(tmpl.Execute(screen.Writer(), p))
	screen.Update()
}

func (p *HumanPlayer) PlayerName(player int) string {
	return []string{"You", "Op1", "Pnr", "Op2"}[player]
}

func (p *HumanPlayer) PrintBid() string {
	if p.bid == nil {
		return "—"
	}
	return fmt.Sprintf("%s by %s", p.bid, p.PlayerName(p.bidder))
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
	delay time.Duration
}

// Random implements Player.
var _ Player = &RandomPlayer{}

func (p *RandomPlayer) NotifyPlayerNum(int)                 {}
func (p *RandomPlayer) NotifyHand(*c.List[Card])            {}
func (p *RandomPlayer) NotifyBid(player int, bid Bid)       {}
func (p *RandomPlayer) NotifyBidWinner(player int, bid Bid) {}
func (p *RandomPlayer) NotifyPlay(player int, card Card)    {}
func (p *RandomPlayer) NotifyTrickWinner(player int)        {}
func (p *RandomPlayer) NotifyHandResult(res HandResult)     {}

func (p *RandomPlayer) Bid() Bid {
	time.Sleep(p.delay)
	// Random player doesn't bid
	return Pass{}
}

func (p *RandomPlayer) Drop3() *c.Set[int] {
	// Random player never wins bid, so we don't need to implement
	panic("RandomPlayer.Drop3 unimplemented")
}

func (p *RandomPlayer) Play(trick *c.List[playInfo], validPlays *c.List[int]) int {
	time.Sleep(p.delay)
	n := rand.Intn(validPlays.Size())
	return E(validPlays.Get(n))
}

func (p *RandomPlayer) JokerSuit() Suit {
	time.Sleep(p.delay)
	return []Suit{Spades, Clubs, Diamonds, Hearts}[rand.Intn(4)]
}
