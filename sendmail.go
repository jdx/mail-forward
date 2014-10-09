package main

import (
	"fmt"
	"io"
	"net/smtp"
)

func SendMail(from string, to []string, data io.ReadCloser) <-chan error {
	done := make(chan error)
	go func() {
		c, err := smtp.Dial("gmail-smtp-in.l.google.com:25")
		if err != nil {
			done <- err
			return
		}
		fmt.Println("Sending email from", from)
		if err := c.Mail(from); err != nil {
			done <- err
			return
		}
		for _, to := range to {
			to = "dickeyxxx@gmail.com"
			fmt.Println("Sending email to", to)
			if err := c.Rcpt(to); err != nil {
				done <- err
				return
			}
		}
		wc, err := c.Data()
		if err != nil {
			done <- err
			return
		}
		_, err = fmt.Fprintf(wc, "This is the email body")
		if err != nil {
			done <- err
			return
		}
		err = wc.Close()
		if err != nil {
			done <- err
			return
		}
		err = c.Quit()
		if err != nil {
			done <- err
		}
		close(done)
	}()
	return done
}
