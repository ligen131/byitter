package controllers

import (
	"byoj/model"
	"byoj/utils/logs"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func checkUserName(userName string) error {
	if userName == "" {
		return errors.New("Empty user_name.")
	}

	result, _ := model.FindUserByName(userName)
	if result.UserName == userName {
		return errors.New("Username have been used.")
	}

	return nil
}

func checkUserEmail(email string) error {
	if email == "" {
		return errors.New("Empty email.")
	}

	result, _ := model.FindUserByEmail(email)
	if result.Email == email {
		return errors.New("Email have been used.")
	}

	return nil
}

func UserRegisterPOST(c echo.Context) error {
	logs.Debug("POST /user/register")

	user := model.User{}
	err := c.Bind(&user)
	if err != nil {
		logs.Warn("Failed to parse request data.", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			Code:    http.StatusBadRequest,
			Message: "Failed to parse request data.",
			Err:     err.Error(),
		})
	}
	logs.Debug("User struct:", zap.Any("user", user))

	err = checkUserName(user.UserName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     "",
		})
	}

	err = checkUserEmail(user.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Err:     "",
		})
	}

	err = model.UserRegister(user.UserName, user.Email, user.PasswordMD5, user.RealName, user.Bio)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorMessage{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user into database.",
			Err:     err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ErrorMessage{
		Code:    http.StatusOK,
		Message: "Register successfully.",
		Err:     "",
	})
}
