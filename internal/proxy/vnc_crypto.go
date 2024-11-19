package proxy

import (
	"crypto/des"
)

func vncEncryptChallenge(challenge []byte, password string) []byte {
	// 确保密码长度为8字节
	key := make([]byte, 8)
	copy(key, []byte(password))

	// 反转每个字节的比特位
	for i := range key {
		key[i] = reverseBits(key[i])
	}

	response := make([]byte, 16)

	// 使用密码加密挑战的两个8字节块
	block, err := des.NewCipher(key)
	if err != nil {
		return nil
	}

	block.Encrypt(response[0:8], challenge[0:8])
	block.Encrypt(response[8:16], challenge[8:16])

	return response
}

func reverseBits(b byte) byte {
	var result byte
	for i := 0; i < 8; i++ {
		result = (result << 1) | (b & 1)
		b >>= 1
	}
	return result
}
