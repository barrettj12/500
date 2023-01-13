package remote

import (
	context "context"

	main "github.com/barrettj12/500"
	"google.golang.org/protobuf/types/known/wrapperspb"

	c "github.com/barrettj12/collections"
)

// RemotePlayer communicates with a remote player using a gRPC client.
type RemotePlayer struct {
	client PlayerClient
}

// RemotePlayer implements Player.
var _ main.Player = &RemotePlayer{}

func (p *RemotePlayer) NotifyPlayerNum(playerNum int) {
	_, err := p.client.NotifyPlayerNum(
		context.Background(),
		&wrapperspb.Int32Value{Value: int32(playerNum)},
	)
	panicIfNotNil(err)
}

func (p *RemotePlayer) NotifyHand(hand *c.List[main.Card]) {
	_, err := p.client.NotifyHand(
		context.Background(),
		encodeHand(hand),
	)
	panicIfNotNil(err)
}

func (p *RemotePlayer) NotifyBid(player int, bid main.Bid)       {}
func (p *RemotePlayer) NotifyBidWinner(player int, bid main.Bid) {}

func (p *RemotePlayer) NotifyPlay(player int, card main.Card) {
	_, err := p.client.NotifyPlay(
		context.Background(),
		&PlayInfo{
			Player: int32(player),
			Card:   encodeCard(card),
		},
	)
	panicIfNotNil(err)
}

func (p *RemotePlayer) NotifyTrickWinner(player int) {
	_, err := p.client.NotifyTrickWinner(
		context.Background(),
		&wrapperspb.Int32Value{Value: int32(player)},
	)
	panicIfNotNil(err)
}

func (p *RemotePlayer) NotifyHandResult(res main.HandResult) {}

func (p *RemotePlayer) Bid() main.Bid {
	// TODO: implement properly
	return main.Pass{}
}

func (p *RemotePlayer) Drop3() *c.Set[int] {
	// Should not be called yet, as we don't have RemotePlayer bidding.
	// TODO: implement properly
	panic("RemotePlayer.Drop3 not implemented")
}

func (p *RemotePlayer) Play(trick *c.List[main.PlayInfo], validPlays *c.List[int]) int {
	resp, err := p.client.Play(
		context.Background(),
		&PlayRequest{
			Trick:      encodeTrick(trick),
			ValidPlays: encodeValidPlays(validPlays),
		},
	)
	panicIfNotNil(err)
	return int(resp.Value)
}

func (p *RemotePlayer) JokerSuit() main.Suit {
	// TODO: implement properly
	panic("RemotePlayer.Drop3 not implemented")
}

func panicIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
