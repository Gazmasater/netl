package main

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/sys/unix"
)

// Определяем состояния TCP вручную
const (
	TCP_ESTABLISHED     = 1
	TCP_SYN_SENT        = 2
	TCP_SYN_RECV        = 3
	TCP_FIN_WAIT1       = 4
	TCP_FIN_WAIT2       = 5
	TCP_TIME_WAIT       = 6
	TCP_CLOSE           = 7
	TCP_CLOSE_WAIT      = 8
	TCP_LAST_ACK        = 9
	TCP_LISTEN          = 10
	TCP_CLOSING         = 11
	SOCK_DIAG_BY_FAMILY = 4 // Определяем SOCK_DIAG_BY_FAMILY вручную
)

func main() {
	// Создаем netlink-сокет для SOCK_DIAG
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_SOCK_DIAG)
	if err != nil {
		log.Fatalf("Error creating netlink socket: %v", err)
	}
	defer unix.Close(fd)

	// Подготавливаем адрес для подписки на события TCP
	saddr := &unix.SockaddrNetlink{
		Family: unix.AF_NETLINK,
		Groups: unix.NETLINK_SOCK_DIAG,
	}

	// Привязываем сокет к адресу
	if err := unix.Bind(fd, saddr); err != nil {
		log.Fatalf("Error binding netlink socket: %v", err)
	}

	log.Println("Netlink socket successfully created and bound. Listening for socket events...")

	// Бесконечный цикл для получения сообщений
	for {
		buf := make([]byte, 4096)

		n, _, err := unix.Recvfrom(fd, buf, 0)
		if err != nil {
			log.Fatalf("Error receiving from netlink socket: %v", err)
		}

		// Разбираем полученные данные
		msgs, err := syscall.ParseNetlinkMessage(buf[:n])
		if err != nil {
			log.Fatalf("Error parsing netlink message: %v", err)
		}

		// Обрабатываем каждое сообщение
		for _, msg := range msgs {

			// Просто выводим тип сообщения
			fmt.Printf("Received netlink message type: %d\n", msg.Header.Type)
		}
	}
}
