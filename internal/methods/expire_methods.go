package methods

import (
	"sync"
	"time"
)

type Expire struct {
	expire map[string]time.Time
	sync.RWMutex
}

func NewExpire() *Expire {
	return &Expire{
		expire: make(map[string]time.Time),
	}
}

func (ex *Expire) SetTTL(key string, ttl int64) {
	ex.Lock()
	defer ex.Unlock()
	if ttl > 0 {
		ex.expire[key] = time.Now().Add(time.Duration(ttl) * time.Millisecond)
	} else {
		delete(ex.expire, key)
	}
}

func (ex *Expire) IsExpired(key string) bool {
	ex.RLock()
	defer ex.RUnlock()
	if expiryTime, exists := ex.expire[key]; exists {
		return time.Now().After(expiryTime)
	}
	return false
}

func (ex *Expire) Remove(key string) {
	ex.Lock()
	defer ex.Unlock()
	delete(ex.expire, key)
}
