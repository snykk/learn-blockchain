package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// RaftState represents the state of a Raft node
type RaftState string

const (
	RaftFollower  RaftState = "follower"
	RaftCandidate RaftState = "candidate"
	RaftLeader    RaftState = "leader"
)

// RaftMessageType represents the type of Raft message
type RaftMessageType string

const (
	RaftRequestVote       RaftMessageType = "request_vote"
	RaftRequestVoteResp   RaftMessageType = "request_vote_response"
	RaftAppendEntries     RaftMessageType = "append_entries"
	RaftAppendEntriesResp RaftMessageType = "append_entries_response"
)

// RaftLogEntry represents a log entry in Raft
type RaftLogEntry struct {
	Index   int64  `json:"index"`
	Term    int64  `json:"term"`
	Command *Block `json:"command"` // Block to add to blockchain
}

// RaftMessage represents a message in Raft consensus
type RaftMessage struct {
	Type   RaftMessageType `json:"type"`
	Term   int64           `json:"term"`
	NodeID string          `json:"node_id"`
	From   string          `json:"from"`

	// For RequestVote
	LastLogIndex int64 `json:"last_log_index,omitempty"`
	LastLogTerm  int64 `json:"last_log_term,omitempty"`
	VoteGranted  bool  `json:"vote_granted,omitempty"`

	// For AppendEntries
	PrevLogIndex int64           `json:"prev_log_index,omitempty"`
	PrevLogTerm  int64           `json:"prev_log_term,omitempty"`
	Entries      []*RaftLogEntry `json:"entries,omitempty"`
	LeaderCommit int64           `json:"leader_commit,omitempty"`
	Success      bool            `json:"success,omitempty"`

	Signature string    `json:"signature"`
	Timestamp time.Time `json:"timestamp"`
}

// RaftNode represents a node in the Raft consensus
type RaftNode struct {
	ID            string
	Peers         []string
	State         RaftState
	CurrentTerm   int64
	VotedFor      string
	VotesReceived int
	Log           []*RaftLogEntry
	CommitIndex   int64
	LastApplied   int64
	LeaderID      string

	// Leader state
	NextIndex  []int64
	MatchIndex []int64

	// Election timeout
	ElectionTimeout   time.Duration
	LastHeartbeat     time.Time
	HeartbeatInterval time.Duration

	mu         sync.RWMutex
	Blockchain *Blockchain
}

// NewRaftNode creates a new Raft node
func NewRaftNode(nodeID string, peers []string, bc *Blockchain) *RaftNode {
	node := &RaftNode{
		ID:                nodeID,
		Peers:             peers,
		State:             RaftFollower,
		CurrentTerm:       0,
		VotedFor:          "",
		VotesReceived:     0,
		Log:               make([]*RaftLogEntry, 0),
		CommitIndex:       0,
		LastApplied:       0,
		LeaderID:          "",
		NextIndex:         make([]int64, len(peers)),
		MatchIndex:        make([]int64, len(peers)),
		ElectionTimeout:   time.Duration(150+rand.Intn(150)) * time.Millisecond, // 150-300ms
		HeartbeatInterval: 50 * time.Millisecond,
		LastHeartbeat:     time.Now(),
		Blockchain:        bc,
	}

	// Initialize next index for leader (will be set when becoming leader)
	for i := range node.NextIndex {
		node.NextIndex[i] = 1
	}

	return node
}

// RequestVote initiates leader election
func (rn *RaftNode) RequestVote() (*RaftMessage, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Transition to candidate state
	rn.State = RaftCandidate
	rn.CurrentTerm++
	rn.VotedFor = rn.ID
	rn.VotesReceived = 1

	// Get last log index and term
	lastLogIndex := int64(len(rn.Log))
	lastLogTerm := int64(0)
	if lastLogIndex > 0 {
		lastLogTerm = rn.Log[lastLogIndex-1].Term
	}

	msg := &RaftMessage{
		Type:         RaftRequestVote,
		Term:         rn.CurrentTerm,
		NodeID:       rn.ID,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  lastLogTerm,
		Timestamp:    time.Now(),
		Signature:    rn.signMessage(),
	}

	return msg, nil
}

