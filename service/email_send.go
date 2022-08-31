package service

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"ntfileupload/config"
)

// SendEmail 发送邮件 --> receiveMailbox: 接收方邮件 subject: 邮件主题 content: 邮件内容
func SendEmail(receiveMailbox string, subject string, content string) {
	// 简单设置 log 参数
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = config.SendEmail

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	//em.To = []string{"xxx@qq.com"}
	em.To = []string{receiveMailbox}

	// 设置主题
	em.Subject = subject

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte(content)

	//设置服务器相关的配置
	err := em.Send(config.EmailAddr, smtp.PlainAuth("", config.Username, config.Password, config.Host))
	if err != nil {
		log.Println("send fail!")
		log.Fatal(err)
	}
	log.Println("send successfully ... ")

}
