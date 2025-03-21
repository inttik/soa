package passhandle_test

import (
	"testing"
	passhandle "users/internal/passhandler"
	"users/oas"

	"github.com/stretchr/testify/assert"
)

func TestPassEncrypting(t *testing.T) {
	pass := oas.PasswordString("abacaba")
	hash, err := passhandle.HashPassword(pass)
	assert.NoErrorf(t, err, "pass should be hashed")
	assert.NotEqual(t, hash, pass)
}

func TestPassValiud(t *testing.T) {
	pass := oas.PasswordString("abacaba")
	hash, err := passhandle.HashPassword(pass)
	assert.NoErrorf(t, err, "pass should be hashed")
	assert.NotEqual(t, hash, pass)

	same := passhandle.ComparePassword(string(pass), hash)
	assert.Equalf(t, true, same, "pass should be accepted")
}
