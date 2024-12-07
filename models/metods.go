package models

import (
	"flag"
	"log"
)

var (
	Port = flag.Int("port", 8080, "give port number")
	Help = flag.Bool("help", false, "show information about this code")
)

var Logger *log.Logger
