package utils

import (
	"math/rand"
	"strings"
)

type RandomBase62StringGenerator int

func (r *RandomBase62StringGenerator) GenerateURL(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var builder strings.Builder

	for i := 0; i < length; i++ {
		index := rand.Int31n(62)
		_, err := builder.WriteString(chars[index : index+1])
		if err != nil {
			return "", err
		}
	}
	return builder.String(), nil
}
