package core

import (
	"fmt"
	"strings"
)

type CustomErrorCode struct {
	Code    string
	Message string
	Params  map[string]interface{}
}

func (self CustomErrorCode) ToString() string {
	return self.Code
}

func (self CustomErrorCode) ErrorMessage() string {
	return self.Message
}

func (self CustomErrorCode) DefaultParams() map[string]interface{} {
	return self.Params
}

func (self CustomErrorCode) FillParams(tmpl string, params map[string]interface{}) string {
	for k, v := range params {
		tmpl = strings.Replace(tmpl, fmt.Sprintf("${%s}", k), fmt.Sprint(v), -1)
	}
	return tmpl
}

var (
	ErrCode_InvalidParam = CustomErrorCode{
		Code:    "InvalidParam",
		Message: "param field:${field} is invalid as ${reason}",
		Params: map[string]interface{}{
			"field":  "-",
			"reason": "-",
		},
	}
	ErrCode_InvalidParams = CustomErrorCode{
		Code:    "InvalidParams",
		Message: "the api params is invalid",
		Params:  nil,
	}
	ErrCode_DbConnectFailed = CustomErrorCode{
		Code:    "DbConnectFailed",
		Message: "database connect failed",
		Params:  nil,
	}
	ErrCode_RecordNotExist = CustomErrorCode{
		Code:    "RecordNotExist",
		Message: "record ${field} not exist",
		Params: map[string]interface{}{
			"field": "",
		},
	}
	ErrCode_RecordExisted = CustomErrorCode{
		Code:    "RecordExisted",
		Message: "record ${field} has existed",
		Params: map[string]interface{}{
			"field": "",
		},
	}
	ErrCode_Unknown = CustomErrorCode{
		Code:    "Unknown",
		Message: "unknown error",
		Params:  nil,
	}
)
