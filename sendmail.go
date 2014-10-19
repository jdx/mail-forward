package main

import (
	"fmt"

	"github.com/dickeyxxx/mail-forward/smtpclient"
)

func SendMail(from string, to []string, lines []string) error {
	c, err := smtpclient.Dial("alt1.gmail-smtp-in.l.google.com:25")
	defer c.Close()
	if err != nil {
		return err
	}
	fmt.Println("Sending email from", from)
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, to := range to {
		to = "dickeyxxx@gmail.com"
		fmt.Println("Sending email to", to)
		if err := c.Rcpt(to); err != nil {
			return err
		}
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	for _, line := range lines {
		if _, err = fmt.Fprintf(wc, line); err != nil {
			return err
		}
	}
	wc.Close()
	return c.Quit()
}
