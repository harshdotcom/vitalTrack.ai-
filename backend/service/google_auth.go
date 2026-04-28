package service

import (
	"errors"
	"strings"
	"vita-track-ai/models"
	"vita-track-ai/repository"
	"vita-track-ai/utility"

	"gorm.io/gorm"
)

var (
	ErrGoogleEmailUnavailable   = errors.New("google account did not provide an email address")
	ErrGoogleEmailNotVerified   = errors.New("google account email is not verified")
	ErrGoogleTokenNotConfigured = errors.New("google login is not configured")
	ErrInvalidGoogleToken       = errors.New("invalid google token")
)

func AuthenticateGoogleUser(idToken string) (*models.UserResponse, string, error) {
	payload, err := utility.VerifyGoogleIDTokenAndGetPayload(idToken)
	if err != nil {
		switch err.Error() {
		case ErrGoogleTokenNotConfigured.Error():
			return nil, "", ErrGoogleTokenNotConfigured
		case ErrInvalidGoogleToken.Error():
			return nil, "", ErrInvalidGoogleToken
		default:
			return nil, "", err
		}
	}

	claims := payload.Claims
	email := strings.TrimSpace(utility.GetClaim("email", claims))
	if email == "" {
		return nil, "", ErrGoogleEmailUnavailable
	}

	if !utility.GetBoolClaim("email_verified", claims) {
		return nil, "", ErrGoogleEmailNotVerified
	}

	name := strings.TrimSpace(utility.GetClaim("name", claims))
	picture := strings.TrimSpace(utility.GetClaim("picture", claims))
	googleID := strings.TrimSpace(payload.Subject)

	userModel, err := repository.GetUserModelByEmail(email)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		userModel = buildGoogleUser(email, name, picture, googleID)
		if _, err := repository.SaveUser(&userModel); err != nil {
			return nil, "", err
		}
	case err != nil:
		return nil, "", err
	default:
		if shouldSyncGoogleAccount(&userModel, name, picture, googleID) {
			if err := repository.UpdateGoogleAccount(&userModel); err != nil {
				return nil, "", err
			}
		}
	}

	token, err := utility.GenerateToken(userModel.Email, userModel.UserId)
	if err != nil {
		return nil, "", err
	}

	userResponse, err := BuildUserResponse(userModel)
	if err != nil {
		return nil, "", err
	}

	return userResponse, token, nil
}

func buildGoogleUser(email string, name string, picture string, googleID string) models.User {
	user := models.User{
		Email:      email,
		Name:       resolveGoogleDisplayName(name, email),
		GoogleId:   optionalStringPointer(googleID),
		IsVerified: true,
	}

	if picture != "" {
		user.LegacyProfilePic = &picture
	}

	return user
}

func shouldSyncGoogleAccount(userModel *models.User, name string, picture string, googleID string) bool {
	updated := false

	if googleID != "" && (userModel.GoogleId == nil || *userModel.GoogleId != googleID) {
		userModel.GoogleId = &googleID
		updated = true
	}

	if !userModel.IsVerified {
		userModel.IsVerified = true
		updated = true
	}

	if strings.TrimSpace(userModel.Name) == "" {
		userModel.Name = resolveGoogleDisplayName(name, userModel.Email)
		updated = true
	}

	if picture != "" && userModel.ProfileImage == nil && (userModel.LegacyProfilePic == nil || strings.TrimSpace(*userModel.LegacyProfilePic) == "") {
		userModel.LegacyProfilePic = &picture
		updated = true
	}

	return updated
}

func resolveGoogleDisplayName(name string, email string) string {
	trimmedName := strings.TrimSpace(name)
	if trimmedName != "" {
		return trimmedName
	}

	localPart, _, found := strings.Cut(email, "@")
	if found && strings.TrimSpace(localPart) != "" {
		return localPart
	}

	return email
}

func optionalStringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
