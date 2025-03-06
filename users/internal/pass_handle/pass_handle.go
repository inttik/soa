package passhandle

import (
	"crypto/md5"
	"io"
	"users/oas"
)

func HashPass(login oas.LoginString, pass oas.PasswordString) string {
	toEncrypt := string(pass) + "salt" + string(login)
	hasher := md5.New()
	io.WriteString(hasher, toEncrypt)
	hash := hasher.Sum(nil)
	return string(hash)
}
