package auth

import (
	"testing"
)

func TestGenerateTokenUniqueness(t *testing.T) {
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := generateToken()
		if err != nil {
			t.Fatalf("generating token: %v", err)
		}
		if len(token) != 64 {
			t.Errorf("token length = %d, want 64", len(token))
		}
		if tokens[token] {
			t.Errorf("duplicate token generated: %s", token)
		}
		tokens[token] = true
	}
}
