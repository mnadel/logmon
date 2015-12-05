package logmon

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"time"
)

type Alerter interface {
	SendAlert(map[string][]string) error
}

type EmailAlerter struct {
	Smtp    string
	From    string
	To      string
	Subject string
}

func (e *EmailAlerter) SendAlert(errors map[string][]string) error {
	log.Println("dialing:", e.Smtp)

	conn, err := net.DialTimeout("tcp", e.Smtp, 15*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, e.Smtp)
	if err != nil {
		return err
	}
	defer c.Close()

	c.Mail(e.From)
	c.Rcpt(e.To)

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	buf := bytes.NewBufferString(fmt.Sprintf("%v", errors))
	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}

	return nil
}
