package passhandle_test

import (
	"testing"
	passhandle "users/internal/passhandler"
	"users/oas"

	"github.com/stretchr/testify/assert"
)

func TestPassEncrypting(t *testing.T) {
	login := oas.LoginString("abacaba")
	pass := oas.PasswordString("abacaba")

	assert.NotEqual(t, passhandle.HashPass(login, pass), pass)
}

func TestPassLoginDependent(t *testing.T) {
	login1 := oas.LoginString("user1")
	login2 := oas.LoginString("user2")
	pass := oas.PasswordString("pass")

	h1 := passhandle.HashPass(login1, pass)
	h2 := passhandle.HashPass(login2, pass)
	assert.NotEqual(t, h1, h2)
}
func TestPassDeterministic(t *testing.T) {
	l1 := oas.LoginString("user")
	p1 := oas.PasswordString("pass")
	h1 := passhandle.HashPass(l1, p1)

	l2 := oas.LoginString("user")
	p2 := oas.PasswordString("pass")
	h2 := passhandle.HashPass(l2, p2)

	assert.Equal(t, h1, h2)
}
