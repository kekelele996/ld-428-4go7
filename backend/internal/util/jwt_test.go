package util

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret-0123456789abcdef"
	token, err := GenerateToken(secret, "usr-1", "Artist", time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.UserID != "usr-1" || claims.Role != "Artist" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	token, _ := GenerateToken("secret-a", "usr-1", "Artist", time.Hour)
	if _, err := ParseToken("secret-b", token); err == nil {
		t.Fatal("expected wrong secret error")
	}
}

func TestParseTokenExpired(t *testing.T) {
	token, _ := GenerateToken("secret-a", "usr-1", "Artist", -time.Hour)
	if _, err := ParseToken("secret-a", token); err == nil {
		t.Fatal("expected expired token error")
	}
}
