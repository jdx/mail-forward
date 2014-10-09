package main

import (
	"fmt"
	"net/smtp"
)

func SendMail(mail *Mail) error {
	c, err := smtp.Dial("gmail-smtp-in.l.google.com:25")
	if err != nil {
		return err
	}
	fmt.Println("Sending email from", mail.From)
	if err := c.Mail(mail.From); err != nil {
		return err
	}
	for _, to := range mail.To {
		to = "jeff@dickeyxxx.com"
		fmt.Println("Sending email to", to)
		if err := c.Rcpt(to); err != nil {
			return err
		}
	}
	fmt.Println("a")
	wc, err := c.Data()
	if err != nil {
		return err
	}
	fmt.Println("a")
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		return err
	}
	err = wc.Close()
	fmt.Println("a")
	if err != nil {
		return err
	}
	fmt.Println("a")
	return c.Quit()
}
