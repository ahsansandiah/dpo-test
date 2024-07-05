package jwtAuth

import (
	"fmt"
	"strings"
	"time"

	"github.com/ahsansandiah/dpo-test/packages/config"
	"github.com/golang-jwt/jwt"
)

type Jwt interface {
	GenerateToken(data *JwtData) (string, *time.Time, error)
	ExtractJwtToken(token string) (*JwtData, error)
	VerifyAccessToken(token string, secretKey string) (*JwtData, error)
}

type Options struct {
	secretKey           string
	accessTokenDuration int
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.secretKey = cfg.JwtSecretKey
	opt.accessTokenDuration = cfg.JwtAccessTokenDuration

	return opt
}

func (o *Options) GenerateToken(data *JwtData) (string, *time.Time, error) {
	jwtPayload := &JwtPayload{
		Reference: data.Reference,
		UserID:    data.UserID,
	}

	expiredTime := time.Now().Local().Add(time.Second * time.Duration(o.accessTokenDuration))
	jwtPayload.StandardClaims.ExpiresAt = expiredTime.Unix()
	jwtPayload.StandardClaims.NotBefore = jwt.TimeFunc().Local().Unix()
	acToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtPayload)

	accessToken, err := acToken.SignedString([]byte(o.secretKey))
	if err != nil {
		return "", nil, err
	}

	return accessToken, &expiredTime, nil
}

func (o *Options) ExtractJwtToken(token string) (*JwtData, error) {
	parsedToken, err := jwt.Parse(token, nil)
	if strings.Contains(err.Error(), "invalid number of segments") {
		return nil, err
	}

	claims, _ := parsedToken.Claims.(jwt.MapClaims)

	userId := int64(claims["ui"].(float64))
	jwtData := &JwtData{
		Reference: claims["reference"].(string),
		UserID:    userId,
	}

	return jwtData, nil
}

func (o *Options) VerifyAccessToken(token string, secretKey string) (*JwtData, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, err
	}

	jwtData := &JwtData{
		Reference: claims["reference"].(string),
	}

	return jwtData, nil
}
