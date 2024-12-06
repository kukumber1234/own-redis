package methods

import "sync"

type KeyValue struct {
	data map[string]string
	sync.RWMutex
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		data: make(map[string]string),
	}
}

func (store *KeyValue) Set(key, value string) {
	store.Lock()
	defer store.Unlock()
	store.data[key] = value
}

func (store *KeyValue) Get(key string) (string, bool) {
	store.RLock()
	defer store.RUnlock()
	value, exists := store.data[key]
	return value, exists
}

func (store *KeyValue) Delete(key string) {
	store.Lock()
	defer store.Unlock()
	delete(store.data, key)
}
