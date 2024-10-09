package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/vishvananda/netlink"
)

func main() {
	handle, err := netlink.NewHandle()
	if err != nil {
		log.Fatalf("Error creating handle: %v", err)
	}
	defer handle.Close()

	// Получаем TCP сокеты
	sockets, err := handle.SocketDiagTCPInfo(netlink.FAMILY_V4)
	if err != nil {
		log.Fatalf("Error getting socket list: %v", err)
	}

	for _, socket := range sockets {
		localAddr := socket.InetDiagMsg.ID.Source
		destAddr := socket.InetDiagMsg.ID.Destination

		uid := socket.InetDiagMsg.UID
		userInfo, err := user.LookupId(fmt.Sprintf("%d", uid))
		if err != nil {
			log.Printf("Error looking up user: %v", err)
		}

		// Проверяем, является ли сокет активным (например, ESTABLISHED)
		if socket.InetDiagMsg.State == netlink.TCP_LISTEN {
			fmt.Printf("Local Address: %s:%d\n", localAddr.String(), socket.InetDiagMsg.ID.SourcePort)
			fmt.Printf("Destination: %s:%d\n", destAddr.String(), socket.InetDiagMsg.ID.DestinationPort)

			if userInfo != nil {
				fmt.Printf("User: %s\n", userInfo.Username)
			} else {
				fmt.Println("User: Unknown")
			}

			fmt.Println() // Добавляет пустую строку между сокетами для удобства
		}
	}
}
