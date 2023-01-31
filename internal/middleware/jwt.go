package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccessDetails struct {
	ID   int64
	User string
}

func verifyToken(r *http.Request) (*jwt.Token, error) {

	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("private_key_sec"), nil
		})
		if err != nil {
			return nil, err
		}
		return token, nil
	} else {
		return nil, errors.New("token not found")
	}
}

func ExtractTokenMetadata(c *gin.Context) (*AccessDetails, error) {
	var (
		err           error
		accessExpires int64
	)
	token, err := verifyToken(c.Request)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessExpires, err = strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		if err != nil {
			goerrors.Log().WithError(err).Error("can't parse claims.exp")
			return nil, err
		}
		id, ok := claims["id"].(float64)
		if !ok {
			goerrors.Log().Error("can't parse claims.id")
			return nil, errors.New("can't parse claims.id")
		}

		user, ok := claims["user"].(string)

		if !ok {
			goerrors.Log().Error("can't parse claims.user")
			return nil, errors.New("can't parse claims.user")
		}

		accessDet := &AccessDetails{
			ID:   int64(id),
			User: user,
		}

		goerrors.Log().Info(c.ClientIP(), "jwt token params:", accessDet, " expires:", time.Unix(accessExpires, 0))

		return accessDet, nil

	}
	return nil, err
}
