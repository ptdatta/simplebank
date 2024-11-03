package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ptdatta/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker,err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)

	username := util.RandomOwner()
	duration := time.Minute

	token,err := maker.CreateToken(username,duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)

	payload,err := maker.VerifyToken(token)
	require.NoError(t,err)
	require.NotEmpty(t,payload)

	require.NotZero(t,payload.ID)
	expTime, err := payload.GetExpirationTime()
	require.NoError(t, err)
	require.Equal(t, payload.ExpiredAt.Unix(), expTime.Unix())

	issuedAt, err := payload.GetIssuedAt()
	require.NoError(t, err)
	require.Equal(t, payload.IssuedAt.Unix(), issuedAt.Unix())

	subject, err := payload.GetSubject()
	require.NoError(t, err)
	require.Equal(t, username, subject)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}