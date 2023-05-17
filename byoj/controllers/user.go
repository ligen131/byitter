package controllers

import (
	"byoj/controllers/auth"
	"byoj/model"
	"byoj/utils/logs"
	"errors"

	"github.com/labstack/echo"
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

func FindUser(c echo.Context, request model.User) (user model.User, err error, isInternalServerError bool) {
	user = request
	if request.ID != 0 {
		user, err = model.FindUserByID(request.ID)
	} else if request.Email != "" {
		user, err = model.FindUserByEmail(request.Email)
	} else if request.UserName != "" {
		user, err = model.FindUserByName(request.UserName)
	} else {
		return user, errors.New("Author ID, email or username is required."), false
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, err, false
		}
		return user, ResponseInternalServerError(c, "Find user failed", err), true
	}
	return user, nil, false
}

func UserRegisterPOST(c echo.Context) error {
	logs.Debug("POST /user/register")

	user := model.User{}
	_ok, err := Bind(c, &user)
	if !_ok {
		return err
	}

	err = checkUserName(user.UserName)
	if err != nil {
		return ResponseBadRequest(c, err.Error(), nil)
	}

	err = checkUserEmail(user.Email)
	if err != nil {
		return ResponseBadRequest(c, err.Error(), nil)
	}

	err = model.UserRegister(user.UserName, user.Email, user.PasswordMD5, user.RealName, user.Bio)
	if err != nil {
		return ResponseInternalServerError(c, "Failed to create user into database.", err)
	}

	return ResponseOK(c, StatusMessage{
		Status: "Register successfully.",
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
	_ok, err := Bind(c, &userRequest)
	if !_ok {
		return err
	}

	user, err, e500 := FindUser(c, model.User{
		ID:       userRequest.ID,
		UserName: userRequest.UserName,
		Email:    userRequest.Email,
	})
	if e500 {
		return err
	}
	if err != nil {
		return ResponseBadRequest(c, "Find user failed.", err)
	}

	if user.Deleted {
		return ResponseBadRequest(c, "This user has been deleted.", nil)
	}

	if !user.Verified {
		return ResponseBadRequest(c, "This user has not been verified.", nil)
	}

	if user.PasswordMD5 != userRequest.PasswordMD5 {
		return ResponseBadRequest(c, "Wrong password.", nil)
	}

	accessTokenString, accessTokenExpireAt, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return ResponseInternalServerError(c, "Generate access token failed.", err)
	}

	refreshTokenString, refreshTokenExpireAt, err := auth.GenerateRefreshToken(&user)
	if err != nil {
		return ResponseInternalServerError(c, "Generate refresh token failed.", err)
	}

	return ResponseOK(c, UserLoginResponse{
		ID:                   user.ID,
		UserName:             user.UserName,
		AccessToken:          accessTokenString,
		AccessTokenExpireAt:  accessTokenExpireAt.Unix(),
		RefreshToken:         refreshTokenString,
		RefreshTokenExpireAt: refreshTokenExpireAt.Unix(),
	})
}

func UserIsAuthGET(c echo.Context) error {
	logs.Debug("GET /user/isauth")

	return ResponseOK(c, StatusMessage{
		Status: "OK",
	})
}

type UserGETResponse struct {
	ID       uint32 `json:"user_id"   `
	UserName string `json:"user_name" `
	Email    string `json:"email"     `
	RealName string `json:"real_name" `
	Bio      string `json:"bio"       `
	Verified bool   `json:"verified"  `
	Deleted  bool   `json:"deleted"   `
}

func UserGET(c echo.Context) error {
	logs.Debug("GET /user")

	userRequest := model.User{}
	_ok, err := Bind(c, &userRequest)
	if !_ok {
		return err
	}

	user, err, e500 := FindUser(c, model.User{
		ID:       userRequest.ID,
		UserName: userRequest.UserName,
		Email:    userRequest.Email,
	})
	if e500 {
		return err
	}
	if err != nil {
		return ResponseBadRequest(c, "Find user failed.", err)
	}

	return ResponseOK(c, UserGETResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		RealName: user.RealName,
		Bio:      user.Bio,
		Verified: user.Verified,
		Deleted:  user.Deleted,
	})
}
