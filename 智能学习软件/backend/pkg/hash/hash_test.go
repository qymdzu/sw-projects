package hash_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/pkg/hash"
)

func TestHashAndVerify_Success(t *testing.T) {
	hashed, err := hash.Hash("Abc12345!")
	require.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NoError(t, hash.Verify(hashed, "Abc12345!"))
}

func TestVerify_Mismatch(t *testing.T) {
	hashed, err := hash.Hash("Abc12345!")
	require.NoError(t, err)
	assert.ErrorIs(t, hash.Verify(hashed, "WrongPassword!"), hash.ErrPasswordMismatch)
}

func TestHash_Empty(t *testing.T) {
	_, err := hash.Hash("")
	assert.Error(t, err)
}

func TestVerify_EmptyPassword(t *testing.T) {
	hashed, err := hash.Hash("Abc12345!")
	require.NoError(t, err)
	// 用空密码校验应失败
	assert.Error(t, hash.Verify(hashed, ""))
}