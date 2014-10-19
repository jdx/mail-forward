package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"regexp"
	"strings"
)

type Client struct {
	Conn *textproto.Conn
	From string
	To   []string
}

var mailFromRegex = regexp.MustCompile("MAIL FROM: ?<(.*)>.*")
var rcptToRegex = regexp.MustCompile("RCPT TO: ?<(.*)>.*")

type cmdFn func(c *Client, input string) error

var commands = map[string]cmdFn{
	"HELO":       cmdHelo,
	"EHLO":       cmdEhlo,
	"STARTTLS":   cmdStartTLS,
	"NOOP":       cmdNoop,
	"MAIL FROM:": cmdMailFrom,
	"RCPT TO:":   cmdRcptTo,
	"DATA":       cmdData,
	"QUIT":       cmdQuit,
}

func handleConn(conn net.Conn) {
	c := &Client{Conn: textproto.NewConn(conn)}
	defer c.Conn.Close()
	c.PrintfLine("220 mx.grandcentralemail.com")
	for {
		input, err := c.Conn.ReadLine()
		fmt.Println("r:", input)
		if err != nil {
			log.Println(err)
			return
		}
		f := parseCmd(input)
		if err := f(c, input); err != nil {
			if err == io.EOF {
				return
			}
			log.Println(err)
			c.PrintfLine("500 unexpected error")
			return
		}
	}
}

func parseCmd(input string) cmdFn {
	upper := strings.ToUpper(input)
	for cmd, fn := range commands {
		if strings.Index(upper, cmd) == 0 {
			return fn
		}
	}
	return cmdUnknown
}

func cmdHelo(c *Client, input string) error {
	return c.PrintfLine("220 mx.grandcentralemail.com")
}

func cmdEhlo(c *Client, input string) error {
	c.PrintfLine("250-mx.grandcentralemail.com")
	c.PrintfLine("250-SIZE 35882577")
	//c.PrintfLine("250-STARTTLS")
	c.PrintfLine("250-8BITMIME")
	return c.PrintfLine("250 SMTPUTF8")
}

func cmdStartTLS(c *Client, input string) error {
	//c.PrintfLine("220 Ready to start TLS")
	//tlsConn := tls.Server(c, TLSConfig)
	//err := tlsConn.Handshake()
	//if err != nil {
	//return err
	//}
	//c = tlsConn
	//c.in = bufio.NewReader(c)
	//c.out = bufio.NewWriter(c)
	return nil
}

func cmdNoop(c *Client, input string) error {
	return c.PrintfLine("250 OK")
}

func cmdMailFrom(c *Client, input string) error {
	c.From = mailFromRegex.FindStringSubmatch(input)[1]
	return c.PrintfLine("250 Accepted")
}

func cmdRcptTo(c *Client, input string) error {
	address := rcptToRegex.FindStringSubmatch(input)[1]
	if !strings.HasSuffix(address, "dickeyxxx.com") {
		return c.PrintfLine("500 Invalid email")
	}
	c.To = append(c.To, address)
	return c.PrintfLine("250 Accepted")
}

func cmdData(c *Client, input string) error {
	c.PrintfLine("354 End data with <CR><LF>.<CR><LF>")
	lines, err := c.Conn.ReadDotLines()
	if err != nil {
		return err
	}
	fmt.Printf("r:\n  %s\n", strings.Join(lines, "\n  "))
	if err := SendMail(c.From, c.To, lines); err != nil {
		return err
	}
	return c.PrintfLine("250 OK")
}

func cmdQuit(c *Client, input string) error {
	c.PrintfLine("221 Bye")
	return io.EOF
}

func cmdUnknown(c *Client, input string) error {
	log.Println("Unrecognized command:", input)
	return c.PrintfLine("500 Unrecognized command")
}

func (c *Client) PrintfLine(format string, args ...interface{}) error {
	fmt.Printf("s: "+format+"\n", args...)
	return c.Conn.PrintfLine(format, args...)
}
