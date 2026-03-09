# timing-attack-comparison

## Project

A Go authentication server (`authserver`) that stores user credentials as SHA-256 hashes and verifies passwords against the stored hashes. The server exposes an HTTP login endpoint and a CLI for testing. User credentials are stored in an in-memory store with thread-safe access.

## Symptoms

Basic authentication works correctly: valid passwords are accepted, invalid passwords are rejected. However, the security test reveals that the time taken to reject a password varies depending on how many characters of the input hash match the stored hash. An input whose hash shares a longer prefix with the stored hash takes measurably longer to compare. The timing ratio between no-match and partial-match cases exceeds the safety threshold, indicating a side-channel vulnerability.

## Bug description

The password hash comparison uses a method that short-circuits on the first differing byte. This leaks information about the stored hash through response timing. An attacker can iteratively guess the hash one character at a time by measuring which inputs take longer to be rejected. The comparison needs to be replaced with a constant-time algorithm that always examines every byte regardless of match position.

## Difficulty

Hard

## Expected turns

8-12
