package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	m "own-redis/models"
)

// KeyValueStore отвечает за хранение пар ключ-значение
type KeyValueStore struct {
	data map[string]string
	mu   sync.RWMutex
}

// Expire управляет сроками действия ключей
type Expire struct {
	expiry map[string]time.Time
	mu     sync.RWMutex
}

// StoreManager объединяет KeyValueStore и Expire
type StoreManager struct {
	kv     *KeyValueStore
	expire *Expire
}

// NewKeyValueStore создает экземпляр KeyValueStore
func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data: make(map[string]string),
	}
}

// NewExpire создает экземпляр Expire
func NewExpire() *Expire {
	return &Expire{
		expiry: make(map[string]time.Time),
	}
}

// NewStoreManager создает экземпляр StoreManager
func NewStoreManager() *StoreManager {
	return &StoreManager{
		kv:     NewKeyValueStore(),
		expire: NewExpire(),
	}
}

// Методы KeyValueStore

func (store *KeyValueStore) Set(key, value string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[key] = value
}

func (store *KeyValueStore) Get(key string) (string, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	value, exists := store.data[key]
	return value, exists
}

func (store *KeyValueStore) Delete(key string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.data, key)
}

// Методы Expire

func (expire *Expire) SetTTL(key string, ttl int64) {
	expire.mu.Lock()
	defer expire.mu.Unlock()
	if ttl > 0 {
		expire.expiry[key] = time.Now().Add(time.Duration(ttl) * time.Millisecond)
	} else {
		delete(expire.expiry, key)
	}
}

func (expire *Expire) IsExpired(key string) bool {
	expire.mu.RLock()
	defer expire.mu.RUnlock()
	if expiryTime, exists := expire.expiry[key]; exists {
		return time.Now().After(expiryTime)
	}
	return false
}

func (expire *Expire) Remove(key string) {
	expire.mu.Lock()
	defer expire.mu.Unlock()
	delete(expire.expiry, key)
}

// Методы StoreManager

func (sm *StoreManager) Set(key, value string, ttl int64) string {
	sm.kv.Set(key, value)
	sm.expire.SetTTL(key, ttl)
	return "OK"
}

func (sm *StoreManager) Get(key string) string {
	if sm.expire.IsExpired(key) {
		sm.kv.Delete(key)
		sm.expire.Remove(key)
		return "(nil)"
	}
	if value, exists := sm.kv.Get(key); exists {
		return value
	}
	return "(nil)"
}

// UDP-сервер

func handleConnection(sm *StoreManager, conn *net.UDPConn, addr *net.UDPAddr, message string) {
	parts := strings.Fields(strings.ToUpper(message))
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	switch command {
	case "PING":
		conn.WriteToUDP([]byte("PONG\n"), addr)
	case "SET":
		if len(parts) < 3 {
			conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'SET' command\n"), addr)
			return
		}
		key := parts[1]
		value := parts[2]
		ttl := int64(0)
		if len(parts) == 5 && strings.ToUpper(parts[3]) == "PX" {
			parsedTTL, err := strconv.ParseInt(parts[4], 10, 64)
			if err == nil {
				ttl = parsedTTL
			}
		}
		response := sm.Set(key, value, ttl)
		conn.WriteToUDP([]byte(response), addr)
	case "GET":
		if len(parts) != 2 {
			conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'GET' command\n"), addr)
			return
		}
		key := parts[1]
		response := sm.Get(key)
		conn.WriteToUDP([]byte(response), addr)
	default:
		conn.WriteToUDP([]byte("(error) ERR unknown command\n"), addr)
	}
}

func startServer(port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Server started on port %d\n", port)

	sm := NewStoreManager()

	buffer := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}
		go handleConnection(sm, conn, clientAddr, string(buffer[:n]))
	}
}

func main() {
	flag.Parse()
	if *m.Help {
		HelpShow()
		return
	}
	if *m.Port == 0 {
		fmt.Println("Port number can not be equal to 0 (port != 0)")
		return
	}
	startServer(*m.Port)
}
