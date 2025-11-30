package main

import (
	"coin-collector/server"
	"log"
)

func main() {
	log.Println("=== Coin Collector Server Starting ===")

	s := server.NewServer()

	if err := s.StartNetwork(); err != nil {
		log.Fatal("Failed to start network:", err)
	}

	go s.StartGameLoop()

	log.Println("Server initialized. Awaiting clients...")

	select {}
}
