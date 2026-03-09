# import-cycle-startup

A Python microservice that manages user accounts and sends notification emails. The service has separate modules for user management and notification delivery. On cold start (first import), the application crashes with an `AttributeError` -- a module attribute that should exist is reported as missing. The error is confusing because the attribute clearly exists in the source code. The crash happens before any request is processed.

## Symptoms

- Running `test.sh` shows an `AttributeError` on startup
- The traceback references a module-level attribute that visibly exists in the source
- The error only occurs on cold start (fresh Python interpreter)
- Individual modules appear correct when read in isolation

## Bug description

Two modules reference each other at import time, creating a circular dependency. Python's import machinery partially initializes one module before the other finishes loading, causing an attribute lookup to fail because the attribute hasn't been defined yet at the point it's accessed. The fix requires breaking the import cycle.

## Difficulty

Easy

## Expected turns

3-5
