/*
 * Copyright (c) 2019.
 */

package cerror

import (
	"encoding/json"
	"fmt"
)

type CError struct {
	format     ErrorMessageFormat     `json:"-"`
	ErrCode    ErrorCode              `json:"code"`
	ErrMessage string                 `json:"message"`
	Params     map[string]interface{} `json:"-"`
}

func NewJsonError(code ErrorCode) *CError {
	return &CError{
		format:     ErrFmt_Json,
		ErrCode:    code,
		ErrMessage: code.ErrorMessage(),
		Params:     code.DefaultParams(),
	}
}

func NewPlainError(code ErrorCode) *CError {
	return &CError{
		format:     ErrFmt_Plain,
		ErrCode:    code,
		ErrMessage: code.ErrorMessage(),
		Params:     code.DefaultParams(),
	}
}

func NewJsonErrorWithParams(code ErrorCode, params map[string]interface{}) *CError {
	return &CError{
		format:     ErrFmt_Json,
		ErrCode:    code,
		ErrMessage: code.ErrorMessage(),
		Params:     params,
	}
}

func NewPlainErrorWithParams(code ErrorCode, params map[string]interface{}) *CError {
	return &CError{
		format:     ErrFmt_Plain,
		ErrCode:    code,
		ErrMessage: code.ErrorMessage(),
		Params:     params,
	}
}

func (ce CError) Error() string {
	if ce.ErrCode == nil {
		return ""
	}
	if len(ce.Params) > 0 {
		ce.ErrMessage = ce.ErrCode.FillParams(ce.ErrCode.ErrorMessage(), ce.Params)
	} else {
		ce.ErrMessage = ce.ErrCode.FillParams(ce.ErrCode.ErrorMessage(), ce.ErrCode.DefaultParams())
	}

	var res = fmt.Sprintf("[code]%s [message]%s", ce.ErrCode.ToString(), ce.ErrMessage)
	if ce.format == ErrFmt_Json {
		tmp, _ := json.Marshal(ce)
		res = string(tmp)
	}
	return res
}

func (ce CError) Code() string {
	if ce.ErrCode == nil {
		return ""
	}
	return ce.ErrCode.ToString()
}

func (ce CError) Message() string {
	if ce.ErrCode == nil {
		return ""
	}
	if len(ce.Params) > 0 {
		ce.ErrMessage = ce.ErrCode.FillParams(ce.ErrCode.ErrorMessage(), ce.Params)
	} else {
		ce.ErrMessage = ce.ErrCode.FillParams(ce.ErrCode.ErrorMessage(), ce.ErrCode.DefaultParams())
	}
	return ce.ErrMessage
}
