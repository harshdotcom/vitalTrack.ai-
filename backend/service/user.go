package service

import "vita-track-ai/repository"

var MAX_ALLOWED_USER_STORAGE int64 = 104857600 //100MB in bytes

// var MAX_ALLOWED_USER_STORAGE int64 = 439870

func exceedStorageLimit(userId int64, fileSize int64) (bool, error) {
	userUsage, err := repository.GetCurrentStorageUsed(userId)
	if err != nil {
		return true, err
	}
	totalStorageUsed := userUsage.TotalStorageUsed

	if totalStorageUsed+fileSize > MAX_ALLOWED_USER_STORAGE {
		return true, nil
	}
	return false, nil
}
