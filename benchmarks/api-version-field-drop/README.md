# api-version-field-drop

## Project

A Go HTTP API server with two versioned endpoints: /v1/users and /v2/users. The v2 API was introduced to add new fields (department, phone_number, is_active) while maintaining backward compatibility with v1. Existing clients should be able to upgrade from v1 to v2 without breaking.

## Symptoms

Clients that upgrade from /v1/users to /v2/users find that certain fields they depended on are no longer present in the response. The v2 response has the new fields (department, phone_number, is_active) but is missing some fields that existed in v1. No error is returned — the fields are simply absent from the JSON response. The migration documentation says v2 is a superset of v1 with only additive changes.

## Bug description

When the v2 response model was created, some fields from the v1 model were accidentally omitted. The conversion function faithfully maps only the fields defined in the v2 struct, so the omitted fields silently disappear. The underlying data store has all the information; it is the v2 serialization layer that drops them.

## Difficulty

Medium

## Expected turns

6-12
