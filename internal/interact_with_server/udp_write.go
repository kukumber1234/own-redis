package write

import (
	"net"
	"strconv"
	"strings"

	me "own-redis/internal/methods"
	mo "own-redis/models"
)

func WriteToServer(sm *me.StoreManager, buf string, addr *net.UDPAddr, conn *net.UDPConn) {
	var setValue []string
	var value string
	var ttl int64
	parts := strings.Fields(buf)
	if len(parts) == 0 {
		return
	}
	command := strings.ToUpper(parts[0])

	switch command {
	case "PING":
		_, err := conn.WriteToUDP([]byte("PONG\n"), addr)
		if err != nil {
			mo.Logger.Println("Error sending PONG", err)
			return
		}
		mo.Logger.Println("Send 'PING' command")

	case "SET":
		if len(parts) < 3 {
			conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'SET' command\n"), addr)
			mo.Logger.Println("(error) ERR wrong number of arguments for 'SET' command")
			return
		}

		key := parts[1]
		switch len(parts) {
		case 3:
			ttl = 0
			value = parts[2]
		default:
			px := false
			if strings.ToUpper(parts[len(parts)-2]) == "PX" {
				parsedTTL, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
				if err == nil {
					ttl = parsedTTL
					px = true
				}
			}

			if px {
				parts = parts[:len(parts)-2]
			}

			for i := 2; i < len(parts); i++ {
				setValue = append(setValue, parts[i])
			}
			value = strings.Join(setValue, " ")
		}

		response := sm.Set(key, value, ttl)
		conn.WriteToUDP([]byte(response+"\n"), addr)
		mo.Logger.Println("Send 'SET' command")

	case "GET":
		if len(parts) != 2 {
			conn.WriteToUDP([]byte("(error) ERR wrong number of arguments for 'GET' command\n"), addr)
			mo.Logger.Println("(error) ERR wrong number of arguments for 'GET' command")
			return
		}

		key := parts[1]
		response := sm.Get(key)
		conn.WriteToUDP([]byte(response+"\n"), addr)
		mo.Logger.Println("Sent 'GET' command")

	default:
		_, err := conn.WriteToUDP([]byte("(error) Undefined command\n"), addr)
		if err != nil {
			return
		}
		mo.Logger.Println("(error) Undefined command")
	}
}
