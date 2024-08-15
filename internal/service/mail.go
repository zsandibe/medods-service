package service

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

func (s *service) NotifyToEmail(email string, oldIp, newIp string) error {
	message := fmt.Sprintf("We noticed that you attempted to update your access token from a new IP address. Your current IP address is [%v], and the IP address from which you tried to perform the update is: [%v]", oldIp, newIp)
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.conf.Smtp.Username)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Attention!")
	msg.SetBody("text/plain", message)

	port, err := strconv.Atoi(s.conf.Smtp.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(s.conf.Smtp.Server, port, s.conf.Smtp.Username, s.conf.Smtp.Password)
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}