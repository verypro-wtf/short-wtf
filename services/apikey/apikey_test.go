package apikey

import (
	"fmt"
	"testing"

	"github.com/verypro-wtf/short-wtf/config"
)

func test_getConfig() config.ApiKeyConfig {
	return config.ApiKeyConfig{
		Prefix:        "short",
		EntropyLength: 50,
		ChecksumLength: 8,
	}
}

var Test_ApiKeys []string

var Test_InvalidKeys = []string{
	"short_15y14y151klrjk1hklsahjkahfjklaj",
	"short_145ih1iuaklhjfajkhfaijf",
	"short_15u1947109p4j1lk41oy3io1j4kl1bno1h",
	"iris_1ihfkalkfhaiolfjaaaalfalk",
	"short_aokahflkajfklahfioajfl;anfajklhfaojklaakolj",
}

func TestGeneration(t *testing.T) {
	KeyGen := New(test_getConfig())
	Test_ApiKeys = make([]string, 0, 50)
	for range 50 {
		apiKey, err := KeyGen.Generate()
		if err != nil {
			t.Error(err)
		}
		Test_ApiKeys = append(Test_ApiKeys, apiKey)
	}
}

func TestValidation(t *testing.T) {
	KeyGen := New(test_getConfig())
	for _, key := range Test_ApiKeys {
		_, err := KeyGen.Validate(key)
		if err != nil {
			fmt.Println(key)
			t.Error(err)
		}
	}
}

func TestInvalidKey(t *testing.T) {
	KeyGen := New(test_getConfig())
	for _, key := range Test_InvalidKeys {
		valid, err := KeyGen.Validate(key)
		if err == nil || valid {
			t.Error(err)
		}
	}
}
