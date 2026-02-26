package utility

// Helper function to safely get strings from claims
func GetClaim(key string, claims map[string]interface{}) string {

	if val, ok := claims[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
