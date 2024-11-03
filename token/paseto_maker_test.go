package token

import (
	"testing"
	"time"

	"github.com/ptdatta/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker,err := NewPasetoMaker(util.RandomString(32))
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

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

