# split-brain-leader-election

## Project

A Go implementation of a leader election protocol for a 5-node distributed cluster. Nodes communicate via message passing and elect a single leader using a voting mechanism similar to Raft. The system includes a network simulation layer that can introduce partitions to test fault tolerance.

## Symptoms

When a network partition splits the cluster into two groups, the safety check reports that multiple nodes believe they are the leader simultaneously. The test runs a simulation where the cluster is partitioned and both sides attempt election. In a correct implementation, only the majority partition should be able to elect a leader. Instead, both partitions sometimes succeed.

## Bug description

The leader election quorum check has an off-by-one error in how it determines whether a candidate has received enough votes. The threshold calculation allows a minority partition to also elect a leader, violating the single-leader invariant. Understanding the relationship between cluster size, partition sizes, and integer arithmetic is required to pinpoint the issue.

## Difficulty

Hard

## Expected turns

8-12