// ProcessRequestVote processes a request vote message
func (rn *RaftNode) ProcessRequestVote(msg *RaftMessage) (*RaftMessage, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	voteGranted := false

	// If term is higher, update term and become follower
	if msg.Term > rn.CurrentTerm {
		rn.CurrentTerm = msg.Term
		rn.State = RaftFollower
		rn.VotedFor = ""
	}

	// Vote for candidate if:
	// 1. Candidate's term is at least as high as ours
	// 2. We haven't voted for anyone else in this term
	// 3. Candidate's log is at least as up-to-date as ours
	if msg.Term == rn.CurrentTerm && (rn.VotedFor == "" || rn.VotedFor == msg.NodeID) {
		lastLogIndex := int64(len(rn.Log))
		lastLogTerm := int64(0)
		if lastLogIndex > 0 {
			lastLogTerm = rn.Log[lastLogIndex-1].Term
		}

		// Check if candidate's log is at least as up-to-date
		if msg.LastLogTerm > lastLogTerm ||
			(msg.LastLogTerm == lastLogTerm && msg.LastLogIndex >= lastLogIndex) {
			voteGranted = true
			rn.VotedFor = msg.NodeID
			rn.LastHeartbeat = time.Now()
		}
	}

	resp := &RaftMessage{
		Type:        RaftRequestVoteResp,
		Term:        rn.CurrentTerm,
		NodeID:      rn.ID,
		From:        rn.ID,
		VoteGranted: voteGranted,
		Timestamp:   time.Now(),
		Signature:   rn.signMessage(),
	}

	return resp, nil
}

// ProcessRequestVoteResponse processes a response to request vote
func (rn *RaftNode) ProcessRequestVoteResponse(msg *RaftMessage) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// Ignore if term is outdated
	if msg.Term < rn.CurrentTerm {
		return nil
	}

	// If term is higher, update and become follower
	if msg.Term > rn.CurrentTerm {
		rn.CurrentTerm = msg.Term
		rn.State = RaftFollower
		rn.VotedFor = ""
		return nil
	}

	// Count vote if granted and we're still a candidate
	if rn.State == RaftCandidate && msg.VoteGranted {
		rn.VotesReceived++

		// Check if we won the election
		if rn.VotesReceived > len(rn.Peers)/2 {
			rn.BecomeLeader()
		}
	}

	return nil
}

// BecomeLeader transitions node to leader state
func (rn *RaftNode) BecomeLeader() {
	rn.State = RaftLeader
	rn.LeaderID = rn.ID

	// Initialize leader state
	lastLogIndex := int64(len(rn.Log))
	for i := range rn.NextIndex {
		rn.NextIndex[i] = lastLogIndex + 1
		rn.MatchIndex[i] = 0
	}

	fmt.Printf("  Node %s became LEADER for term %d\n", rn.ID[:16]+"...", rn.CurrentTerm)
}

// AppendEntries creates an append entries message (heartbeat or log replication)
func (rn *RaftNode) AppendEntries(entries []*RaftLogEntry, prevLogIndex, prevLogTerm, leaderCommit int64) (*RaftMessage, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if rn.State != RaftLeader {
		return nil, fmt.Errorf("only leader can send append entries")
	}

	msg := &RaftMessage{
		Type:         RaftAppendEntries,
		Term:         rn.CurrentTerm,
		NodeID:       rn.ID,
		PrevLogIndex: prevLogIndex,
		PrevLogTerm:  prevLogTerm,
		Entries:      entries,
		LeaderCommit: leaderCommit,
		Timestamp:    time.Now(),
		Signature:    rn.signMessage(),
	}

	return msg, nil
}

// ProcessAppendEntries processes an append entries message
func (rn *RaftNode) ProcessAppendEntries(msg *RaftMessage) (*RaftMessage, error) {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	success := false

	// If term is higher, update term and become follower
	if msg.Term > rn.CurrentTerm {
		rn.CurrentTerm = msg.Term
		rn.State = RaftFollower
		rn.VotedFor = ""
	}

	// Update leader
	if msg.Term == rn.CurrentTerm {
		rn.LeaderID = msg.NodeID
		rn.LastHeartbeat = time.Now()
	}

	// Check if log is consistent
	if msg.Term == rn.CurrentTerm {
		// Check if previous log entry matches
		if msg.PrevLogIndex == 0 || (msg.PrevLogIndex <= int64(len(rn.Log)) && rn.Log[msg.PrevLogIndex-1].Term == msg.PrevLogTerm) {
			success = true

			// Append new entries
			if len(msg.Entries) > 0 {
				// Find conflict
				for i, entry := range msg.Entries {
					logIndex := msg.PrevLogIndex + 1 + int64(i)
					if logIndex <= int64(len(rn.Log)) {
						// Check for conflict
						if rn.Log[logIndex-1].Term != entry.Term {
							// Remove conflicting and subsequent entries
							rn.Log = rn.Log[:logIndex-1]
							rn.Log = append(rn.Log, entry)
						}
					} else {
						// Append new entry
						rn.Log = append(rn.Log, entry)
					}
				}
			}

			// Update commit index
			if msg.LeaderCommit > rn.CommitIndex {
				lastLogIndex := int64(len(rn.Log))
				if msg.LeaderCommit < lastLogIndex {
					rn.CommitIndex = msg.LeaderCommit
				} else {
					rn.CommitIndex = lastLogIndex
				}

				// Apply committed entries to blockchain
				rn.applyCommittedEntries()
			}
		}
	}

	resp := &RaftMessage{
		Type:      RaftAppendEntriesResp,
		Term:      rn.CurrentTerm,
		NodeID:    rn.ID,
		From:      rn.ID,
		Success:   success,
		Timestamp: time.Now(),
		Signature: rn.signMessage(),
	}

	return resp, nil
}

