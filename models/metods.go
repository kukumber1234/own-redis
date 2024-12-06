package models

import (
	"flag"
	"sync"
	"time"
)

var (
	Port = flag.Int("port", 8080, "give port number")
	Help = flag.Bool("help", false, "show information about this code")
)

type KeyValue struct {
	data map[string]string
	sync.RWMutex
}

type Expire struct {
	expire map[string]time.Time
	sync.RWMutex
}

type StoreManager struct {
	keyValue *KeyValue
	expire   *Expire
}
