# deadlock-transfer

## Project

A Java banking application that supports concurrent fund transfers between accounts. Multiple threads process transfer requests simultaneously using synchronization for thread safety.

## Symptoms

When running concurrent transfers between accounts, the program hangs and eventually times out. The test runs 200 transfers across 8 threads with a 10-second timeout. The program exits with code 2 indicating the transfers never completed. Individual transfers work fine in isolation; the problem only occurs under concurrent load with transfers going in both directions between the same accounts.

## Bug description

The transfer logic acquires locks on two accounts to ensure atomicity, but the order in which locks are acquired depends on the direction of the transfer. When two threads try to transfer between the same pair of accounts in opposite directions simultaneously, they can each hold one lock while waiting for the other, creating a circular dependency that never resolves.

## Difficulty

Medium

## Expected turns

6-12
