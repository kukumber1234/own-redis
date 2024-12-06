package main

import (
	"flag"
	"fmt"

	s "own-redis/internal/start"
	m "own-redis/models"
)

func main() {
	flag.Parse()

	if *m.Help {
		HelpShow()
		return
	}
	if *m.Port == 0 {
		fmt.Println("Port can not be equal to 0")
		return
	}

	s.StartServer(*m.Port)
}
