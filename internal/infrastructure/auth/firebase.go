package auth

import (
	"context"
	"fmt"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/yohannes/sint-new-backend/internal/domain"
)

type FirebaseVerifier struct {
	client *firebaseauth.Client
}

func NewFirebaseVerifier(client *firebaseauth.Client) *FirebaseVerifier {
	return &FirebaseVerifier{client: client}
}

func (verifier *FirebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (domain.User, error) {
	token, err := verifier.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return domain.User{}, fmt.Errorf("verify Firebase ID token: %w", err)
	}
	email, _ := token.Claims["email"].(string)
	displayName, _ := token.Claims["name"].(string)
	return domain.User{
		FirebaseUID: token.UID,
		Email:       email,
		DisplayName: displayName,
	}, nil
}
