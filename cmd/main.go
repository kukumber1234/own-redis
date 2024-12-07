package main

import (
	"flag"
	"log"

	l "own-redis/internal/logger"
	s "own-redis/internal/start_server"
	m "own-redis/models"
)

func main() {
	flag.Parse()

	if *m.Help {
		HelpShow()
		return
	}
	if err := l.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize the logger: %v", err)
	}
	if *m.Port == 0 {
		m.Logger.Println("Port can not be equal to 0")
		return
	}

	s.StartServer(*m.Port)
}
