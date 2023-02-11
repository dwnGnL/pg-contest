package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
)

type JwtToken[claim MyClaim] struct {
	secretKey string
}

func New[claim MyClaim](key string) JwtToken[claim] {
	return JwtToken[claim]{secretKey: key}
}

type MyClaim interface {
	Valid() error
}

func (j *JwtToken[claim]) verifyToken(bearerToken string) (claim, error) {
	var nilClaim claim
	goerrors.Log().Println("start verifyToken")

	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		goerrors.Log().Println("start jwt.Parse")
		token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(j.secretKey), nil
		})
		if err != nil {
			return nilClaim, err
		}
		goerrors.Log().Println("start token.Claims")

		if v, ok := token.Claims.(jwt.MapClaims); ok {
			jsonbody, err := json.Marshal(v)
			if err != nil {
				return nilClaim, err
			}

			if err := json.Unmarshal(jsonbody, &nilClaim); err != nil {
				return nilClaim, err
			}
			return nilClaim, nil
		} else {

			goerrors.Log().Info(token.Claims)
			goerrors.Log().Printf("%#v", nilClaim)
			return nilClaim, errors.New("claims not valid")
		}
		//return token, nil
	} else {
		goerrors.Log().Println("token not found")
		return nilClaim, errors.New("token not found")
	}
}

func (j *JwtToken[claim]) ExtractTokenMetadata(bearerToken string) (claim, error) {
	var nilClaim claim
	goerrors.Log().Println("start ExtractTokenMetadata")
	token, err := j.verifyToken(bearerToken)
	if err != nil {
		return nilClaim, err
	}
	goerrors.Log().Println("start Valid")

	if token.Valid() != nil {
		goerrors.Log().Info("token not valid")
		return nilClaim, errors.New("token not valid")
	}
	goerrors.Log().Info("token: ", token)

	return token, err
}
