package main

import "fmt"

func HelpShow() {
	fmt.Println(`$ ./own-redis --help
	Own Redis
	
	Usage:
	  own-redis [--port <N>]
	  own-redis --help
	
	Options:
	  --help       Show this screen.
	  --port N     Port number.`)
}
