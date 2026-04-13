package service

import (
	"strings"
	"vita-track-ai/models"
)

func BuildUserResponse(user models.User) (*models.UserResponse, error) {
	profilePicURL, err := resolveUserProfilePicture(user)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		UserId:     user.UserId,
		Email:      user.Email,
		GoogleId:   user.GoogleId,
		Name:       user.Name,
		Age:        user.Age,
		Gender:     user.Gender,
		ProfilePic: profilePicURL,
		DOB:        user.DOB,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}, nil
}

func resolveUserProfilePicture(user models.User) (*string, error) {
	if user.ProfileImage != nil {
		url, err := GenerateSignedURL(user.ProfileImage.Bucket, user.ProfileImage.ObjectKey)
		if err != nil {
			return nil, err
		}
		return &url, nil
	}

	if user.LegacyProfilePic == nil || *user.LegacyProfilePic == "" {
		return nil, nil
	}

	if isExternalURL(*user.LegacyProfilePic) {
		return user.LegacyProfilePic, nil
	}

	url, err := GenerateSignedURL(defaultBucketName(), *user.LegacyProfilePic)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func isExternalURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}
