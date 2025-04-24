package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	minSecretKeySize = 12
)

// JWT分为header、payload、signature三部分
// 其中signature是通过非对称或对称数字签名算法得到的，但在jwt包中可以直接创建token
// 而不是单独创建signature后将JWT的三部分拼接在一起

type JWTMaker struct {
	secretKey string //用来创建token的密钥
}

func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

// token的创建，使用NewWithClaims需要指定使用的签名算法、再输入payload，这样就能创建Token
// 使用创建的Token类型的变量的SignedString，以密钥作为参数，SignedString返回创建完整的token
// 其中duration指的是创建和过期之间间隔的时间
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))

}

func (jwtmaker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	jwtKeyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(jwtmaker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, jwtKeyFunc)

	if err != nil {
		if errors.Is(jwt.ErrTokenExpired, err) {
			return nil, jwt.ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
