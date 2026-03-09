package cluster

import (
	"fmt"
	"sync"
	"time"
)

// NodeState represents the state of a node in the cluster.
type NodeState int

const (
	Follower  NodeState = iota
	Candidate
	Leader
)

func (s NodeState) String() string {
	switch s {
	case Follower:
		return "follower"
	case Candidate:
		return "candidate"
	case Leader:
		return "leader"
	default:
		return "unknown"
	}
}

// Node represents a single node in the cluster.
type Node struct {
	ID       int
	State    NodeState
	Term     int
	VotedFor int // -1 means no vote cast this term
	mu       sync.Mutex

	// Communication channels
	inbox chan Message
}

// Message types for leader election protocol.
type MessageType int

const (
	VoteRequest MessageType = iota
	VoteResponse
	Heartbeat
)

// Message represents a message between nodes.
type Message struct {
	Type     MessageType
	From     int
	To       int
	Term     int
	Granted  bool
}

// NewNode creates a new cluster node.
func NewNode(id int) *Node {
	return &Node{
		ID:       id,
		State:    Follower,
		Term:     0,
		VotedFor: -1,
		inbox:    make(chan Message, 100),
	}
}

// StartElection initiates a leader election from this node.
func (n *Node) StartElection(peers []*Node, network *Network) int {
	n.mu.Lock()
	n.Term++
	n.State = Candidate
	n.VotedFor = n.ID
	currentTerm := n.Term
	n.mu.Unlock()

	fmt.Printf("Node %d: starting election for term %d\n", n.ID, currentTerm)

	votes := 1 // vote for self
	totalNodes := len(peers) + 1
	var votesMu sync.Mutex

	// Request votes from all peers
	var wg sync.WaitGroup
	for _, peer := range peers {
		wg.Add(1)
		go func(p *Node) {
			defer wg.Done()

			// Send vote request through network (may be partitioned)
			resp := network.SendVoteRequest(n.ID, p.ID, currentTerm)
			if resp != nil && resp.Granted {
				votesMu.Lock()
				votes++
				votesMu.Unlock()
			}
		}(peer)
	}

	// Wait for responses with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		fmt.Printf("Node %d: election timeout, proceeding with %d votes\n", n.ID, votes)
	}

	votesMu.Lock()
	finalVotes := votes
	votesMu.Unlock()

	fmt.Printf("Node %d: received %d/%d votes for term %d\n", n.ID, finalVotes, totalNodes, currentTerm)

	// Check if we have enough votes for a quorum
	if finalVotes >= totalNodes/2 {
		n.mu.Lock()
		n.State = Leader
		n.mu.Unlock()
		fmt.Printf("Node %d: became LEADER for term %d with %d votes\n", n.ID, finalVotes, currentTerm)
	}

	return finalVotes
}

// HandleVoteRequest processes an incoming vote request.
func (n *Node) HandleVoteRequest(from int, term int) *Message {
	n.mu.Lock()
	defer n.mu.Unlock()

	response := &Message{
		Type: VoteResponse,
		From: n.ID,
		To:   from,
		Term: n.Term,
	}

	if term > n.Term {
		// Higher term — update and grant vote
		n.Term = term
		n.VotedFor = from
		n.State = Follower
		response.Term = term
		response.Granted = true
		fmt.Printf("Node %d: granted vote to %d for term %d\n", n.ID, from, term)
	} else if term == n.Term && (n.VotedFor == -1 || n.VotedFor == from) {
		// Same term, haven't voted yet (or already voted for this candidate)
		n.VotedFor = from
		response.Granted = true
		fmt.Printf("Node %d: granted vote to %d for term %d\n", n.ID, from, term)
	} else {
		response.Granted = false
		fmt.Printf("Node %d: denied vote to %d for term %d (already voted for %d)\n",
			n.ID, from, term, n.VotedFor)
	}

	return response
}

// IsLeader returns true if this node believes it is the leader.
func (n *Node) IsLeader() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.State == Leader
}

// GetTerm returns the current term.
func (n *Node) GetTerm() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.Term
}
