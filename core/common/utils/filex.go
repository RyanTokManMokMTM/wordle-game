package utils

import (
	"os"
	"strings"
)

func LoadFile(file string) ([]byte, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DataBytesToList(dataBytes []byte, separatedBy string) []string {
	return strings.Split(string(dataBytes), separatedBy)
}
