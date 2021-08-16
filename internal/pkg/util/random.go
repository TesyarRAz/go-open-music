package util

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandomText(n int) string {
	result := ""

	for i := 0; i < n; i++ {
		result += string(letters[rand.Intn(len(letters))])
	}

	return result
}
