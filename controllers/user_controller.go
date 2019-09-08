package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/zouxinjiang/le/api"
	"github.com/zouxinjiang/le/core"
	"github.com/zouxinjiang/le/pkgs/cerror"
	"github.com/zouxinjiang/le/pkgs/constraint"
	"github.com/zouxinjiang/le/types"
)

type UserController struct {
	core.Controller
}

func (self UserController) GetMyInfo(c echo.Context) error {
	loginUser := c.Get("LoginUser").(types.UserSession)
	userApi := api.UserApi{}
	user, err := userApi.GetUserByUserName(loginUser.UserName)
	if err != nil {
		return err
	}
	return self.RespJson(c, user)
}

func (self UserController) UpdateUserInfo(c echo.Context) error {
	var params = struct {
		Name   string `json:"Name" form:"Name" query:"Name"`
		Email  string `json:"Email" form:"Email" query:"Email" constraint:"type:email"`
		Mobile string `json:"Mobile" form:"Mobile" query:"Mobile" constraint:"type:mobile"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if err := constraint.Valid(&params); err != nil {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  err.Error(),
			"reason": "value is not right",
		})
	}
	fields := map[string]interface{}{}
	if params.Name != "" {
		fields["Name"] = params.Name
	}
	if params.Mobile != "" {
		fields["Mobile"] = params.Mobile
	}
	if params.Email != "" {
		fields["Email"] = params.Email
	}
	userApi := api.UserApi{}
	loginUser := c.Get("LoginUser").(types.UserSession)
	userinfo, _ := userApi.GetUserByUserName(loginUser.UserName)

	err := userApi.UpdateUser(userinfo.Id, fields)
	if err != nil {
		return err
	}
	return self.RespJson(c, "")
}

func (self UserController) ChangeMyPassword(c echo.Context) error {
	var params = struct {
		OldPassword string `json:"OldPassword" form:"OldPassword" query:"OldPassword"`
		NewPassword string `json:"NewPassword" form:"NewPassword" query:"NewPassword"`
	}{}
	if err := c.Bind(&params); err != nil {
		return cerror.NewJsonError(core.ErrCode_InvalidParams)
	}
	if params.OldPassword == "" {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  "OldPassword",
			"reason": " value must specific",
		})
	}
	if params.NewPassword == "" {
		return cerror.NewJsonErrorWithParams(core.ErrCode_InvalidParam, map[string]interface{}{
			"field":  "NewPassword",
			"reason": " value must specific",
		})
	}
	loginUser := c.Get("LoginUser").(types.UserSession)
	userApi := api.UserApi{}
	uinfo, err := userApi.GetUserByUserName(loginUser.UserName)
	if err != nil {
		return err
	}
	err = userApi.ChangeUserPassword(uinfo.Id, params.OldPassword, params.NewPassword)
	if err != nil {
		return err
	}
	return self.RespJson(c, "")
}
