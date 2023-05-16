package controllers

import (
	"byoj/controllers/auth"
	"byoj/model"
	"byoj/utils/logs"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Failed to parse request data.",
				Err:     err.Error(),
			},
		})
	}
	logs.Debug("User struct:", zap.Any("user", user))

	err = checkUserName(user.UserName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: err.Error(),
				Err:     "",
			},
		})
	}

	err = checkUserEmail(user.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: err.Error(),
				Err:     "",
			},
		})
	}

	err = model.UserRegister(user.UserName, user.Email, user.PasswordMD5, user.RealName, user.Bio)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Failed to create user into database.",
				Err:     err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data: StatusMessage{
			Status: "Register successfully.",
		},
	})
}

type UserLoginResponse struct {
	ID                   uint32 `json:"user_id"`
	UserName             string `json:"user_name"`
	AccessToken          string `json:"access_token"`
	AccessTokenExpireAt  int64  `json:"access_token_expiration_time"`
	RefreshToken         string `json:"refresh_token"`
	RefreshTokenExpireAt int64  `json:"refresh_token_expiration_time"`
}

func UserLoginPOST(c echo.Context) error {
	logs.Debug("POST /user/login")

	userRequest := model.User{}
	err := c.Bind(&userRequest)
	if err != nil {
		logs.Warn("Failed to parse request data.", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Failed to parse request data.",
				Err:     err.Error(),
			},
		})
	}
	logs.Debug("User struct:", zap.Any("user", userRequest))

	user := model.User{}
	if userRequest.ID != 0 {
		user, err = model.FindUserByID(userRequest.ID)
	} else if userRequest.Email != "" {
		user, err = model.FindUserByEmail(userRequest.Email)
	} else if userRequest.UserName != "" {
		user, err = model.FindUserByName(userRequest.UserName)
	} else {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "User ID, email or username is required.",
				Err:     "",
			},
		})
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, ResponseStruct{
				Code:    http.StatusBadRequest,
				Message: "Bad Request",
				Data: ErrorMessage{
					Message: "User not found.",
					Err:     err.Error(),
				},
			})
		}
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Find user failed.",
				Err:     err.Error(),
			},
		})
	}

	if user.Deleted {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "This user has been deleted.",
				Err:     "",
			},
		})
	}

	if !user.Verified {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "This user has not been verified.",
				Err:     "",
			},
		})
	}

	if user.PasswordMD5 != userRequest.PasswordMD5 {
		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
			Data: ErrorMessage{
				Message: "Wrong password.",
				Err:     "",
			},
		})
	}

	accessTokenString, accessTokenExpireAt, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Generate access token failed.",
				Err:     err.Error(),
			},
		})
	}

	refreshTokenString, refreshTokenExpireAt, err := auth.GenerateRefreshToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStruct{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Data: ErrorMessage{
				Message: "Generate refresh token failed.",
				Err:     err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data: UserLoginResponse{
			ID:                   user.ID,
			UserName:             user.UserName,
			AccessToken:          accessTokenString,
			AccessTokenExpireAt:  accessTokenExpireAt.Unix(),
			RefreshToken:         refreshTokenString,
			RefreshTokenExpireAt: refreshTokenExpireAt.Unix(),
		},
	})
}

func UserIsAuthGET(c echo.Context) error {
	logs.Debug("GET /user/isauth")
	return c.JSON(http.StatusOK, ResponseStruct{
		Code:    http.StatusOK,
		Message: "OK",
		Data: StatusMessage{
			Status: "OK",
		},
	})
}
