package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"smart-learning/pkg/jwt"
)

func TestManager_GenerateAndParse(t *testing.T) {
	mgr := jwt.NewManager("test-secret-key-1234567890", time.Hour, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair("user-uuid-1", "student")
	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)
	assert.Greater(t, pair.ExpiresIn, int64(0))

	claims, err := mgr.ParseToken(pair.AccessToken)
	require.NoError(t, err)
	assert.Equal(t, "user-uuid-1", claims.UserID)
	assert.Equal(t, "student", claims.Role)
	assert.Equal(t, jwt.TypeAccess, claims.Type)
}

func TestManager_RefreshToken(t *testing.T) {
	mgr := jwt.NewManager("test-secret-key-1234567890", time.Hour, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair("user-uuid-1", "student")
	require.NoError(t, err)

	newPair, err := mgr.RefreshToken(pair.RefreshToken)
	require.NoError(t, err)
	assert.NotEqual(t, pair.AccessToken, newPair.AccessToken)
	assert.NotEqual(t, pair.RefreshToken, newPair.RefreshToken)
}

func TestManager_RefreshWithAccessToken_Fails(t *testing.T) {
	mgr := jwt.NewManager("test-secret-key-1234567890", time.Hour, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair("user-uuid-1", "student")
	require.NoError(t, err)

	_, err = mgr.RefreshToken(pair.AccessToken)
	assert.ErrorIs(t, err, jwt.ErrTokenWrongType)
}

func TestManager_ParseExpiredToken(t *testing.T) {
	mgr := jwt.NewManager("test-secret-key-1234567890", time.Millisecond, 24*time.Hour)
	pair, err := mgr.GenerateTokenPair("user-uuid-1", "student")
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	_, err = mgr.ParseToken(pair.AccessToken)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestManager_ParseInvalidToken(t *testing.T) {
	mgr := jwt.NewManager("test-secret-key-1234567890", time.Hour, 24*time.Hour)
	_, err := mgr.ParseToken("invalid-token")
	assert.Error(t, err)
}

func TestManager_ParseWithWrongSecret(t *testing.T) {
	mgr1 := jwt.NewManager("secret-one", time.Hour, 24*time.Hour)
	mgr2 := jwt.NewManager("secret-two", time.Hour, 24*time.Hour)
	pair, err := mgr1.GenerateTokenPair("user-uuid-1", "student")
	require.NoError(t, err)

	_, err = mgr2.ParseToken(pair.AccessToken)
	assert.Error(t, err)
}