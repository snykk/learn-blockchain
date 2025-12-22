package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// MessageType represents the type of message in the network
type MessageType string

const (
	MessageTypeBlockchain  MessageType = "blockchain"
	MessageTypeBlock       MessageType = "block"
	MessageTypeTransaction MessageType = "transaction"
	MessageTypePing        MessageType = "ping"
	MessageTypePong        MessageType = "pong"
)

// Message represents a message sent between nodes
type Message struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	From      string      `json:"from"`
}

// Node represents a blockchain node in the network
type Node struct {
	Address    string
	Port       int
	Blockchain *Blockchain
	Peers      map[string]bool // Map of peer addresses
	mu         sync.RWMutex
	listener   net.Listener
	running    bool
}

// NewNode creates a new node
func NewNode(address string, port int) *Node {
	return &Node{
		Address:    address,
		Port:       port,
		Blockchain: NewBlockchain(),
		Peers:      make(map[string]bool),
		running:    false,
	}
}

// AddPeer adds a peer to the node's peer list
func (n *Node) AddPeer(peerAddress string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Peers[peerAddress] = true
}

// RemovePeer removes a peer from the node's peer list
func (n *Node) RemovePeer(peerAddress string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.Peers, peerAddress)
}

// GetAddress returns the full address of the node
func (n *Node) GetAddress() string {
	return fmt.Sprintf("%s:%d", n.Address, n.Port)
}

// Start starts the node server
func (n *Node) Start() error {
	listener, err := net.Listen("tcp", n.GetAddress())
	if err != nil {
		return err
	}

	n.listener = listener
	n.running = true

	fmt.Printf("Node started on %s\n", n.GetAddress())

	go n.acceptConnections()

	return nil
}

// Stop stops the node server
func (n *Node) Stop() {
	n.running = false
	if n.listener != nil {
		n.listener.Close()
	}
}

// acceptConnections accepts incoming connections
func (n *Node) acceptConnections() {
	for n.running {
		conn, err := n.listener.Accept()
		if err != nil {
			if n.running {
				fmt.Printf("Error accepting connection: %v\n", err)
			}
			continue
		}

		go n.handleConnection(conn)
	}
}

// handleConnection handles an incoming connection
func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var msg Message

	if err := decoder.Decode(&msg); err != nil {
		fmt.Printf("Error decoding message: %v\n", err)
		return
	}

	n.processMessage(msg, conn)
}

// processMessage processes incoming messages
func (n *Node) processMessage(msg Message, conn net.Conn) {
	switch msg.Type {
	case MessageTypePing:
		// Respond with pong
		pong := Message{
			Type:      MessageTypePong,
			Timestamp: time.Now(),
			From:      n.GetAddress(),
		}
		n.sendMessage(pong, conn)

	case MessageTypeBlockchain:
		// Receive and merge blockchain data
		if receivedBlocks, err := n.parseBlocksFromMessage(msg); err == nil {
			fmt.Printf("Received blockchain data from %s (%d blocks)\n", msg.From, len(receivedBlocks))
			if err := n.Blockchain.MergeBlockchain(receivedBlocks); err != nil {
				fmt.Printf("Error merging blockchain: %v\n", err)
			}
		} else {
			fmt.Printf("Error parsing blockchain data: %v\n", err)
		}

	case MessageTypeBlock:
		// Receive and add new block
		if block, err := n.parseBlockFromMessage(msg); err == nil {
			fmt.Printf("Received new block #%d from %s\n", block.Index, msg.From)
			if err := n.Blockchain.AddReceivedBlock(block); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		} else {
			fmt.Printf("Error parsing block: %v\n", err)
		}

	case MessageTypeTransaction:
		// Receive and add new transaction to mempool
		if tx, err := n.parseTransactionFromMessage(msg); err == nil {
			fmt.Printf("Received new transaction from %s\n", msg.From)
			if err := n.Blockchain.AddTransactionToMempool(tx); err != nil {
				fmt.Printf("Error adding transaction to mempool: %v\n", err)
			} else {
				fmt.Printf("Transaction added to mempool: %s\n", tx.String())
			}
		} else {
			fmt.Printf("Error parsing transaction: %v\n", err)
		}
	}
}

