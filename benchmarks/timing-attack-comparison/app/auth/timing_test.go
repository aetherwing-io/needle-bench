package auth

import (
	"testing"
	"time"
)

func TestCompareHashes_Correct(t *testing.T) {
	hash := hashPassword("test123")
	if !compareHashes(hash, hash) {
		t.Error("identical hashes should match")
	}
}

func TestCompareHashes_Mismatch(t *testing.T) {
	hash1 := hashPassword("test123")
	hash2 := hashPassword("test456")
	if compareHashes(hash1, hash2) {
		t.Error("different hashes should not match")
	}
}

func TestCompareHashes_DifferentLength(t *testing.T) {
	if compareHashes("abc", "abcd") {
		t.Error("different length strings should not match")
	}
}

func TestTimingSafety(t *testing.T) {
	// This test checks whether the comparison leaks timing information.
	// A constant-time comparison should take approximately the same time
	// regardless of how many characters match.

	target := hashPassword("secret-password-42")

	// Create test hashes with increasing prefix matches
	// Match 0 chars, 8 chars, 16 chars, 32 chars of the target hash
	testCases := []struct {
		name       string
		matchChars int
	}{
		{"no_match", 0},
		{"8_chars", 8},
		{"16_chars", 16},
		{"32_chars", 32},
	}

	results := make(map[string]time.Duration)

	for _, tc := range testCases {
		input := makePartialMatch(target, tc.matchChars)

		// Warm up
		for i := 0; i < 1000; i++ {
			compareHashes(target, input)
		}

		// Measure
		const iterations = 100000
		start := time.Now()
		for i := 0; i < iterations; i++ {
			compareHashes(target, input)
		}
		elapsed := time.Since(start)
		results[tc.name] = elapsed

		t.Logf("  %s: %v (%d iterations)", tc.name, elapsed, iterations)
	}

	// Check if timing varies significantly with match length
	// A vulnerable implementation will show increasing time with more matches
	noMatch := results["no_match"]
	fullPrefix := results["32_chars"]

	// If 32-char match takes >30% longer than 0-char match, timing leaks
	ratio := float64(fullPrefix) / float64(noMatch)
	t.Logf("  Timing ratio (32_chars / no_match): %.3f", ratio)

	if ratio > 1.3 {
		t.Errorf("TIMING LEAK: 32-char prefix match took %.1f%% longer than no match (ratio: %.3f)",
			(ratio-1)*100, ratio)
	}
}
