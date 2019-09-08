package core

import (
	"strings"
)

func IsDbErrorRecordNotFount(err error) bool {
	if err == nil {
		return false
	}
	if idx := strings.Index(err.Error(), "record not found"); idx >= 0 {
		return true
	}
	return false
}

func IsDbErrorUnique(err error) bool {
	if err == nil {
		return false
	}
	if idx := strings.Index(err.Error(), ""); idx >= 0 {
		return true
	}
	return false
}

func IsDbErrorForeignKey(err error) bool {
	if err == nil {
		return false
	}
	if idx := strings.Index(err.Error(), ""); idx >= 0 {
		return true
	}
	return false
}
