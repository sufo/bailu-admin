/**
 *
 * @Desc https://github.com/casdoor/go-sms-sender
 * 短信发送
 */

package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/sufo/bailu-admin/app/config"
)

type AliyunClient struct {
	template string
	sign     string
	core     *dysmsapi.Client
}

type AliyunResult struct {
	RequestId string
	Message   string
}

// func New(accessId string, accessKey string, sign string, template string) (*AliyunClient, error) {
// func New(c config.SMS) (*AliyunClient, error) {
func New() (*AliyunClient, error) {
	c := config.Conf.SMS
	client, err := dysmsapi.NewClientWithAccessKey(c.RegionId, c.AppKey, c.AppSecret)
	if err != nil {
		return nil, err
	}

	aliyunClient := &AliyunClient{
		template: c.TemplateCode,
		core:     client,
		sign:     c.SignName,
	}

	return aliyunClient, nil
}

func (c *AliyunClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	requestParam, err := json.Marshal(param)
	if err != nil {
		return err
	}

	if len(targetPhoneNumber) < 1 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	phoneNumbers := bytes.Buffer{}
	phoneNumbers.WriteString(targetPhoneNumber[0])
	for _, s := range targetPhoneNumber[1:] {
		phoneNumbers.WriteString(",")
		phoneNumbers.WriteString(s)
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phoneNumbers.String()
	request.TemplateCode = c.template
	request.TemplateParam = string(requestParam)
	request.SignName = c.sign

	response, err := c.core.SendSms(request)
	if response.Code != "OK" {
		aliyunResult := AliyunResult{}
		_ = json.Unmarshal(response.GetHttpContentBytes(), &aliyunResult)
		return fmt.Errorf(aliyunResult.Message)
	}
	return err
}
