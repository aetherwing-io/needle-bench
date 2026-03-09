package main

// MigrationGuide documents the changes between v1 and v2 API.
//
// v2 Changes:
//   - Added: department, phone_number, is_active fields
//   - Removed: (none documented)
//
// Note: All v1 fields should be preserved in v2. The v2 API is
// supposed to be a superset of v1 — additive changes only.
// If any v1 fields are missing from v2, that's a regression.

const MigrationNotes = `
API v1 -> v2 Migration Guide

New fields in v2:
  - department (string): User's department
  - phone_number (string): User's phone number
  - is_active (bool): Whether user account is active

All existing v1 fields are preserved in v2.
Clients can safely upgrade from v1 to v2 without code changes.
`
