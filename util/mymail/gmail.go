package mymail

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type (
	ReqSmtpServer func() (smtpHost string, smtpPort int)
)

func SendByGmail(fromMail, fromPass, toMail, title, content string, smtpCb ReqSmtpServer) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	if smtpCb != nil {
		smtpHost, smtpPort = smtpCb()
	}

	emailauth := smtp.PlainAuth("", fromMail, fromPass, smtpHost)
	receivers := []string{toMail}

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", toMail, title, content)) // your message

	// send out the email
	if err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), //convert port number from int to string
		emailauth,
		fromMail,
		receivers,
		message,
	); err != nil {
		return err
	}

	return nil
}
