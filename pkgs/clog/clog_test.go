package clog

import (
	"testing"
)

func TestClog_Println(t *testing.T) {
	l := NewClog()
	l.AddCustomFormatFunc("a", func(level LoggerLevel, skip int) string {
		return "hahaha"
	})
	l.SetFormat(l.GetFormat() + " $a")
	l.Error("error")
	l.Info("info")

	SetFormat(GetFormat() + " $a")
	Info("aaaa", len(support))
	Error("xxxx")
}

func TestClog_Println2(t *testing.T) {
	Info("aaaa")
	Error("xxxx")
}
