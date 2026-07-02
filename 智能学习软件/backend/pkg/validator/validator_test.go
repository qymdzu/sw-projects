package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"smart-learning/pkg/validator"
)

func TestRequireString(t *testing.T) {
	assert.NoError(t, validator.RequireString("name", "张三"))
	assert.Error(t, validator.RequireString("name", ""))
	assert.Error(t, validator.RequireString("name", "   "))
}

func TestPhone(t *testing.T) {
	assert.NoError(t, validator.Phone("13800138000"))
	assert.NoError(t, validator.Phone("15912345678"))
	assert.Error(t, validator.Phone("12345"))
	assert.Error(t, validator.Phone("23800138000"))
	assert.Error(t, validator.Phone(""))
}

func TestEmail(t *testing.T) {
	assert.NoError(t, validator.Email("a@b.com"))
	assert.NoError(t, validator.Email("user.name+tag@example.co.uk"))
	assert.NoError(t, validator.Email("")) // 允许空
	assert.Error(t, validator.Email("not-an-email"))
}

func TestPassword(t *testing.T) {
	assert.NoError(t, validator.Password("Abc12345"))
	assert.NoError(t, validator.Password("strong1Password"))
	assert.Error(t, validator.Password("short1"))
	assert.Error(t, validator.Password("NoDigits!"))
	assert.Error(t, validator.Password("nodigitshere"))
	assert.Error(t, validator.Password("12345678"))
}

func TestRole(t *testing.T) {
	assert.NoError(t, validator.Role("student"))
	assert.NoError(t, validator.Role("teacher"))
	assert.NoError(t, validator.Role("admin"))
	assert.NoError(t, validator.Role("parent"))
	assert.Error(t, validator.Role("guest"))
	assert.Error(t, validator.Role(""))
}