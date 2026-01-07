/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package mail

import (
	"errors"
	"github.com/sufo/bailu-admin/pkg/log"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

type Options struct {
	MailHost string
	MailPort int
	MailUser string // 发件人
	MailPass string // 发件人密码
	MailTo   string // 收件人 多个用,分割
	Subject  string // 邮件主题
	Body     string // 邮件内容
}

func Send(o *Options) error {
	var err error
	if err = validate(o); err != nil {
		log.L.Error(err)
		return err
	}

	m := gomail.NewMessage()

	//设置发件人
	m.SetHeader("From", o.MailUser)

	//设置发送给多个用户
	mailArrTo := strings.Split(o.MailTo, ",")
	m.SetHeader("To", mailArrTo...)

	//设置邮件主题
	m.SetHeader("Subject", o.Subject)

	//设置邮件正文
	m.SetBody("text/html", o.Body)

	d := gomail.NewDialer(o.MailHost, o.MailPort, o.MailUser, o.MailPass)

	//return d.DialAndSend(m)
	//发送（失败最多三次重发）
	maxTimes := 3
	i := 0
	for i < maxTimes {
		err = d.DialAndSend(m)
		if err == nil {
			break
		}
		i += 1
		time.Sleep(2 * time.Second)
		if i < maxTimes {
			log.L.Errorf("mail#发送消息失败#%s#消息内容-%s", err.Error(), o.Body)
		}
	}
	return err
}

func validate(o *Options) error {
	if o.MailHost == "" {
		return errors.New("#mail#Host为空")
	}
	if o.MailPort == 0 {
		return errors.New("#mail#Port为空")
	}
	if o.MailUser == "" {
		return errors.New("#mail#User为空")
	}
	if o.MailPass == "" {
		log.L.Error("#mail#Password为空")
		return errors.New("#mail#Password为空")
	}
	if o.MailTo == "" {
		return errors.New("#mail#To为空")
	}
	return nil
}
