package utilities

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	uuid "github.com/satori/go.uuid"
)

func IsNilInteface(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

func GetSha(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func GenerateUuid() string {
	u := uuid.NewV4()
	us := u.String()
	return us
}

func GenerateGenricUserId() string {
	rand.Seed(time.Now().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var userBrokerID string
	for i := 0; i < 4; i++ {
		userBrokerID += string(chars[rand.Intn(len(chars))])
	}
	for i := 0; i < 6; i++ {
		userBrokerID += fmt.Sprintf("%v0", rand.Intn(10))
	}
	userBrokerID = "EM-" + userBrokerID
	return userBrokerID
}

func GenOtp(num int) string {
	return RandStringBytesMaskImprSrc(num, NumberBytes)
}

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	NumberBytes = "0123456789"

	LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int, bytes string) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(bytes) {
			b[i] = bytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