// ProcessAppendEntriesResponse processes response to append entries
func (rn *RaftNode) ProcessAppendEntriesResponse(msg *RaftMessage, peerIndex int) error {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	// If term is higher, update term and become follower
	if msg.Term > rn.CurrentTerm {
		rn.CurrentTerm = msg.Term
		rn.State = RaftFollower
		rn.VotedFor = ""
		return nil
	}

	// Only process if we're leader and term matches
	if rn.State == RaftLeader && msg.Term == rn.CurrentTerm {
		if msg.Success {
			// Update match index and next index for this peer
			rn.MatchIndex[peerIndex] = rn.NextIndex[peerIndex] - 1
			rn.NextIndex[peerIndex]++

			// Try to commit more entries
			rn.updateCommitIndex()
		} else {
			// Decrement next index and retry
			if rn.NextIndex[peerIndex] > 1 {
				rn.NextIndex[peerIndex]--
			}
		}
	}

	return nil
}

// updateCommitIndex updates commit index based on match indices
func (rn *RaftNode) updateCommitIndex() {
	for n := rn.CommitIndex + 1; n <= int64(len(rn.Log)); n++ {
		count := 0
		for _, matchIndex := range rn.MatchIndex {
			if matchIndex >= n {
				count++
			}
		}

		// If majority of peers have replicated this entry
		if count > len(rn.Peers)/2 {
			// Only commit if entry is from current term
			if rn.Log[n-1].Term == rn.CurrentTerm {
				rn.CommitIndex = n
				rn.applyCommittedEntries()
			}
		}
	}
}

// applyCommittedEntries applies committed log entries to blockchain
func (rn *RaftNode) applyCommittedEntries() {
	for rn.LastApplied < rn.CommitIndex {
		rn.LastApplied++
		entry := rn.Log[rn.LastApplied-1]

		// Apply the block to blockchain
		if entry.Command != nil {
			// Check if block already exists
			blockExists := false
			for _, block := range rn.Blockchain.Blocks {
				if block.Hash == entry.Command.Hash {
					blockExists = true
					break
				}
			}

			if !blockExists {
				rn.Blockchain.Blocks = append(rn.Blockchain.Blocks, entry.Command)
				fmt.Printf("    Applied committed block #%d to blockchain\n", entry.Command.Index)
			}
		}
	}
}

// StartElection starts a leader election
func (rn *RaftNode) StartElection() error {
	fmt.Printf("\n  Starting leader election (term %d)...\n", rn.CurrentTerm)

	// Send request vote to all peers
	requestVoteMsg, err := rn.RequestVote()
	if err != nil {
		return fmt.Errorf("failed to create request vote: %v", err)
	}

	fmt.Printf("  → Node %s is requesting votes\n", rn.ID[:16]+"...")
	fmt.Printf("    Term: %d, Last log index: %d, Last log term: %d\n",
		requestVoteMsg.Term, requestVoteMsg.LastLogIndex, requestVoteMsg.LastLogTerm)

	// Simulate receiving votes from peers
	votes := 1 // Vote for self
	for _, peer := range rn.Peers {
		if peer == rn.ID {
			continue
		}

		// Simulate peer response
		voteMsg := &RaftMessage{
			Type:        RaftRequestVoteResp,
			Term:        rn.CurrentTerm,
			NodeID:      peer,
			From:        peer,
			VoteGranted: true, // Assume peers grant vote for simulation
			Timestamp:   time.Now(),
		}

		if err := rn.ProcessRequestVoteResponse(voteMsg); err != nil {
			fmt.Printf("    Warning: Failed to process vote from node %s: %v\n", peer[:16]+"...", err)
			continue
		}

		if voteMsg.VoteGranted {
			votes++
			fmt.Printf("    Received vote from node %s\n", peer[:16]+"...")
		}

		if rn.State == RaftLeader {
			break
		}
	}

	fmt.Printf("    Total votes: %d/%d (majority: %d)\n", votes, len(rn.Peers)+1, len(rn.Peers)/2+1)

	if rn.State != RaftLeader {
		return fmt.Errorf("failed to win election")
	}

	return nil
}

