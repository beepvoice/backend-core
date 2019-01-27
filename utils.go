package main

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ttacon/libphonenumber"
)

func RandomHex() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("unable to generate 16 bytes of randomness")
	}
	return hex.EncodeToString(b)
}

func ParsePhone(phone string) (string, error) {
	num, err := libphonenumber.Parse(phone, "")
	if err != nil {
		return "", err
	}
	return libphonenumber.Format(num, libphonenumber.INTERNATIONAL), nil
}
