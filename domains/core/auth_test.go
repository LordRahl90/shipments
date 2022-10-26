package core

import (
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	SigningSecret = "hello world"
	userID, email := uuid.NewString(), gofakeit.Email()
	td := &TokenData{
		UserID: userID,
		Email:  email,
	}

	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)
	assert.Len(t, strings.Split(token, "."), 3) // token should always divided into 3 parts
}

func TestDecodeToken(t *testing.T) {
	SigningSecret = "hello world"
	userID, email := uuid.NewString(), gofakeit.Email()
	td := &TokenData{
		UserID: userID,
		Email:  email,
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	tk, err := Decode(token)
	require.NoError(t, err)
	require.NotNil(t, tk)
	assert.Equal(t, td.UserID, tk.UserID)
	assert.Equal(t, td.Email, tk.Email)
}