// sendMessage sends a message to a connection
func (n *Node) sendMessage(msg Message, conn net.Conn) error {
	encoder := json.NewEncoder(conn)
	return encoder.Encode(msg)
}

// BroadcastBlockchain broadcasts the blockchain to all peers
func (n *Node) BroadcastBlockchain() {
	n.mu.RLock()
	peers := make([]string, 0, len(n.Peers))
	for peer := range n.Peers {
		peers = append(peers, peer)
	}
	n.mu.RUnlock()

	msg := Message{
		Type:      MessageTypeBlockchain,
		Data:      n.Blockchain.Blocks,
		Timestamp: time.Now(),
		From:      n.GetAddress(),
	}

	for _, peer := range peers {
		n.SendToPeer(peer, msg)
	}
}

// SendToPeer sends a message to a specific peer
func (n *Node) SendToPeer(peerAddress string, msg Message) error {
	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %s: %v", peerAddress, err)
	}
	defer conn.Close()

	return n.sendMessage(msg, conn)
}

// PingPeer pings a peer to check if it's alive
func (n *Node) PingPeer(peerAddress string) error {
	msg := Message{
		Type:      MessageTypePing,
		Timestamp: time.Now(),
		From:      n.GetAddress(),
	}

	return n.SendToPeer(peerAddress, msg)
}

// SyncBlockchain synchronizes blockchain with a peer
func (n *Node) SyncBlockchain(peerAddress string) error {
	msg := Message{
		Type:      MessageTypeBlockchain,
		Data:      n.Blockchain.Blocks,
		Timestamp: time.Now(),
		From:      n.GetAddress(),
	}

	return n.SendToPeer(peerAddress, msg)
}

// parseBlocksFromMessage parses blocks from a message
func (n *Node) parseBlocksFromMessage(msg Message) ([]*Block, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}

	var blocks []*Block
	if err := json.Unmarshal(dataBytes, &blocks); err != nil {
		return nil, err
	}

	return blocks, nil
}

// parseBlockFromMessage parses a block from a message
func (n *Node) parseBlockFromMessage(msg Message) (*Block, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}

	var block Block
	if err := json.Unmarshal(dataBytes, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

// parseTransactionFromMessage parses a transaction from a message
func (n *Node) parseTransactionFromMessage(msg Message) (*Transaction, error) {
	dataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		return nil, err
	}

	var tx Transaction
	if err := json.Unmarshal(dataBytes, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

// BroadcastBlock broadcasts a new block to all peers
func (n *Node) BroadcastBlock(block *Block) {
	n.mu.RLock()
	peers := make([]string, 0, len(n.Peers))
	for peer := range n.Peers {
		peers = append(peers, peer)
	}
	n.mu.RUnlock()

	msg := Message{
		Type:      MessageTypeBlock,
		Data:      block,
		Timestamp: time.Now(),
		From:      n.GetAddress(),
	}

	for _, peer := range peers {
		n.SendToPeer(peer, msg)
	}
}

// BroadcastTransaction broadcasts a new transaction to all peers
func (n *Node) BroadcastTransaction(tx *Transaction) {
	n.mu.RLock()
	peers := make([]string, 0, len(n.Peers))
	for peer := range n.Peers {
		peers = append(peers, peer)
	}
	n.mu.RUnlock()

	msg := Message{
		Type:      MessageTypeTransaction,
		Data:      tx,
		Timestamp: time.Now(),
		From:      n.GetAddress(),
	}

	for _, peer := range peers {
		n.SendToPeer(peer, msg)
	}
}
