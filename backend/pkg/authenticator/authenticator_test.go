package authenticator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestAuthenticator_VerifyAndExtract(t *testing.T) {
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	cases := []struct {
		name     string
		tokenTTL time.Duration
		wantErr  bool
	}{
		{
			name:     "valid token",
			tokenTTL: time.Hour,
			wantErr:  false,
		},
		{
			name:     "expired token",
			tokenTTL: -1 * time.Hour,
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			auth := NewAuthenticator(pk, pk, tc.tokenTTL, tc.tokenTTL)

			randomUserID, _ := uuid.NewV7()

			access, _, err := auth.SignTokens(auth.IssueTokens(randomUserID.String()))
			if err != nil {
				t.Fatalf("error signing tokens: %v", err)
			}

			claims, err := auth.VerifyAndExtract(access)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error, got none")
			}

			if err != nil {
				return
			}

			if sub, err := claims.GetSubject(); err != nil || sub != randomUserID.String() {
				t.Error("unexpected claims")
			}
		})
	}
}
