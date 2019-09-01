package twofactor

import (
	"errors"
)

type UnknownFactor string

func init() {
	register("unknown", newUnknownFactor)
}

func newUnknownFactor() TwoFactor {
	return new(UnknownFactor)
}

func (UnknownFactor) Do(params map[string]string) (string, error) {
	return "unknown", errors.New("unknown two factor")
}
