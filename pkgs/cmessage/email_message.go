package cmessage

import (
	"errors"
	"fmt"
	"github.com/zouxinjiang/le/pkgs/clog"
	"net/smtp"
	"strconv"
)

type EmailMessage struct {
	username   string
	host       string
	port       int
	initResult error
	auth       smtp.Auth
	tmplate    CMessageTemplate
}

func init() {
	register("email", func() CMessage {
		return &EmailMessage{}
	})
}

func (self *EmailMessage) SetTemplate(template CMessageTemplate) {
	self.tmplate = template
}

func (self *EmailMessage) Send(from string, to []string, params map[string]string) error {
	if self.initResult != nil {
		return self.initResult
	}
	content := self.tmplate.FillParams(params)
	return smtp.SendMail(fmt.Sprintf("%s:%d", self.host, self.port), self.auth, self.username, to, []byte(content))
}

func (self *EmailMessage) Initialize(params map[string]string) {
	clog.Debug(params, params["UserName"], params["Password"], params["Host"])
	if params["UserName"] == "" || params["Password"] == "" || params["Host"] == "" {
		self.initResult = errors.New("initial failed, UserName or Password or Host is required")
	}
	if params["Port"] != "" {
		p, _ := strconv.Atoi(params["Port"])
		if p == 0 {
			p = 25
		}
		self.port = p
	} else {
		self.port = 25
	}
	host := params["Host"]
	self.host = host
	self.username = params["UserName"]
	self.auth = smtp.PlainAuth("", params["UserName"], params["Password"], host)
}
