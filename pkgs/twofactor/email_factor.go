package twofactor

import (
	"fmt"
	"github.com/zouxinjiang/le/pkgs/cmessage"
	"strings"
)

type emailFactorTemplate string

func (e emailFactorTemplate) FillParams(params map[string]string) (content string) {
	content = "Content-Type: text/plain; charset=utf-8\r\nFrom: ${from}\r\nTo: ${to}\r\nSubject: LE 邮箱验证码\r\n\r\n[LE] 您此次验证码为:${code}"
	for k, v := range params {
		content = strings.Replace(content, fmt.Sprintf("${%s}", k), v, -1)
	}
	return content
}

type EmailFactor struct {
}

func init() {
	register("email", newEmailFactor)
}

func newEmailFactor() TwoFactor {
	return &EmailFactor{}
}

func (e EmailFactor) Do(params map[string]string) (addr string, code string, err error) {
	from := params["from"]
	to := params["to"]
	verifyCode := params["code"]
	msg := cmessage.New("email")
	msg.Initialize(params)
	msg.SetTemplate(new(emailFactorTemplate))
	return to, verifyCode, msg.Send(from, []string{to}, map[string]string{
		"code": verifyCode,
		"from": from,
		"to":   to,
	})
}
