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
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user1" {
		t.Errorf("expected user1, got %s", claims.UserID)
	}
	if claims.Role != RoleAdmin {
		t.Errorf("expected admin, got %s", claims.Role)
	}
}

func TestValidateToken_Expired(t *testing.T) {
	token, _ := GenerateToken("user1", RoleViewer, -time.Hour)
	_, err := ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	_, err := ValidateToken("not.a.token")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestValidateToken_TamperedSignature(t *testing.T) {
	token, _ := GenerateToken("user1", RoleAdmin, time.Hour)
	tampered := token[:len(token)-4] + "xxxx"
	_, err := ValidateToken(tampered)
	if err == nil {
		t.Error("expected error for tampered signature")
	}
}

func TestIssueRefreshToken_NonEmpty(t *testing.T) {
	store := NewRefreshStore()
	token := store.IssueRefreshToken("test@example.com")
	if token == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestValidateRefresh_ValidToken(t *testing.T) {
	store := NewRefreshStore()
	token := store.IssueRefreshToken("test@example.com")
	email, err := store.ValidateRefresh(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if email != "test@example.com" {
		t.Errorf("expected test@example.com, got %s", email)
	}
}

func TestValidateRefresh_InvalidToken(t *testing.T) {
	store := NewRefreshStore()
	_, err := store.ValidateRefresh("bogus-token")
	if err == nil {
		t.Error("expected error for invalid refresh token")
	}
}

func TestRotateRefresh_InvalidatesOld(t *testing.T) {
	store := NewRefreshStore()
	old := store.IssueRefreshToken("rotate@example.com")
	_, _, err := store.RotateRefresh(old)
	if err != nil {
		t.Fatalf("unexpected error on rotate: %v", err)
	}
	_, err = store.ValidateRefresh(old)
	if err == nil {
		t.Error("expected error validating rotated token")
	}
}

func TestIsRevoked_AfterRevoke(t *testing.T) {
	token := "test-session-token-" + time.Now().Format("150405")
	if IsRevoked(token) {
		t.Fatal("token should not be revoked before Revoke()")
	}
	Revoke(token)
	if !IsRevoked(token) {
		t.Error("expected token to be revoked after Revoke()")
	}
}
