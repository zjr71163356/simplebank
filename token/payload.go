package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		Id:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil

}

// Valid checks if the token payload is valid
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return jwt.ErrTokenExpired
	}

	return nil
}

// GetExpirationTime returns the expiration time claim.
func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.ExpiredAt), nil
}

// GetIssuedAt returns the issued at time claim.
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(payload.IssuedAt), nil
}

// GetNotBefore returns the not before time claim.
// Payload does not implement this claim, so it returns nil.
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil // Or return an error if this claim is mandatory
}

// GetIssuer returns the issuer claim.
// Payload does not implement this claim, so it returns an empty string.
func (payload *Payload) GetIssuer() (string, error) {
	return "", nil // Or return an error if this claim is mandatory
}

// GetSubject returns the subject claim.
// Payload does not implement this claim, so it returns an empty string.
func (payload *Payload) GetSubject() (string, error) {
	return "", nil // Or return an error if this claim is mandatory
}

// GetAudience returns the audience claim.
// Payload does not implement this claim, so it returns nil.
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil // Or return an error if this claim is mandatory
}
