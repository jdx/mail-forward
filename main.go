package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:25")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("connection received")
		bufin := bufio.NewReader(conn)
		bufout := bufio.NewWriter(conn)
		bufout.WriteString("220 mail.dickey.xxx\r\n")
		bufout.Flush()
		reply, err := bufin.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(reply)
		conn.Close()
	}
}
