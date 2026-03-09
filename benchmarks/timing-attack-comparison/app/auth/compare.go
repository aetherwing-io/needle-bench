package auth

// compareHashes compares two hex-encoded hash strings for equality.
func compareHashes(stored, input string) bool {
	if len(stored) != len(input) {
		return false
	}
	return stored == input
}
