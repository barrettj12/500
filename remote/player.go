package remote

import (
	context "context"

	main "github.com/barrettj12/500"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

// RemoteController is the gRPC server side, which is run on the remote
// player's machine.
type RemoteController struct {
	UnimplementedPlayerServer

	player main.Player
}

var _ PlayerServer = &RemoteController{}

func (c *RemoteController) NotifyPlayerNum(_ context.Context, n *wrapperspb.Int32Value) (*emptypb.Empty, error) {
	c.player.NotifyPlayerNum(int(n.Value))
	return nil, nil
}

func (c *RemoteController) NotifyHand(_ context.Context, h *Hand) (*emptypb.Empty, error) {
	c.player.NotifyHand(decodeHand(h))
	return nil, nil
}

func (c *RemoteController) NotifyPlay(_ context.Context, pi *PlayInfo) (*emptypb.Empty, error) {
	c.player.NotifyPlay(
		int(pi.Player),
		decodeCard(pi.Card),
	)
	return nil, nil
}

func (c *RemoteController) NotifyTrickWinner(_ context.Context, winner *wrapperspb.Int32Value) (*emptypb.Empty, error) {
	c.player.NotifyTrickWinner(int(winner.Value))
	return nil, nil
}

func (c *RemoteController) Play(_ context.Context, req *PlayRequest) (*wrapperspb.Int32Value, error) {
	n := c.player.Play(
		decodeTrick(req.Trick),
		decodeValidPlays(req.ValidPlays),
	)
	return &wrapperspb.Int32Value{Value: int32(n)}, nil
}
