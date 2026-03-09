# timezone-scheduling

A Python scheduling service that manages recurring events for users across different timezones. The service stores events with their scheduled time and timezone, then determines which events should fire during a given time window. Events near midnight boundaries are being silently skipped -- a daily 11:30 PM event in US/Eastern never fires, while a noon event works fine. The scheduler processes all events but some just never match the current window.

## Symptoms

- Running `test.sh` shows that late-night events (near midnight) are missed
- The scheduler reports "0 events due" for time windows that should contain events
- Events scheduled during daytime hours work correctly
- The bug only manifests when the event's local time is in a different calendar day than UTC

## Bug description

The scheduler compares timestamps incorrectly when determining which events fall within the current execution window. The comparison logic mixes timezone representations, causing events whose local time and UTC time fall on different calendar days to be silently skipped.

## Difficulty

Easy

## Expected turns

3-6
