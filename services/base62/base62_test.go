package base62

import "testing"

var (
	EncodeStringMap = map[string]string{
		"abcdefg":    "CBhsccspaX",
		"testing":    "CaDSseieGr",
		"abc1234567": "CShsiFFZWUAwJJ",
	}
	InvalidInput = []string{
		"abc$%^&*()",
		"abc 123",
		"",
	}
)

func TestEncode(t *testing.T) {
	for in, out := range EncodeStringMap {
		encoded := Encode([]byte(in))
		if encoded != out {
			t.Errorf("Expected %s, got %s for input %s", out, encoded, in)
		}
	}
}

func TestDecode(t *testing.T) {
	for in, out := range EncodeStringMap {
		decoded, err := Decode(out)
		if err != nil {
			t.Errorf("Error decoding %s: %v", out, err)
			continue
		}
		if decoded != in {
			t.Errorf("Expected %s, got %s for encoded %s", in, decoded, out)
		}
	}
}

func TestDecodeInvalid(t *testing.T) {
	for _, in := range InvalidInput {
		_, err := Decode(in)
		if err == nil {
			t.Errorf("Expected error for invalid input %s, but got none", in)
		}
	}
}

func TestGenerateRandom(t *testing.T) {
	length := 10
	for i := 0; i < 10; i++ {
		randomStr, err := GenerateRandom(length)
		if err != nil {
			t.Errorf("Error generating random string of length %d: %v", length, err)
			continue
		}
		if len(randomStr) != length {
			t.Errorf("Generated string length %d does not match expected length %d", len(randomStr), length)
		}
	}
}
