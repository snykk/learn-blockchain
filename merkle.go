package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	Root *MerkleNode
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
	Hash  []byte
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{
		Left:  left,
		Right: right,
		Data:  data,
	}

	if node.Left == nil && node.Right == nil {
		// Leaf node: hash the data
		hash := sha256.Sum256(data)
		node.Hash = hash[:]
	} else {
		// Internal node: hash the concatenation of left and right children
		prevHashes := append(node.Left.Hash, node.Right.Hash...)
		hash := sha256.Sum256(prevHashes)
		node.Hash = hash[:]
	}

	return node
}

// NewMerkleTree creates a new Merkle tree from transactions
func NewMerkleTree(transactions []*Transaction) *MerkleTree {
	if len(transactions) == 0 {
		return &MerkleTree{Root: nil}
	}

	var nodes []*MerkleNode

	// Create leaf nodes from transactions
	for _, tx := range transactions {
		txHash := tx.Hash()
		node := NewMerkleNode(nil, nil, txHash)
		nodes = append(nodes, node)
	}

	// Build tree from bottom up
	for len(nodes) > 1 {
		var level []*MerkleNode

		// Process pairs of nodes
		for i := 0; i < len(nodes); i += 2 {
			var left, right *MerkleNode
			left = nodes[i]

			if i+1 < len(nodes) {
				right = nodes[i+1]
			} else {
				// Odd number of nodes: duplicate the last node
				right = nodes[i]
			}

			// Create parent node
			parent := NewMerkleNode(left, right, nil)
			level = append(level, parent)
		}

		nodes = level
	}

	return &MerkleTree{Root: nodes[0]}
}

// GetRootHash returns the root hash as a hex string
func (mt *MerkleTree) GetRootHash() string {
	if mt.Root == nil {
		return ""
	}
	return hex.EncodeToString(mt.Root.Hash)
}
