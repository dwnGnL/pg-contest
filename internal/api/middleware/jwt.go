package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccessDetails struct {
	AccessUuid string
	MerchId    int64
	ExtTranID  string
	Phone      string
	Amount     float64
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

			return []byte(settings.AppSettings.TokenParams.AccessSecret), nil
		})
		if err != nil {
			return nil, err
		}
		return token, nil
	} else {
		return nil, errors.New("token not found")
	}
}

func extractTokenMetadata(c *gin.Context) (*AccessDetails, error) {
	token, err := verifyToken(c.Request)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		accessExpires, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
		merchId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["merch_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		extTranID, ok := claims["tran_id"].(string)
		phone, ok := claims["phone"].(string)
		amount, err := strconv.ParseFloat(fmt.Sprintf("%f", claims["amount"]), 64)

		accessDet := &models.AccessDetails{
			AccessUuid: accessUuid,
			MerchId:    merchId,
			ExtTranID:  extTranID,
			Phone:      phone,
			Amount:     amount,
		}

		logger.Info.Println(c.ClientIP(), "jwt token params:", settings.MakePsevdoJson(accessDet), " expires:", time.Unix(accessExpires, 0))

		return accessDet, nil

	}
	return nil, err
}

func getUserFromContext(c *gin.Context) UserRole {
	claims := ExtractClaims(c)
	return UserRole{
		UserID: int64(claims["userID"].(float64)),
		//TerritoryID: int64(claims["territoryID"].(float64)),
		Role:     claims["role"].(string),
		UserName: claims["userName"].(string),
	}
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(jwt.MapClaims)
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}
