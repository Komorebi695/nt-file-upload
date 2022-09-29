package service

import (
	"github.com/jordan-wright/email"
	"gopkg.in/gomail.v2"
	"log"
	"net/smtp"
	"ntfileupload/config"
	"strings"
)

// SendEmail 发送邮件 --> receiveMailbox: 接收方邮件 subject: 邮件主题 content: 邮件内容
func SendEmail(receiveMailbox string, subject string, content string) (err error) {
	// 简单设置 log 参数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = config.SendEmail

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{receiveMailbox}

	// 设置主题
	em.Subject = subject

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte(content)

	//设置服务器相关的配置
	err = em.Send("smtp.qq.com:25", smtp.PlainAuth("", config.Username, config.Password, config.Host))
	if err != nil {
		log.Println("send fail!")
		return err
	}
	log.Println("send successfully ... ")
	return nil
}

func SendEmail2(receiveMailbox string, subject string, content string) (err error) {
	m := gomail.NewMessage()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	m.SetHeader(`From`, config.SendEmail)
	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	m.SetHeader(`To`, receiveMailbox)
	// 设置主题
	m.SetHeader(`Subject`, subject)
	// 简单设置文件发送的内容
	m.SetBody(`text/html`, content)

	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	// Send
	if err := d.DialAndSend(m); err != nil {
		//panic(err)
		return err
	}
	log.Println("send successfully ... ")
	return nil
}

// SendToMail 主题 内容 类型 回调地址 收件人地址 抄送地址 密送地址
func SendToMail(subject, body, mailtype, replyToAddress string, to, cc, bcc []string) error {
	hp := strings.Split(config.AliyunHost, ":")
	auth := smtp.PlainAuth("", config.AliyunUser, config.AliyunPassword, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	ccAddress := strings.Join(cc, ";")
	bccAddress := strings.Join(bcc, ";")
	toAddress := strings.Join(to, ";")
	msg := []byte("To: " + toAddress + "\r\n" +
		"From: " + config.AliyunUser + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Reply-To: " + replyToAddress + "\r\n" +
		"Cc: " + ccAddress + "\r\n" +
		"Bcc: " + bccAddress + "\r\n" +
		contentType + "\r\n\r\n" +
		body)

	sendTo := MergeSlice(to, cc)
	sendTo = MergeSlice(sendTo, bcc)
	err := smtp.SendMail(config.AliyunHost, auth, config.AliyunUser, sendTo, msg)
	return err
}

func MergeSlice(s1 []string, s2 []string) []string {
	slice := make([]string, len(s1)+len(s2))
	copy(slice, s1)
	copy(slice[len(s1):], s2)
	return slice
}
