# kernel-panic-ioctl

## Project

A C program (`devctl`) that simulates a kernel device driver's ioctl interface. The program manages a virtual device with a key-value configuration store. It accepts ioctl-style commands to set and get configuration values, mimicking the pattern used in real kernel drivers where user-space pointers must be validated before dereferencing.

## Symptoms

The self-test (normal operations with valid inputs) passes without issues. When running the fuzz test with crafted inputs including NULL pointers and oversized length fields, the program crashes with a segmentation fault. In a real kernel driver, this would be a kernel panic triggered by dereferencing an unvalidated user-space pointer. The crash occurs during specific ioctl operations when certain fields in the request structure contain invalid values.

## Bug description

The ioctl handler validates some input pointers but not others. While the top-level dispatch correctly checks for NULL device and request pointers, the individual command handlers have inconsistent validation. Some pointer fields are checked before use, while others are dereferenced directly. An attacker who controls the request structure can trigger a crash by setting specific fields to NULL. The fix requires ensuring all pointer fields are validated before any dereference, similar to how a kernel driver must use `copy_from_user`/`copy_to_user` for every user-space pointer access.

## Difficulty

Hard

## Expected turns

8-12
