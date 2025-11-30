package main

import (
	"log"

	"coin-collector/client"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := client.NewGame()

	if err := game.Init(); err != nil {
		log.Fatal("Failed to initialize client:", err)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Coin Collector - Client")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
