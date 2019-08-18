/*
 * Copyright (c) 2019.
 */

package clog

import (
	"io"
)

type LoggerLevel int

const (
	Lvl_Info LoggerLevel = 1 << iota
	Lvl_Warning
	Lvl_Error
	Lvl_Debug
)

var str2lvl = map[LoggerLevel]string{
	Lvl_Debug:   "DEBUG",
	Lvl_Error:   "ERROR",
	Lvl_Info:    "INFO",
	Lvl_Warning: "WARNING",
}

func (ll LoggerLevel) String() string {
	return str2lvl[ll]
}

type FormatFunc func(level LoggerLevel, skip int) string

type Logger interface {
	//设置显示日志的最低等级
	SetShowLevel(level LoggerLevel)
	Println(level LoggerLevel, v ...interface{})
	Print(level LoggerLevel, v ...interface{})
	Printf(level LoggerLevel, fmt string, v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
	//设置输出位置
	SetWriter(wd io.Writer)
	SetPrefix(prefix string)
	//增加调用层数
	AddCallerLevel()
	//设置日志格式
	SetFormat(format string)
	GetFormat() string
	//添加自定义日志格式函数
	AddCustomFormatFunc(name string, fn FormatFunc)
}
