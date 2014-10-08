package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	addr := "0.0.0.0:" + port()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(&Client{
			Conn: conn,
			in:   bufio.NewReader(conn),
			out:  bufio.NewWriter(conn),
		})
	}
}

type Client struct {
	net.Conn
	in  *bufio.Reader
	out *bufio.Writer
}

func handleConn(c *Client) {
	defer c.Close()
	fmt.Println("connection received")
	c.writeline("220 mail.dickey.xxx")
	c.out.Flush()
	for {
		input := c.readline()
		cmd := strings.ToUpper(input)
		switch {
		case strings.Index(cmd, "EHLO") == 0:
			c.writeline("250-mail.dickey.xxx")
			c.writeline("250-SIZE 35882577")
			c.writeline("250-8BITMIME")
			c.writeline("250-ENHANCEDSTATUSCODES")
			c.writeline("250-SMTPUTF8")
		case strings.Index(cmd, "NOOP") == 0:
			c.writeline("250 OK")
		case strings.Index(cmd, "QUIT") == 0:
			c.writeline("221 Bye")
			return
		}
		c.out.Flush()
	}
}

func (c *Client) writeline(s string) {
	c.out.WriteString(s + "\r\n")
	fmt.Println("snt:", s)
}

func (c *Client) readline() string {
	line, err := c.in.ReadString('\n')
	if err != nil {
		panic(err)
	}
	line = strings.TrimSpace(line)
	fmt.Println("rcv:", line)
	return line
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "25"
	}
	return port
}
