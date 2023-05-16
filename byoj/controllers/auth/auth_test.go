package auth_test

import (
	"byoj/controllers/auth"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestAuth(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
		ID:       1,
		UserName: "ligen131",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	})
	t.Log("token = ", token)
	tokenString, err := token.SignedString([]byte("123"))
	t.Log("signed = ", tokenString, err)
	// token2, _ := jwt.ParseWithClaims(tokenString, auth.Claims{}, )
	cl := new(auth.Claims)
	token2, _ := jwt.ParseWithClaims(tokenString, cl, func(token *jwt.Token) (interface{}, error) {
		return []byte("123"), nil
	})
	t.Log(cl.ID, cl.UserName, token2)
	t.FailNow()
}
