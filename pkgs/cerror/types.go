/*
 * Copyright (c) 2019.
 */

package cerror

type ErrorMessageFormat string

const (
	ErrFmt_Json  ErrorMessageFormat = "json"
	ErrFmt_Plain ErrorMessageFormat = "plain"
)

type ErrorCode interface {
	ToString() string
	ErrorMessage() string
	FillParams(tmpl string, params map[string]interface{}) string
	DefaultParams() map[string]interface{}
}
