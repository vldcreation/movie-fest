package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`

	// payload.CustomClaims
	Meta CustomClaims `json:"meta"`
}

type CustomClaims struct {
	Data map[string]any `json:"metas"`
}

type Opt func(*Payload)

func WithCustomClaims(meta CustomClaims) Opt {
	return func(p *Payload) {
		if meta.Data == nil {
			meta.Data = make(map[string]any)
		}

		p.Meta = meta
	}
}

func WithCustomClaimsData(key string, value any) Opt {
	return func(p *Payload) {
		p.Meta.Data[key] = value
	}
}

func (p *Payload) GetCustomClaims(key string) (any, bool) {
	value, ok := p.Meta.Data[key]
	return value, ok
}

func NewPayloadWithOpts(username string, duration time.Duration, opts ...Opt) (*Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(payload)
	}
	return payload, nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{p.Username}, nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Now()), nil
}

func (p *Payload) GetIssuer() (string, error) {
	return "github.com/vldcreation/movie-fest", nil
}

func (p *Payload) GetSubject() (string, error) {
	return "auth", nil
}
