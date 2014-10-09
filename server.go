package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
	in   *bufio.Reader
	out  *bufio.Writer
}

type Mail struct {
	From string
	To   []string
}

var commands = map[string]func(c *Client, mail *Mail, input string) error{
	"HELO":       cmdHelo,
	"EHLO":       cmdEhlo,
	"STARTTLS":   cmdStartTLS,
	"NOOP":       cmdNoop,
	"MAIL FROM:": cmdMailFrom,
	"RCPT TO:":   cmdRcptTo,
	"DATA":       cmdData,
	"QUIT":       cmdQuit,
}

func handleConn(c *Client) {
	defer c.conn.Close()
	mail := &Mail{}
	c.writeline("220 mail.dickey.xxx")
	for {
		input, err := c.in.ReadString('\n')
		input = strings.TrimSpace(input)
		fmt.Println("c:", input)
		if err != nil {
			log.Println(err)
			return
		}
		err = runCommand(c, mail, input)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println(err)
			c.writeline("500 unexpected error")
			return
		}
	}
}

func runCommand(c *Client, mail *Mail, input string) error {
	upperInput := strings.ToUpper(input)
	for cmd, fn := range commands {
		if strings.Index(upperInput, cmd) == 0 {
			return fn(c, mail, input)
		}
	}
	return cmdUnknown(c, mail, input)
}

func cmdHelo(c *Client, mail *Mail, input string) error {
	c.writeline("220 mail.dickey.xxx")
	return nil
}

func cmdEhlo(c *Client, mail *Mail, input string) error {
	c.writeline("250-mail.dickey.xxx")
	c.writeline("250-SIZE 35882577")
	//c.writeline("250-STARTTLS")
	c.writeline("250-8BITMIME")
	c.writeline("250-ENHANCEDSTATUSCODES")
	c.writeline("250 SMTPUTF8")
	return nil
}

func cmdStartTLS(c *Client, mail *Mail, input string) error {
	c.writeline("220 Ready to start TLS")
	tlsConn := tls.Server(c.conn, TLSConfig)
	err := tlsConn.Handshake()
	if err != nil {
		return err
	}
	c.conn = tlsConn
	c.in = bufio.NewReader(c.conn)
	c.out = bufio.NewWriter(c.conn)
	return nil
}

func cmdNoop(c *Client, mail *Mail, input string) error {
	c.writeline("250 OK")
	return nil
}

func cmdMailFrom(c *Client, mail *Mail, input string) error {
	mail.From = input[8:]
	c.writeline("250 Accepted")
	return nil
}

func cmdRcptTo(c *Client, mail *Mail, input string) error {
	mail.To = append(mail.To, input)
	c.writeline("250 Accepted")
	return nil
}

func cmdData(c *Client, mail *Mail, input string) error {
	err := SendMail(mail)
	if err != nil {
		return err
	}
	c.writeline("354 End data with <CR><LF>.<CR><LF>")
	for {
		c.conn.SetDeadline(time.Now().Add(time.Minute))
		line, err := c.in.ReadString('\n')
		if err != nil {
			return err
		}
		if line == ".\r\n" {
			c.writeline("250 OK")
			return nil
		}
	}
}

func cmdQuit(c *Client, mail *Mail, input string) error {
	c.writeline("221 Bye")
	return io.EOF
}

func cmdUnknown(c *Client, mail *Mail, input string) error {
	c.writeline("500 Unrecognized command")
	log.Println("Unrecognized command:", input)
	return nil
}

func (c *Client) writeline(s string) {
	c.out.WriteString(s + "\r\n")
	c.out.Flush()
	fmt.Println("s:", s)
}