// ReplicateLog replicates log entries to followers
func (rn *RaftNode) ReplicateLog(block *Block) error {
	rn.mu.Lock()

	if rn.State != RaftLeader {
		rn.mu.Unlock()
		return fmt.Errorf("only leader can replicate log")
	}

	// Create log entry
	entry := &RaftLogEntry{
		Index:   int64(len(rn.Log)) + 1,
		Term:    rn.CurrentTerm,
		Command: block,
	}

	rn.Log = append(rn.Log, entry)

	// Get previous log index and term
	prevLogIndex := int64(len(rn.Log) - 1)
	prevLogTerm := int64(0)
	if prevLogIndex > 0 {
		prevLogTerm = rn.Log[prevLogIndex-1].Term
	}

	entries := []*RaftLogEntry{entry}
	rn.mu.Unlock()

	fmt.Printf("\n  Replicating block #%d to followers...\n", block.Index)

	// Replicate to all peers
	successCount := 0
	for i, peer := range rn.Peers {
		if peer == rn.ID {
			continue
		}

		_, err := rn.AppendEntries(entries, prevLogIndex, prevLogTerm, rn.CommitIndex)
		if err != nil {
			fmt.Printf("    Failed to replicate to node %s: %v\n", peer[:16]+"...", err)
			continue
		}

		// Simulate follower response
		resp := &RaftMessage{
			Type:      RaftAppendEntriesResp,
			Term:      rn.CurrentTerm,
			NodeID:    peer,
			From:      peer,
			Success:   true, // Assume success for simulation
			Timestamp: time.Now(),
		}

		if err := rn.ProcessAppendEntriesResponse(resp, i); err != nil {
			fmt.Printf("    Failed to process response from node %s: %v\n", peer[:16]+"...", err)
			continue
		}

		successCount++
		if successCount <= 3 {
			fmt.Printf("    Replicated to node %s\n", peer[:16]+"...")
		}
	}

	fmt.Printf("    Successfully replicated to %d/%d peers\n", successCount, len(rn.Peers))

	return nil
}

// SendHeartbeat sends heartbeat to followers
func (rn *RaftNode) SendHeartbeat() error {
	rn.mu.Lock()

	if rn.State != RaftLeader {
		rn.mu.Unlock()
		return fmt.Errorf("only leader can send heartbeat")
	}

	// Get previous log index and term
	prevLogIndex := int64(len(rn.Log))
	prevLogTerm := int64(0)
	if prevLogIndex > 0 {
		prevLogTerm = rn.Log[prevLogIndex-1].Term
	}

	rn.mu.Unlock()

	// Send empty append entries (heartbeat) to all peers
	for i, peer := range rn.Peers {
		if peer == rn.ID {
			continue
		}

		_, err := rn.AppendEntries([]*RaftLogEntry{}, prevLogIndex, prevLogTerm, rn.CommitIndex)
		if err != nil {
			continue
		}

		// Simulate follower response
		resp := &RaftMessage{
			Type:      RaftAppendEntriesResp,
			Term:      rn.CurrentTerm,
			NodeID:    peer,
			From:      peer,
			Success:   true,
			Timestamp: time.Now(),
		}

		rn.ProcessAppendEntriesResponse(resp, i)
	}

	return nil
}

// CheckElectionTimeout checks if election timeout has expired
func (rn *RaftNode) CheckElectionTimeout() bool {
	rn.mu.RLock()
	defer rn.mu.RUnlock()

	// Only followers and candidates check timeout
	if rn.State == RaftLeader {
		return false
	}

	return time.Since(rn.LastHeartbeat) > rn.ElectionTimeout
}

// signMessage creates a simple signature for a message
func (rn *RaftNode) signMessage() string {
	data := fmt.Sprintf("%s:%d:%d", rn.ID, rn.CurrentTerm, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GetStatus returns the current status of the Raft node
func (rn *RaftNode) GetStatus() string {
	rn.mu.RLock()
	defer rn.mu.RUnlock()

	return fmt.Sprintf("Node: %s, State: %s, Term: %d, Leader: %s, Log length: %d, CommitIndex: %d",
		rn.ID[:16]+"...", rn.State, rn.CurrentTerm, rn.LeaderID[:16]+"...", len(rn.Log), rn.CommitIndex)
}

// IsLeader checks if this node is the leader
func (rn *RaftNode) IsLeader() bool {
	rn.mu.RLock()
	defer rn.mu.RUnlock()
	return rn.State == RaftLeader
}
