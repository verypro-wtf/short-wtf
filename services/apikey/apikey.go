package apikey

import (
	"errors"
	"fmt"
	"strings"

	"github.com/verypro-wtf/short-wtf/config"
	"github.com/verypro-wtf/short-wtf/services/base62"
	"github.com/verypro-wtf/short-wtf/services/crc32"
)

var (
	ErrInvalidApiKeyLength    = errors.New("Invalid API Key length")
	ErrInvalidApiKeyPrefix    = errors.New("Invalid API Key prefix")
	ErrInvalidApiKeyChecksum  = errors.New("Invalid API Key checksum")
	ErrInvalidApiKeyStructure = errors.New("Invalid API Key structure")
)

type ApiKeyHandler interface {
	Generate() (string, error)
	Validate(rawKey string) (bool, error)
	Parse(rawKey string) (ApiKey, error)
	Format(key ApiKey) string
}

type ApiKeyHandle struct {
	prefix         string
	entropyLength  int
	checksumLength int
}

type ApiKey struct {
	Prefix   string
	Entropy  string
	Checksum string
}

func New(config config.ApiKeyConfig) ApiKeyHandler {
	return &ApiKeyHandle{
		prefix:         config.Prefix,
		entropyLength:  config.EntropyLength,
		checksumLength: config.ChecksumLength,
	}
}

func (gen ApiKeyHandle) Format(apiKey ApiKey) string {
	return fmt.Sprintf("%s_%s%s", apiKey.Prefix, apiKey.Entropy, apiKey.Checksum)
}

func (gen ApiKeyHandle) Parse(key string) (ApiKey, error) {
	if len(key) != (len(gen.prefix) + 1 + gen.entropyLength + gen.checksumLength) {
		return ApiKey{}, ErrInvalidApiKeyLength
	}

	if !strings.HasPrefix(key, gen.prefix) {
		return ApiKey{}, ErrInvalidApiKeyPrefix
	}

	items := strings.Split(key, "_")
	if len(items) != 2 {
		return ApiKey{}, ErrInvalidApiKeyStructure
	}

	return ApiKey{
		Prefix:   items[0],
		Entropy:  items[1][:gen.entropyLength],
		Checksum: items[1][gen.entropyLength:],
	}, nil
}

func (gen ApiKeyHandle) Generate() (string, error) {
	entropy, err := base62.GenerateRandom(gen.entropyLength)
	if err != nil {
		return "", err
	}
	checksum := crc32.CreateChecksum(entropy)
	apiKey := ApiKey{
		Prefix:   gen.prefix,
		Entropy:  entropy,
		Checksum: checksum,
	}
	fmt.Println(len(checksum))
	return gen.Format(apiKey), nil
}

func (gen ApiKeyHandle) Validate(rawKey string) (bool, error) {
	apiKey, err := gen.Parse(rawKey)
	if err != nil {
		return false, err
	}

	expectedChecksum := crc32.CreateChecksum(apiKey.Entropy)
	if expectedChecksum != apiKey.Checksum {
		return false, ErrInvalidApiKeyChecksum
	}

	return true, nil
}
