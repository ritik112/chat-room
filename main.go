package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"strconv"
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
	var users [25]string
	var rights [25]string
	// users := make(chan string)
	i := 0

	go func() {

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
			}
			conn.Write([]byte("Enter your name in format \nusername yourname\n\n"))
			conns <- conn
		}
	}()

	for {
		select {
		//read incoming connections
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
					if strings.Contains(m, "username"){
						words := strings.Fields(m)
						userlist := "Userlist\n"
						if len(words)>2 {
							users[i] = words[2]
							rights[i] = words[1]
							for j, user := range users {
								if !(user ==""){
									userlist +=fmt.Sprintf("%v - %v \n", j, user)
								}
							}
							userlist +=fmt.Sprintf("Chat to particular user write\nconectto userid 'your message\n\nGroup chat\ngroupmsg no. of user userids your msg\n\n")
						}
						msgs <- fmt.Sprintf(userlist)
					}else if (strings.Contains(m, "conectto") && !(users[i]=="") && strings.Compare(rights[i],"all")==0){
						words := strings.Fields(m)
						if len(words)>2 {
							for usercon,k := range aconns {
								_, err := strconv.ParseInt(words[1], 0, 64)
								if err == nil {
									k2 := strconv.Itoa(k+1)
									fmt.Println(k2)
									if !(strings.Compare(words[1],k2)==0){
										fmt.Println("false")
									}else if (users[k+1]=="") {
										usermsg := fmt.Sprintf("Client %v doesn't exist \n\n", (k+1))
										conn.Write([]byte(usermsg))
									}else{
										messg :="" 
										for msg := 2;  msg<=len(words)-1; msg++ {
											messg += " "+words[msg]
										}
										usermsg := fmt.Sprintf("Client %v:%v \n\n", i, messg)
										usercon.Write([]byte(usermsg))
									}

								}else{
									usermsg := fmt.Sprintf("Invalid user id\n\n")
									conn.Write([]byte(usermsg))
								}								
							}
						}
					}else if (strings.Contains(m, "groupmsg") && !(users[i]=="") && strings.Compare(rights[i],"all")==0){
						words := strings.Fields(m)
						if len(words)>3 {
							g1, err := strconv.ParseInt(words[1],10,64)
							fmt.Println(g1)
							if err == nil {
								for m := int64(1);  m<=g1; m++ {
									for usercon,k := range aconns {
										_, err := strconv.ParseInt(words[m+1], 0, 64)
										if err == nil {
											k2 := strconv.Itoa(k+1)
											fmt.Println(k2,k+1,strings.Compare(words[m+1],k2))
											if !(strings.Compare(words[m+1],k2)==0){
												fmt.Println("false")
											}else if (users[k+1]=="") {
												usermsg := fmt.Sprintf("Client %v doesn't exist \n\n", (k+1))
												conn.Write([]byte(usermsg))
											}else{
												messg :="" 
												for msg := g1+2;  msg<=int64(len(words)-1); msg++ {
													messg += " "+words[msg]
												}
												usermsg := fmt.Sprintf("Client %v:%v \n\n", i, messg)
												usercon.Write([]byte(usermsg))
											}
										}								
									}
									fmt.Printf("Welcome %d times\n",i)
								}

							}else{
								fmt.Println("Invalid number")
							}
						}
					}else if (!(users[i]=="") && strings.Compare(rights[i],"all")==0){
						msgs <- fmt.Sprintf("Client %v:%v ", i, m)
					}else if !(strings.Compare(rights[i],"all")==0){
						conn.Write([]byte("You denied\n"))
					}else{
						conn.Write([]byte("Please first write your username\n"))
					}
				}
				// Done reading
				dconns <- conn
			}(conn, i)
		case msg := <-msgs:
			//Broadcast to all connections
			for conn,i := range aconns {
				if !(users[i+1]==""){
					conn.Write([]byte(msg))
				}
				
			}

		case dconn := <-dconns:
			log.Printf("Client %v is Disconnected \n", aconns[dconn])
			delete(aconns, dconn)
		}

	}
}
