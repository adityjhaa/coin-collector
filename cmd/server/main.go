package main

import (
	"log"

	"coin-collector/server"
)

func main() {
	log.Println("=== Coin Collector Server Starting ===")

	s := server.NewServer()
	if err := s.StartNetwork(); err != nil {
		log.Fatal("network start:", err)
	}

	go s.StartGameLoop()

	log.Println("server running; waiting for clients...")
	select {}
}
