package utility

import (
	"context"
	"errors"
	"os"
	"strings"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDTokenAndGetPayload(token string) (*idtoken.Payload, error) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	if clientID == "" {
		return nil, errors.New("google login is not configured")
	}

	payload, err := idtoken.Validate(
		context.Background(),
		token,
		clientID,
	)

	if err != nil {
		return nil, errors.New("invalid google token")
	}

	return payload, nil
}
