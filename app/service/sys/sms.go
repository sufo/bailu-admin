/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc
 */

package sys

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"math/rand"
	"bailu/app/config"
	"bailu/pkg/sms"
	"bailu/pkg/store"
	"time"
)

var SMSSet = wire.NewSet(wire.Struct(new(SMSService), "*"))

type SMSService struct {
	Store store.IStore
	Sms   *sms.AliyunClient
}

// 发送短信
func (s *SMSService) SendSMS(ctx context.Context, phone string, dialCode string) error {
	//随机生成6位数验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	params := map[string]string{}
	params["code"] = code
	number := dialCode + phone
	err := s.Sms.SendMessage(params, number)
	if err == nil {
		expired := config.Conf.SMS.Expired
		if e := s.Store.Set(number, code, time.Duration(expired)); e != nil {
			_ = s.Store.Set(number, code, time.Duration(expired)) //重新插入一次
		}
	}
	return err
}
