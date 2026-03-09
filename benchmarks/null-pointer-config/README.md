# null-pointer-config

## Project

A Go HTTP service that serves data with optional metrics reporting, configured via a JSON config file.

## Symptoms

When the config file enables the metrics feature (`enable_metrics: true`), the server crashes with a nil pointer dereference panic on the first request to the `/status` endpoint. The `/health` and `/data` endpoints may work fine, but `/status` crashes the process.

## Bug description

The configuration loading logic does not ensure all required sub-structures are initialized when their corresponding feature flags are enabled. The server compiles and starts successfully, but accessing a feature that depends on an uninitialized config section causes a runtime panic.

## Difficulty

Easy

## Expected turns

3-5
