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

type Mail struct {
	From string
	To   []string
}

var addressRegex = regexp.MustCompile("MAIL FROM: ?<(.*)>.*")

type cmdFn func(c *textproto.Conn, mail *Mail, input string) error

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
	c := textproto.NewConn(conn)
	defer c.Close()
	mail := &Mail{}
	c.PrintfLine("220 mail.dickey.xxx")
	for {
		input, err := c.ReadLine()
		fmt.Println("c:", input)
		if err != nil {
			log.Println(err)
			return
		}
		f := parseCmd(input)
		if err := f(c, mail, input); err != nil {
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

func cmdHelo(c *textproto.Conn, mail *Mail, input string) error {
	c.PrintfLine("220 mail.dickey.xxx")
	return nil
}

func cmdEhlo(c *textproto.Conn, mail *Mail, input string) error {
	c.PrintfLine("250-mail.dickey.xxx")
	c.PrintfLine("250-SIZE 35882577")
	//c.PrintfLine("250-STARTTLS")
	c.PrintfLine("250-8BITMIME")
	c.PrintfLine("250 SMTPUTF8")
	return nil
}

func cmdStartTLS(c *textproto.Conn, mail *Mail, input string) error {
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

func cmdNoop(c *textproto.Conn, mail *Mail, input string) error {
	c.PrintfLine("250 OK")
	return nil
}

func cmdMailFrom(c *textproto.Conn, mail *Mail, input string) error {
	mail.From = addressRegex.FindStringSubmatch(input)[1]
	c.PrintfLine("250 Accepted")
	return nil
}

func cmdRcptTo(c *textproto.Conn, mail *Mail, input string) error {
	mail.To = append(mail.To, input)
	c.PrintfLine("250 Accepted")
	return nil
}

func cmdData(c *textproto.Conn, mail *Mail, input string) error {
	c.PrintfLine("354 End data with <CR><LF>.<CR><LF>")
	lines, err := c.ReadDotLines()
	if err != nil {
		return err
	}
	fmt.Printf("email received:\n %s\n", lines)
	if err := SendMail(mail.From, mail.To, lines); err != nil {
		return err
	}
	return c.PrintfLine("250 OK")
}

func cmdQuit(c *textproto.Conn, mail *Mail, input string) error {
	c.PrintfLine("221 Bye")
	return io.EOF
}

func cmdUnknown(c *textproto.Conn, mail *Mail, input string) error {
	c.PrintfLine("500 Unrecognized command")
	log.Println("Unrecognized command:", input)
	return nil
}
