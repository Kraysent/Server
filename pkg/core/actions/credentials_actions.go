package actions

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func GenerateToken(seed *int) string {
	if seed != nil {
		rand.Seed(int64(*seed))
	}

	rnd := rand.Int()
	tokenBytes := md5.Sum([]byte(fmt.Sprint(rnd)))

	return fmt.Sprintf("%x", tokenBytes)
}

func HashPassword(password string, salt int) string {
	passwordHashHex := md5.Sum([]byte(fmt.Sprintf("%s%d", password, salt)))

	return fmt.Sprintf("%x", passwordHashHex)
}

func GenerateSalt(seed *int) int {
	if seed != nil {
		rand.Seed(int64(*seed))
	}

	return rand.Intn(1000000)
}
