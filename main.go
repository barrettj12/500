package main

import (
	"math/rand"
	"time"

	"github.com/barrettj12/500/controller"
	"github.com/barrettj12/500/player"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	ct := controller.Controller{
		Players: [4]player.Player{
			&player.HumanPlayer{},
			&player.RandomPlayer{Delay: player.SLEEP},
			&player.RandomPlayer{Delay: player.SLEEP},
			&player.RandomPlayer{Delay: player.SLEEP},
		},
	}
	ct.Play()
}
