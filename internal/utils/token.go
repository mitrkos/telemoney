package utils

import (
	"os"
)

func GetTokenFromFile(fileName string) (string, error) {
	token, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
