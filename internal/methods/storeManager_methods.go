package methods

type StoreManager struct {
	keyValue *KeyValue
	expire   *Expire
}

func NewStoreManager() *StoreManager {
	return &StoreManager{
		keyValue: NewKeyValue(),
		expire:   NewExpire(),
	}
}

func (sm StoreManager) Set(key, value string, ttl int64) string {
	sm.keyValue.Set(key, value)
	sm.expire.SetTTL(key, ttl)
	return "OK"
}

func (sm StoreManager) Get(key string) string {
	if sm.expire.IsExpired(key) {
		sm.keyValue.Delete(key)
		sm.expire.Remove(key)
		return "(nil)"
	}
	if value, exists := sm.keyValue.Get(key); exists {
		return value
	}
	return "(nil)"
}
