package token

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/vldcreation/movie-fest/pkg/util"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker("a_32_byte_secret_key_for_testing")
	require.NoError(t, err)

	username := "vld"
	duration := 1 * time.Minute

	token, err := maker.CreateToken(username, duration, map[string]any{"foo": "bar"})
	require.NoError(t, err)

	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)

	require.NotEmpty(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), 5*time.Second)

	bar, ok := payload.GetCustomClaims("foo")
	require.True(t, ok)
	require.Equal(t, "bar", bar)
}

func TestExpiredTokenPaseto(t *testing.T) {
	maker, err := NewPasetoMaker("a_32_byte_secret_key_for_testing")
	require.NoError(t, err)

	username := "vld"
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration, nil)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)

	err = payload.Valid()
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidTokenPaseto(t *testing.T) {
	payload, err := NewPayload("vld", 1*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewPasetoMaker(util.RandString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.True(t, strings.HasSuffix(err.Error(), ErrInvalidToken.Error()))
	require.Nil(t, payload)
}
