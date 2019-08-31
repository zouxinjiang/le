package core

import (
	"github.com/labstack/echo"
	"github.com/zouxinjiang/le/types"
	"net/http"
)

type Controller struct {
}

func (c Controller) Init() error {
	return nil
}

func (c Controller) Install() error {
	return nil
}

func (c Controller) Start(params map[string]string) error {
	return nil
}

func (c Controller) Stop(params map[string]string) error {
	return nil
}

//============================= 所有controller 公共方法=======================================
func (self Controller) RespJson(c echo.Context, data interface{}) error {
	res := types.ApiResponseStructure{
		Code:    "Success",
		Message: "",
		Data:    data,
	}
	return c.JSON(http.StatusOK, res)
}

func (self Controller) RespJsonWithCount(c echo.Context, data interface{}, count int) error {
	res := types.ApiResponseStructure{
		Code:    "Success",
		Message: "",
		Count:   &count,
		Data:    data,
	}
	return c.JSON(http.StatusOK, res)
}
