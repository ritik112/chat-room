package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}

	//channels for incoming connections, dead connections and messages
	aconns := make(map[net.Conn]int)
	conns := make(chan net.Conn)
	dconns := make(chan net.Conn)
	msgs := make(chan string)
	i := 0

	go func() {

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
			}
			conns <- conn

		}
	}()

	for {
		select {
		//read incoming connections
			// go ritik go
		case conn := <-conns:
			aconns[conn] = i
			i++
			//connected, Read messages
			go func(conn net.Conn, i int) {
				rd := bufio.NewReader(conn)
				for {
					m, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					msgs <- fmt.Sprintf("Client %v:%v ", i, m)
				}
				// Done reading
				dconns <- conn
			}(conn, i)
		case msg := <-msgs:
			//Broadcast to all connections
			for conn := range aconns {
				conn.Write([]byte(msg))
			}

		case dconn := <-dconns:
			log.Printf("Client %v is Disconnected \n", aconns[dconn])
			delete(aconns, dconn)
		}

	}
}
