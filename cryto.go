package main

import (
	"crypto/md5"
	"fmt"
)

const TOKEN_KEY string = "TAPA"

func GenToken(user string, salt string) string {
	token := fmt.Sprintf("u=%s&salt=%s&key=%x", user, salt, md5.Sum([]byte(TOKEN_KEY)))
	token = fmt.Sprintf("%x", md5.Sum([]byte(token)))
	return token
}

func VerifyToken(user string, salt string, userToken string) bool {
	token := GenToken(user, salt)
	if token == userToken {
		return true
	}
	return false
}
