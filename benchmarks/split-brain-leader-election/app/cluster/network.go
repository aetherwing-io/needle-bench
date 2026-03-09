package cluster

import (
	"fmt"
	"sync"
)

// Network simulates a network that can be partitioned.
type Network struct {
	nodes map[int]*Node
	// partition defines which nodes can communicate
	// if partition is nil, all nodes can communicate
	partitions [][]int
	mu         sync.RWMutex
}

// NewNetwork creates a new network with the given nodes.
func NewNetwork(nodes []*Node) *Network {
	nodeMap := make(map[int]*Node)
	for _, n := range nodes {
		nodeMap[n.ID] = n
	}
	return &Network{
		nodes: nodeMap,
	}
}

// Partition splits the network into two groups.
// Nodes in different groups cannot communicate.
func (net *Network) Partition(groupA, groupB []int) {
	net.mu.Lock()
	defer net.mu.Unlock()
	net.partitions = [][]int{groupA, groupB}
	fmt.Printf("Network partitioned: %v | %v\n", groupA, groupB)
}

// Heal removes all partitions.
func (net *Network) Heal() {
	net.mu.Lock()
	defer net.mu.Unlock()
	net.partitions = nil
	fmt.Println("Network healed")
}

// CanCommunicate returns true if node A can reach node B.
func (net *Network) CanCommunicate(a, b int) bool {
	net.mu.RLock()
	defer net.mu.RUnlock()

	if net.partitions == nil {
		return true
	}

	// Both must be in the same partition group
	for _, group := range net.partitions {
		aInGroup := false
		bInGroup := false
		for _, id := range group {
			if id == a {
				aInGroup = true
			}
			if id == b {
				bInGroup = true
			}
		}
		if aInGroup && bInGroup {
			return true
		}
	}

	return false
}

// SendVoteRequest sends a vote request from one node to another,
// respecting network partitions.
func (net *Network) SendVoteRequest(from, to, term int) *Message {
	if !net.CanCommunicate(from, to) {
		return nil // message dropped due to partition
	}

	net.mu.RLock()
	target, exists := net.nodes[to]
	net.mu.RUnlock()

	if !exists {
		return nil
	}

	return target.HandleVoteRequest(from, term)
}
