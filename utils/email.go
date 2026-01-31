package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	From     string //发件人邮箱
	Password string // 发件人邮箱授权码
	Host     string // SMTP 服务器地址 ，例如 smtp.qq.com
	Port     int    // SMTP 端口， 例如 587
}

func SendEmailHTML(to, subject, htmlbody string, cfg EmailConfig) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// 建立 TLS 连接（关键）
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         cfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return fmt.Errorf("TLS 连接失败: %v", err)
	}

	// 创建 SMTP 客户端
	c, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP 客户端失败:%v", err)
	}
	defer c.Quit()

	//认证
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败:%v", err)
	}

	//设置发件人和收件人
	if err = c.Mail(cfg.From); err != nil {
		return fmt.Errorf("设置发件人失败 :%v", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败:%v", err)
	}

	// 写入邮件内容
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("创建邮件数据失败:%v", err)
	}

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html;charset=UTF-8\r\n\r\n%s",
		cfg.From, to, subject, htmlbody,
	)

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败:%v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭写入通道失败：%v", err)
	}

	return nil
}
