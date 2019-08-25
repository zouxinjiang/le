package clog

import (
	"io"
)

type Logger interface {
	SetShowLevel(level LogLevel)
	SetPrefix(prefix string)
	SetWriter(writer io.Writer)

	Println(level LogLevel, v ...interface{})
	Print(level LogLevel, v ...interface{})
	Printf(level LogLevel, fmt string, v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
	//设置日志格式
	SetFormat(format string)
	GetFormat() string
	//添加自定义日志格式函数
	AddCustomFormatFunc(name string, fn FormatFunc)
}
