package main

import (
	"fmt"
	"log"
	"os/user"
	"strconv"

	"github.com/vishvananda/netlink"
)

func main() {
	handle, err := netlink.NewHandle()
	if err != nil {
		log.Fatalf("Error creating handle: %v", err)
	}
	defer handle.Close()

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

		inode := strconv.Itoa(int(socket.InetDiagMsg.INode))
		UserID := socket.InetDiagMsg.UID

		fmt.Printf("Local Address: %s:%d\n", localAddr.String(), socket.InetDiagMsg.ID.SourcePort)
		fmt.Printf("Destination: %s:%d\n", destAddr.String(), socket.InetDiagMsg.ID.DestinationPort)
		fmt.Printf("INode: %s\n", inode)
		fmt.Println("State:", socket.InetDiagMsg.State)
		// Вывод размера очереди получения
		fmt.Printf("Receive Queue (RQueue): %d\n", socket.InetDiagMsg.RQueue) // Количество байтов в очереди получения
		fmt.Printf("Send Queue (WQueue): %d\n", socket.InetDiagMsg.WQueue)    // Количество байтов в очереди отправки

		if userInfo != nil {
			fmt.Printf("UserID: %d\n", UserID)
			fmt.Printf("User: %s\n", userInfo.Username)
		} else {
			fmt.Println("User: Unknown")
		}

		fmt.Println()
	}
}
