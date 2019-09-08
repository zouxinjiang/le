package main

import (
	"fmt"
	"github.com/zouxinjiang/le/pkgs/clog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/types"
)

func AppErrorHandleFunc(err error, c echo.Context) {
	clog.Println(clog.Lvl_Debug, err, reflect.TypeOf(err))
	var (
		code  = http.StatusInternalServerError
		xcode interface{}
		msg   interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		xcode = he.Code
		msg = he.Message
	} else if ce, ok := err.(cerror.CError); ok {
		xcode = ce.Code()
		msg = ce.Message()
	} else if ce, ok := err.(*cerror.CError); ok {
		xcode = ce.Code()
		msg = ce.Message()
	} else {
		xcode = code
		msg = http.StatusText(code)
	}

	res := types.ApiResponseStructure{
		Code:    fmt.Sprint(xcode),
		Message: fmt.Sprint(msg),
	}
	// Send response
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			c.NoContent(200)
		} else {
			c.JSON(200, res)
		}
	}
}
