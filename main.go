package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var TLSConfig *tls.Config

func main() {
	TLSConfig = setupTLS()
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

func setupTLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair("./certs/ssl.pem", "./certs/ssl.key")
	if err != nil {
		log.Println("Error loading certificate:", err)
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ServerName:   "mail.dickey.xxx",
	}
	config.Rand = rand.Reader
	return config
}

type Client struct {
	net.Conn
	in       *bufio.Reader
	out      *bufio.Writer
	mailFrom string
	mailTo   []string
}

func handleConn(c *Client) {
	defer c.Close()
	fmt.Println("connection received")
	c.writeline("220 mail.dickey.xxx")
	for {
		input := c.readline()
		cmd := strings.ToUpper(input)
		switch {
		case strings.Index(cmd, "HELO") == 0:
			c.writeline("220 mail.dickey.xxx SMTP")
		case strings.Index(cmd, "EHLO") == 0:
			c.writeline("250-mail.dickey.xxx")
			c.writeline("250-SIZE 35882577")
			//c.writeline("250-STARTTLS")
			c.writeline("250-8BITMIME")
			c.writeline("250-ENHANCEDSTATUSCODES")
			c.writeline("250 SMTPUTF8")
		case strings.Index(cmd, "STARTTLS") == 0:
			c.writeline("220 Ready to start TLS")
			c.upgradeTLS()
		case strings.Index(cmd, "NOOP") == 0:
			c.writeline("250 OK")
		case strings.Index(cmd, "MAIL FROM:") == 0:
			c.mailFrom = input[8:]
			c.writeline("250 Accepted")
		case strings.Index(cmd, "RCPT TO:") == 0:
			c.mailTo = append(c.mailTo, input[8:])
			c.writeline("250 Accepted")
		case strings.Index(cmd, "DATA") == 0:
			c.writeline("354 End data with <CR><LF>.<CR><LF>")
			err := c.readdata()
			if err != nil {
				log.Println(err)
				c.writeline("500 unexpected error")
				return
			}
			c.writeline("250 OK: Queued as 298892")
		case strings.Index(cmd, "QUIT") == 0:
			c.writeline("221 Bye")
			return
		default:
			c.writeline("500 unrecognized command")
			log.Println("Unrecognized:", input)
		}
	}
}

func (c *Client) writeline(s string) {
	c.out.WriteString(s + "\r\n")
	c.out.Flush()
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

func (c *Client) readdata() error {
	for {
		c.SetDeadline(time.Now().Add(time.Minute))
		line, err := c.in.ReadString('\n')
		if err != nil {
			return err
		}
		if line == ".\r\n" {
			return nil
		}
	}
}

func (c *Client) upgradeTLS() {
	tlsConn := tls.Server(c, TLSConfig)
	err := tlsConn.Handshake()
	if err != nil {
		panic(err)
	}
	c.Conn = tlsConn
	c.in = bufio.NewReader(c.Conn)
	c.out = bufio.NewWriter(c.Conn)
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "25"
	}
	return port
}
