/**
 * @copyright www.ruomm.com
 * @author 牛牛-wanruome@126.com
 * @create 2024/5/31 16:52
 * @version 1.0
 */
package aliyunx

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type AliyunEmailClient struct {
	AccessKeyId     string `yaml:"accessKeyId" xref:"AccessKeyId;tidy" validate:"min=1,max=50" xvalid_error:"应用秘钥编号必须填写"`
	AccessKeySecret string `yaml:"accessKeySecret" xref:"AccessKeySecret;tidy" validate:"min=1,max=100" xvalid_error:"应用秘钥信息必须填写"`
	EmailAccount    string `yaml:"emailAccount" xref:"EmailAccount;tidy" validate:"email,min=1,max=50" xvalid_error:"邮箱账户必须填写"`
	EmailAlias      string `yaml:"emailAlias" xref:"EmailAlias;tidy" validate:"omitempty,min=1,max=50" xvalid_error:"发件人别名填写错误"`
	EmailSubject    string `yaml:"emailSubject" xref:"EmailSubject;tidy" validate:"min=1,max=50" xvalid_error:"邮件主题填写错误"`
}

type EmailInfo struct {
	EmailAccount string // 邮箱账户必须填写，没有则取配置里面的参数
	EmailAlias   string // 发件人别名必须填写，没有则取配置里面的参数
	EmailSubject string // 邮件主题填写必须填写，没有则取配置里面的参数
	ToAddress    string // 送达地址必须填写
	HtmlBody     string // 邮件内容必须填写
}

func (u *AliyunEmailClient) SendEmail(info EmailInfo) error {
	cfg := &openapi.Config{
		AccessKeyId:     &u.AccessKeyId,
		AccessKeySecret: &u.AccessKeySecret,
	}
	cfg.Endpoint = tea.String("dm.aliyuncs.com")
	client, err := dm.NewClient(cfg)
	if err != nil {
		return err
	}

	emailAccount, err := u.verifyParam(info.EmailAccount, u.EmailAccount, "邮箱账户")
	if err != nil {
		return err
	}
	emailAlias, err := u.verifyParam(info.EmailAlias, u.EmailAlias, "发件人别名")
	if err != nil {
		return err
	}
	emailSubject, err := u.verifyParam(info.EmailSubject, u.EmailSubject, "邮件主题")
	if err != nil {
		return err
	}
	if len(info.ToAddress) <= 0 {
		return errors.New(fmt.Sprintf("阿里云邮件发送失败，%s参数必须填写！", "送达地址"))
	}
	if len(info.HtmlBody) <= 0 {
		return errors.New(fmt.Sprintf("阿里云邮件发送失败，%s参数必须填写！", "邮件内容"))
	}
	param := &dm.SingleSendMailRequest{
		AccountName:    &emailAccount,
		AddressType:    tea.Int32(1),
		ReplyToAddress: tea.Bool(false),
		Subject:        &emailSubject,
		ToAddress:      &info.ToAddress,
		FromAlias:      &emailAlias,
		HtmlBody:       &info.HtmlBody,
	}
	runtime := &util.RuntimeOptions{}
	_, err = client.SingleSendMailWithOptions(param, runtime)
	if err != nil {
		return err
	}
	return nil
}

func (u *AliyunEmailClient) verifyParam(paramInfo string, paramConfig string, errTag string) (string, error) {
	if len(paramInfo) > 0 {
		return paramInfo, nil
	}
	if len(paramConfig) > 0 {
		return paramConfig, nil
	}
	return "", errors.New(fmt.Sprintf("阿里云邮件发送失败，%s参数必须填写！", errTag))
}
