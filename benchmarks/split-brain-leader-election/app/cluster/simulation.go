package cluster

import (
	"fmt"
	"time"
)

// SimulationResult holds the outcome of a leader election simulation.
type SimulationResult struct {
	SplitBrain bool
	Leaders    []int
	PartitionA []int
	PartitionB []int
}

// RunSimulation runs a leader election scenario with a network partition.
func RunSimulation() SimulationResult {
	// Create a 5-node cluster
	nodes := make([]*Node, 5)
	for i := 0; i < 5; i++ {
		nodes[i] = NewNode(i)
	}

	network := NewNetwork(nodes)

	// Partition: nodes [0,1] and [2,3,4]
	partA := []int{0, 1}
	partB := []int{2, 3, 4}
	network.Partition(partA, partB)

	// Both sides try to elect a leader simultaneously
	// Node 0 campaigns in partition A
	// Node 2 campaigns in partition B
	done := make(chan struct{}, 2)

	go func() {
		peers := []*Node{nodes[1], nodes[2], nodes[3], nodes[4]}
		nodes[0].StartElection(peers, network)
		done <- struct{}{}
	}()

	go func() {
		peers := []*Node{nodes[0], nodes[1], nodes[3], nodes[4]}
		nodes[2].StartElection(peers, network)
		done <- struct{}{}
	}()

	// Wait for both elections
	<-done
	<-done

	// Small settle time
	time.Sleep(50 * time.Millisecond)

	// Check for split-brain
	var leaders []int
	for _, n := range nodes {
		if n.IsLeader() {
			leaders = append(leaders, n.ID)
		}
	}

	result := SimulationResult{
		PartitionA: partA,
		PartitionB: partB,
		Leaders:    leaders,
	}

	if len(leaders) > 1 {
		result.SplitBrain = true
	} else if len(leaders) == 0 {
		fmt.Println("WARNING: no leader elected")
	}

	return result
}

// CheckSafety runs multiple simulation rounds to check for split-brain.
func CheckSafety() bool {
	for i := 0; i < 20; i++ {
		result := RunSimulation()
		if result.SplitBrain {
			fmt.Printf("Split-brain found on round %d\n", i+1)
			return false
		}
	}
	return true
}
