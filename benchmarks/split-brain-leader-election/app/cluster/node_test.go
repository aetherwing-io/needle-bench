package cluster

import (
	"testing"
)

func TestNewNode(t *testing.T) {
	n := NewNode(42)
	if n.ID != 42 {
		t.Errorf("expected ID 42, got %d", n.ID)
	}
	if n.State != Follower {
		t.Errorf("expected Follower state, got %v", n.State)
	}
	if n.VotedFor != -1 {
		t.Errorf("expected VotedFor -1, got %d", n.VotedFor)
	}
}

func TestHandleVoteRequest_HigherTerm(t *testing.T) {
	n := NewNode(1)
	resp := n.HandleVoteRequest(2, 5)
	if !resp.Granted {
		t.Error("expected vote to be granted for higher term")
	}
	if n.Term != 5 {
		t.Errorf("expected term to be updated to 5, got %d", n.Term)
	}
}

func TestHandleVoteRequest_AlreadyVoted(t *testing.T) {
	n := NewNode(1)
	n.Term = 3
	n.VotedFor = 5

	resp := n.HandleVoteRequest(2, 3)
	if resp.Granted {
		t.Error("expected vote to be denied (already voted for another)")
	}
}

func TestNodeStateString(t *testing.T) {
	if Follower.String() != "follower" {
		t.Errorf("expected 'follower', got '%s'", Follower.String())
	}
	if Leader.String() != "leader" {
		t.Errorf("expected 'leader', got '%s'", Leader.String())
	}
}

// This test passes with the bug because it uses 3 nodes (no partition)
func TestElection_NoPartition(t *testing.T) {
	nodes := make([]*Node, 3)
	for i := 0; i < 3; i++ {
		nodes[i] = NewNode(i)
	}
	network := NewNetwork(nodes)

	peers := []*Node{nodes[1], nodes[2]}
	votes := nodes[0].StartElection(peers, network)

	if votes < 2 {
		t.Errorf("expected at least 2 votes, got %d", votes)
	}
	if !nodes[0].IsLeader() {
		t.Error("expected node 0 to be leader")
	}
}
