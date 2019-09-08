package cmessage

import (
	"errors"
)

type UnknownMessage string

func init() {
	register("unknown", newUnknownMessage)
}

func newUnknownMessage() CMessage {
	return new(UnknownMessage)
}

func (u UnknownMessage) Initialize(params map[string]string) {
}

func (u UnknownMessage) SetTemplate(template CMessageTemplate) {
}

func (u UnknownMessage) Send(from string, to []string, params map[string]string) error {
	return errors.New("not support cmessage")
}
