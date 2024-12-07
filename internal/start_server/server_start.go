package start

import (
	"fmt"
	"net"

	w "own-redis/internal/interact_with_server"
	me "own-redis/internal/methods"
	mo "own-redis/models"
)

func StartServer(port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		mo.Logger.Println("Error creating address", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Server started on port %d\n", *mo.Port)
	mo.Logger.Printf("Server started on port %d\n", *mo.Port)
	sm := me.NewStoreManager()

	buf := make([]byte, 1024)
	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			mo.Logger.Println("Error read from udp", err)
			return
		}
		go w.WriteToServer(sm, string(buf[:rlen]), remote, conn)
	}
}
