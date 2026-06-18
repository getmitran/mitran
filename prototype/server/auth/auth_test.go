package auth

import (
	"os"
	"testing"
	"time"
)

func init() {
	os.Setenv("MITRAN_JWT_SECRET", "test-secret-key-for-unit-tests")
	jwtSecret = []byte(os.Getenv("MITRAN_JWT_SECRET"))
}

func TestGenerateAndValidateToken(t *testing.T) {
	token, err := GenerateToken("user1", RoleAdmin, time.Hour)
	if err != nil { t.Fatal(err) }
	claims, err := ValidateToken(token)
	if err != nil { t.Fatal(err) }
	if claims.UserID != "user1" { t.Errorf("expected user1, got %s", claims.UserID) }
	if claims.Role != RoleAdmin { t.Errorf("expected admin, got %s", claims.Role) }
}

func TestValidateToken_Expired(t *testing.T) {
	token, _ := GenerateToken("user1", RoleViewer, -time.Hour)
	_, err := ValidateToken(token)
	if err == nil { t.Error("expected error for expired token") }
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("not.a.token")
	if err == nil { t.Error("expected error for invalid token") }
}

func TestValidateToken_TamperedSignature(t *testing.T) {
	token, _ := GenerateToken("user1", RoleAdmin, time.Hour)
	tampered := token[:len(token)-4] + "xxxx"
	_, err := ValidateToken(tampered)
	if err == nil { t.Error("expected error for tampered signature") }
}
