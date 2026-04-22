package utility

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	"vita-track-ai/models"
)

// Helper function to safely get strings from claims
func GetClaim(key string, claims map[string]interface{}) string {

	if val, ok := claims[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func GenerateOTP() models.OneTimePassword {
	var oneTimePassword models.OneTimePassword
	rand.Seed(time.Now().UnixNano())
	otpStr := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiry := time.Now().Add(5 * time.Minute)

	oneTimePassword.OTP = &otpStr
	oneTimePassword.OTPExpiresAt = &expiry

	return oneTimePassword
}

func EncodeCursor(c models.Cursor) (string, error) {
	data, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(data), nil
}

func DecodeCursor(cursorStr string) (*models.Cursor, error) {
	if cursorStr == "" {
		return nil, nil
	}

	data, err := base64.StdEncoding.DecodeString(cursorStr)
	if err != nil {
		return nil, err
	}

	var c models.Cursor
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
