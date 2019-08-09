package main

import (
	"crypto/sha256"
)

func simpleHash(msg ...[]byte) []byte {
	hasher := sha256.New()
	for i := range msg {
		hasher.Write(msg[i])
	}
	return hasher.Sum(nil)
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func reverseStr(s string) string {
	newStr := ""
	for i := len(s) - 1; i >= 0; i-- {
		newStr += string(s[i])
	}
	return newStr
}
