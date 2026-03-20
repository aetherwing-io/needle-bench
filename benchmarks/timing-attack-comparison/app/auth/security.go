package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

// makePartialMatch creates a hash-length string that matches the first n chars of target.
func makePartialMatch(target string, matchChars int) string {
	if matchChars <= 0 {
		// Completely different hash
		h := sha256.Sum256([]byte("completely-different-input"))
		return hex.EncodeToString(h[:])
	}

	result := make([]byte, len(target))
	copy(result, []byte(target[:matchChars]))
	// Fill rest with characters guaranteed to not match
	for i := matchChars; i < len(result); i++ {
		if target[i] == 'x' {
			result[i] = 'y'
		} else {
			result[i] = 'x'
		}
	}
	return string(result)
}

// RunSecurityTest is the exported test function called from main.
func RunSecurityTest() bool {
	target := hashPassword("secret-password-42")

	matchLengths := []int{0, 4, 8, 16, 24, 32, 48, 60}
	type measurement struct {
		matchLen int
		duration time.Duration
	}

	var measurements []measurement

	fmt.Println("Running timing analysis...")

	for _, matchLen := range matchLengths {
		input := makePartialMatch(target, matchLen)

		// Warm up
		for i := 0; i < 5000; i++ {
			compareHashes(target, input)
		}

		// Measure over many iterations for statistical significance
		const iterations = 500000
		start := time.Now()
		for i := 0; i < iterations; i++ {
			compareHashes(target, input)
		}
		elapsed := time.Since(start)

		measurements = append(measurements, measurement{matchLen, elapsed})
		fmt.Printf("  Match %2d chars: %v\n", matchLen, elapsed)
	}

	// Sort by match length
	sort.Slice(measurements, func(i, j int) bool {
		return measurements[i].matchLen < measurements[j].matchLen
	})

	// Check for timing correlation: if more matching characters = more time,
	// the comparison is vulnerable to timing attacks
	baseline := measurements[0].duration
	worst := measurements[len(measurements)-1].duration
	ratio := float64(worst) / float64(baseline)

	fmt.Printf("\nTiming ratio (max_match / no_match): %.3f\n", ratio)

	// A constant-time comparison should have a ratio close to 1.0
	// A vulnerable comparison (using ==) will show increasing time with more matches
	if ratio > 1.2 {
		fmt.Printf("VULNERABLE: timing varies by %.1f%% based on match length\n", (ratio-1)*100)
		fmt.Println("This enables a timing side-channel attack on the password hash")
		return false
	}

	fmt.Println("OK: comparison appears to be constant-time")
	return true
}
