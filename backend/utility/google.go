package utility

import (
	"context"
	"errors"
	"os"

	"google.golang.org/api/idtoken"
)

func VerifyGoogkeIDTokenAndGetPayLoad(token string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(
		context.Background(),
		token,
		os.Getenv("GOOGLE_CLIENT_ID"),
	)

	if err != nil {
		return nil, errors.New("Invalid google token")
	}

	return payload, nil
}
