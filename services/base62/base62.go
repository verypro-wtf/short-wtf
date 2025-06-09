package base62

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var (
	encodeMap = []rune(charset)
	decodeMap = make(map[rune]int)
	base      = big.NewInt(int64(len(charset)))
)

func init() {
	for i, c := range charset {
		decodeMap[c] = i
	}
}

func Encode(input []byte) string {
	num := new(big.Int).SetBytes(input)
	if num.Sign() == 0 {
		return string(encodeMap[0])
	}

	var result []rune
	mod := new(big.Int)
	for num.Sign() > 0 {
		num.DivMod(num, base, mod)
		result = append([]rune{encodeMap[mod.Int64()]}, result...)
	}

	return string(result)
}

func Decode(encoded string) (string, error) {
	if len(encoded) == 0 {
		return "", errors.New("input cannot be empty")
	}
	num := big.NewInt(0)

	for _, c := range encoded {
		val, ok := decodeMap[c]
		if !ok {
			return "", errors.New("invalid character in input: " + string(c))
		}
		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(val)))
	}

	return string(num.Bytes()), nil
}

func GenerateRandom(length int) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(charset)))

	for i := range length {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
