package ReverseProxy

import (
	"fmt"
	"net"
	"strconv"
)

type UpstreamBridge struct {
	address string
	port    int
	conn    *net.TCPConn
}

func Start(remoteAddress string, remotePort int) {

	address, err := net.ResolveTCPAddr("tcp", "localhost:"+strconv.Itoa(25565))
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		connection.SetNoDelay(true)
		connection.SetKeepAlive(true)

		go handleConnection(connection, UpstreamBridge{
			address: remoteAddress,
			port:    remotePort,
		})
	}
}

func handleConnection(conn *net.TCPConn, upstream UpstreamBridge) {
	fmt.Println("Handle new connection:", conn.RemoteAddr())
	err := upstream.openConnection()
	if err != nil {
		conn.Close()
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to upstream")

	go func() {
		for {
			buf := make([]byte, 1024)
			size, err := upstream.conn.Read(buf)
			if err != nil {
				fmt.Println("Closed connection:", err.Error())
				break
			}
			buf = buf[:size]
			if size > 0 {
				fmt.Println("Write to client")
				conn.Write(buf)
			}
		}
	}()

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Closed connection:", err.Error())
			break
		}
		buf = buf[:size]
		if size > 0 {
			fmt.Println("Write to server")
			upstream.conn.Write(buf)
		}
	}
}

func (upstream *UpstreamBridge) openConnection() error {
	address, err := net.ResolveTCPAddr("tcp", upstream.address+":"+strconv.Itoa(upstream.port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		return err
	}
	conn.SetNoDelay(true)
	conn.SetKeepAlive(true)
	upstream.conn = conn
	return nil
}
