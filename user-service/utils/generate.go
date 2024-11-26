package utils

import (
	"strconv"

	"golang.org/x/exp/rand"
)

func GenerateCode() string {
	code := rand.Intn(999999) + 100000
	return strconv.Itoa(code)
}
