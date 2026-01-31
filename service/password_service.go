package service

import (
	"bookstore/dao"
	"bookstore/utils"
	"errors"
	"fmt"
	"time"
)

// SendResetCodeService 发送邮箱验证码(支持用户/管理员)
func SendResetCodeService(email string) error {

	// 生成验证码
	code := utils.GenerateVerifyCode(6)

	if err := dao.SaveResetCode(email, code, 5*time.Minute); err != nil {
		return errors.New("保存验证码失败:" + err.Error())
	}
	// 组织邮件内容(HTML 格式)
	htmlBody := fmt.Sprintf(`
	<h2>找回密码</h2>
	<p>您好，您正在尝试找回密码。</p>
	<p>您的验证码是: <b style="font-size:20px;">%s</b> </p>
	<p>验证码有效期为 <b>5分钟</b>，请尽快使用。 </p>
`, code)

	// 发送邮件
	err := utils.SendEmailHTML(email, "找回密码验证码", htmlBody, utils.EmailConfig{
		From:     "wangkxin@foxmail.com",
		Password: "qkoxaullrqhtebfc",
		Host:     "smtp.qq.com",
		Port:     465,
	})
	if err != nil {
		return errors.New("验证码邮件发送失败:" + err.Error())
	}
	return nil
}

// RecoverPasswordService 用户提交验证码 + 新密码
func RecoverPasswordService(userType string, email string, code string, newPassword string) error {
	// 验证验证码
	storedCode, err := dao.GetResetCode(email)
	if err != nil || storedCode != code {
		return errors.New("验证码错误或已过期")
	}

	// 密码加密
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	switch userType {
	case "user":
		err = dao.UpdateUserPasswordByEmail(email, hashed)
	case "admin":
		err = dao.UpdatedAdminPasswordByEmail(email, hashed)
	default:
		return errors.New("未知用户类型")
	}
	if err != nil {
		return errors.New("更新密码失败")
	}

	// 删除验证码
	_ = dao.DeleteResetCode(email)
	return nil
}
