package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	SOCK_DIAG_BY_FAMILY = 4 // Определяем SOCK_DIAG_BY_FAMILY вручную
	AF_INET             = 2 // IPv4
)

type NetlinkRequest struct {
	Header unix.NlMsghdr
	Data   struct {
		Family int8
		Pad    [3]byte
	}
}

type TcpInfo struct {
	State uint8 // Состояние сокета
	// Добавьте другие поля, которые вы хотите использовать
}

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

	log.Println("Netlink socket successfully created and bound. Sending request for TCP socket info...")

	// Подготовка запроса для получения информации о TCP-сокетах
	req := NetlinkRequest{
		Header: unix.NlMsghdr{
			Type:  SOCK_DIAG_BY_FAMILY,
			Flags: unix.NLM_F_REQUEST | unix.NLM_F_DUMP,
			Seq:   1,
			Pid:   uint32(unix.Getpid()),
		},
	}
	req.Data.Family = AF_INET // Запрашиваем информацию о IPv4-сокетах

	// Устанавливаем длину заголовка
	req.Header.Len = uint32(unsafe.Sizeof(req.Header) + unsafe.Sizeof(req.Data))

	// Создаем срез байтов для отправки запроса
	buf := make([]byte, req.Header.Len)
	*(*NetlinkRequest)(unsafe.Pointer(&buf[0])) = req

	// Отправляем запрос
	if err := unix.Sendto(fd, buf, 0, saddr); err != nil {
		log.Fatalf("Error sending netlink request: %v", err)
	}

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
			if msg.Header.Type == unix.NLMSG_DONE {
				log.Println("End of messages")
				return
			}

			// Извлекаем информацию о сокете
			if len(msg.Data) >= syscall.SizeofSockaddrInet4 {
				tcpInfo := TcpInfo{}
				// Заполняем tcpInfo из данных сообщения
				copy((*[1 << 20]byte)(unsafe.Pointer(&tcpInfo))[:], msg.Data[unsafe.Sizeof(unix.NlMsghdr{}):])
				state := tcpInfo.State

				// Проверяем состояние и выводим только активные сокеты
				if state != 1 { // Пример: состояние "1" может быть TCP_CLOSE
					fmt.Printf("Active TCP socket with state: %d\n", state)
				}
			}
		}
	}
}
