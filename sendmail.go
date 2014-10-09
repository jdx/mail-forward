package main

import (
	"fmt"
	"net/smtp"
)

func SendMail(mail *Mail) error {
	c, err := smtp.Dial("smtp.google.com:25")
	if err != nil {
		return err
	}
	if err := c.Mail(mail.From); err != nil {
		return err
	}
	for _, to := range mail.To {
		if err := c.Rcpt(to); err != nil {
			return err
		}
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(wc, "This is the email body")
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
