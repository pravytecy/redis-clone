package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/pravytecy/redis-clone/config"
)

func RunTcpSyncServer() {
	log.Println("starting a TCP sync server", config.Host, config.Port)
	var con_clients int = 0
	lstnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := lstnr.Accept()
		if err != nil {
			panic(err)
		}
		con_clients += 1
		log.Println("client connected with address", conn.RemoteAddr(), "no of clients", con_clients)

		for {
			c, err := readCommand(conn)
			if err != nil {
				conn.Close()
				con_clients -= 1
				log.Println("client got disconnected", conn.RemoteAddr(), "no of clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			if err := respond(c, conn); err != nil {
				log.Print("err write", err)
			}
		}
	}

}

func respond(c string, conn net.Conn) error {
	_, err := conn.Write([]byte(c))
	if err != nil {
		return err
	}
	return nil
}

func readCommand(conn net.Conn) (string, error) {
	var b []byte = make([]byte, 512)
	n, err := conn.Read(b[:])
	if err != nil {
		return "", err
	}
	return string(b[:n]), nil
}
