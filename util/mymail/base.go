package mymail

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/any-call/gobase/util/mylog"
	"log"
	"net"
	"net/smtp"
	"strings"
)

// 端口 25：用于 SMTP 服务器之间的通信，通常未加密。
// 端口 465：用于 SSL/TLS 加密的 SMTP（旧的标准）
// 端口 587：用于客户端和服务器之间通过 STARTTLS 加密的 SMTP。

var (
	MailPort []int = []int{25, 465, 587}
)

func GetMXRecords(domain string) ([]string, error) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get MX records for domain %s: %v", domain, err)
	}

	var servers []string
	for _, mx := range mxRecords {
		servers = append(servers, mx.Host)
	}
	return servers, nil
}

// 检测 SMTP 服务器支持的端口和加密方式
func DetectSMTPPortsAndEncryption(server string, port int) (encryFlag bool, err error) {
	address := fmt.Sprintf("%s:%d", server, port)

	// 建立 TCP 连接
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 读取服务器的欢迎信息
	response, _ := reader.ReadString('\n')

	// 发送 EHLO 命令
	sendCommand(writer, "EHLO localhost")
	response = readResponse(reader)

	// 检查是否支持 STARTTLS
	if strings.Contains(response, "STARTTLS") {
		sendCommand(writer, "STARTTLS")
		response = readResponse(reader)
		if strings.HasPrefix(response, "220") {
			// 开始加密
			tlsConn := tls.Client(conn, &tls.Config{
				InsecureSkipVerify: true, // 忽略证书验证，仅用于测试
				ServerName:         server,
			})
			defer func() {
				_ = tlsConn.Close()
			}()

			// 切换到加密后的连接
			reader = bufio.NewReader(tlsConn)
			writer = bufio.NewWriter(tlsConn)
			sendCommand(writer, "EHLO localhost")
			response = readResponse(reader)
			encryFlag = true
		} else {
			encryFlag = true
			err = fmt.Errorf("Failed to upgrade to TLS")
		}
	} else {
		mylog.Info("Server does not support STARTTLS")
		encryFlag = false
	}

	// 发送 QUIT 命令
	sendCommand(writer, "QUIT")
	return encryFlag, err
}

func SendEmail(from, to, subject, body string, smtpPort int) error {
	if smtpPort == 0 {
		smtpPort = 25 //default
	}
	// 从收件人邮箱地址获取域名部分
	recipientDomain := strings.Split(to, "@")[1]

	// 获取目标邮件服务器
	mxRecords, err := GetMXRecords(recipientDomain)
	if err != nil {
		return fmt.Errorf("Error getting MX records: %v", err)
	}

	// 获取第一个 MX 记录
	smtpServer := mxRecords[0]
	mylog.Infof("Sending email via SMTP server: %s", smtpServer)

	// 配置邮件内容
	msg := []byte("Subject: " + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"\r\n" +
		body)

	// 使用 PlainAuth 创建身份验证
	//auth := smtp.PlainAuth("", from, "your-email-password", smtpServer)

	// 发送邮件
	err = smtp.SendMail(smtpServer+fmt.Sprintf(":%d", smtpPort), nil, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("Failed to send email: %v", err)
	}

	return nil
}

// 发送 SMTP 命令
func sendCommand(writer *bufio.Writer, cmd string) {
	fmt.Fprintf(writer, "%s\r\n", cmd)
	writer.Flush()
}

// 读取 SMTP 响应
func readResponse(reader *bufio.Reader) string {
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	return response
}
