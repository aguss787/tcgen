package rand

import "math/rand"

const charset = "abcdefghijklmnopqrstuvwxyz"

var seededRand *rand.Rand = rand.New(rand.NewSource(91825479412))

func Intn(n int) int {
	return seededRand.Intn(n)
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func Shuffle(length int, f func(i, j int)) {
	seededRand.Shuffle(length, f)
}