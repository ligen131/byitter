package middleware

import (
	"byoj/controllers"
	"byoj/controllers/auth"
	"byoj/model"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func TokenVerificationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := auth.GetClaimsFromHeader(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, controllers.ResponseStruct{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data: controllers.ErrorMessage{
					Message: "Invalid bearer token in header.",
					Err:     err.Error(),
				},
			})
		}
		if claims.Valid() != nil {
			return c.JSON(http.StatusUnauthorized, controllers.ResponseStruct{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data: controllers.ErrorMessage{
					Message: "Invalid jwt token.",
					Err:     claims.Valid().Error(),
				},
			})
		}

		if claims.ExpiresAt < time.Now().Unix() {
			return c.JSON(http.StatusUnauthorized, controllers.ResponseStruct{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data: controllers.ErrorMessage{
					Message: "Token expired.",
					Err:     "",
				},
			})
		}

		user, err := model.FindUserByID(claims.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.JSON(http.StatusUnauthorized, controllers.ResponseStruct{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized",
					Data: controllers.ErrorMessage{
						Message: "User in token not found.",
						Err:     err.Error(),
					},
				})
			}
			return c.JSON(http.StatusInternalServerError, controllers.ResponseStruct{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
				Data: controllers.ErrorMessage{
					Message: "Find user by ID failed.",
					Err:     err.Error(),
				},
			})
		}
		if user.UserName != claims.UserName {
			return c.JSON(http.StatusUnauthorized, controllers.ResponseStruct{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data: controllers.ErrorMessage{
					Message: "UserID does not match username.",
					Err:     err.Error(),
				},
			})
		}

		return next(c)
	}
}
