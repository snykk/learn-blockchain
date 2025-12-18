package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Basic Blockchain Implementation ===")

	// Create a new blockchain
	fmt.Println("1. Creating new blockchain...")
	bc := NewBlockchain()
	time.Sleep(1 * time.Second)

	// Add some blocks
	fmt.Println("\n2. Adding blocks to the blockchain...")
	bc.AddBlock("Transaction 1: Alice sends 10 coins to Bob")
	time.Sleep(1 * time.Second)

	bc.AddBlock("Transaction 2: Bob sends 5 coins to Charlie")
	time.Sleep(1 * time.Second)

	bc.AddBlock("Transaction 3: Charlie sends 3 coins to Alice")
	time.Sleep(1 * time.Second)

	// Display the blockchain
	fmt.Println("\n3. Displaying the blockchain:")
	fmt.Println("==========================================")
	bc.Print()

	// Validate the blockchain
	fmt.Println("\n4. Validating the blockchain...")
	if bc.IsValid() {
		fmt.Println("Blockchain is valid!")
	} else {
		fmt.Println("Blockchain is invalid!")
	}

	// Test tampering detection - Scenario 1: Hacker changes data without recalculating hash
	fmt.Println("\n5. Testing tampering detection...")

	// Save original data and hash for restoration
	originalData := bc.Blocks[2].Data
	originalHash := bc.Blocks[2].Hash

	fmt.Println("\n   Scenario 1: Hacker modifies data WITHOUT recalculating hash")
	fmt.Println("   Modifying data in Block #2...")
	bc.Blocks[2].Data = "Tampered data: Hacker steals 1000 coins"
	// Note: Hash is NOT recalculated

	fmt.Println("\n   Validating blockchain after tampering...")
	if bc.IsValid() {
		fmt.Println("   Blockchain validation passed (this should not happen!)")
	} else {
		fmt.Println("   Blockchain validation failed (tampering detected!)")
		fmt.Println("   Reason: Hash mismatch - stored hash doesn't match calculated hash")
	}

	// Restore the blockchain for scenario 2
	fmt.Println("\n   Restoring blockchain...")
	bc.Blocks[2].Data = originalData
	bc.Blocks[2].Hash = originalHash

	// Test tampering detection - Scenario 2: Hacker changes data AND recalculates hash
	fmt.Println("\n   Scenario 2: Hacker modifies data AND recalculates hash")
	fmt.Println("   Modifying data in Block #2...")
	bc.Blocks[2].Data = "Tampered data: Hacker steals 1000 coins"
	fmt.Println("   Recalculating hash (without mining)...")
	bc.Blocks[2].Hash = bc.Blocks[2].CalculateHash()

	fmt.Println("\n   Validating blockchain after tampering...")
	if bc.IsValid() {
		fmt.Println("   Blockchain validation passed (this should not happen!)")
	} else {
		fmt.Println("   Blockchain validation failed (tampering detected!)")
		fmt.Println("   Reason: Proof of work invalid - hash doesn't meet difficulty requirement")
	}

	fmt.Println("\n=== Demo Complete ===")
}
