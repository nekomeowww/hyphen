package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
)

// RandBytes 根据给定的长度生成字节，长度默认为 32
func RandBytes(length ...int) ([]byte, error) {
	b := make([]byte, 32)
	if len(length) != 0 {
		b = make([]byte, length[0])
	}
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// RandomHashString 生成随机 SHA128 字符串，最大长度为 128
func RandomHashString(length ...int) string {
	b, _ := RandBytes(1024)
	if len(length) != 0 {
		sliceLength := length[0]
		if length[0] > 64 {
			sliceLength = 64
		}
		if length[0] <= 0 {
			sliceLength = 64
		}

		return fmt.Sprintf("%x", sha512.Sum512(b))[:sliceLength]
	}

	return fmt.Sprintf("%x", sha512.Sum512(b))
}
