package mymail

import (
	"fmt"
	"net/smtp"
	"strconv"
)

type (
	ReqSmtpServer func() (smtpHost string, smtpPort int)
)

func SendByGmail(fromMail, fromPass, toMail, title, content string, smtpCb ReqSmtpServer, isHTML bool) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	if smtpCb != nil {
		smtpHost, smtpPort = smtpCb()
	}

	emailauth := smtp.PlainAuth("", fromMail, fromPass, smtpHost)
	receivers := []string{toMail}

	// Define MIME headers for email
	mimeType := "text/plain"
	if isHTML {
		mimeType = "text/html"
	}

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", toMail, title, mimeType, content)) // your message

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
