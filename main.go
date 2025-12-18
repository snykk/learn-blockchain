package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Enhanced Blockchain Implementation ===")
	fmt.Println("Features: Transactions, Merkle Tree, Wallet & Signing, Balance System")
	fmt.Println()

	// Create wallets
	fmt.Println("1. Creating wallets...")
	aliceWallet, err := NewWallet()
	if err != nil {
		fmt.Printf("Error creating Alice's wallet: %v\n", err)
		return
	}
	fmt.Printf("   Alice's wallet: %s\n", aliceWallet.Address)

	bobWallet, err := NewWallet()
	if err != nil {
		fmt.Printf("Error creating Bob's wallet: %v\n", err)
		return
	}
	fmt.Printf("   Bob's wallet: %s\n", bobWallet.Address)

	charlieWallet, err := NewWallet()
	if err != nil {
		fmt.Printf("Error creating Charlie's wallet: %v\n", err)
		return
	}
	fmt.Printf("   Charlie's wallet: %s\n", charlieWallet.Address)
	time.Sleep(1 * time.Second)

	// Create a new blockchain
	fmt.Println("\n2. Creating new blockchain...")
	bc := NewBlockchain()
	time.Sleep(1 * time.Second)

	// Give initial balances using coinbase transactions
	fmt.Println("\n3. Distributing initial balances (coinbase transactions)...")
	coinbase1 := bc.AddCoinbaseTransaction(aliceWallet.Address, 100.0)
	coinbase2 := bc.AddCoinbaseTransaction(bobWallet.Address, 50.0)
	coinbase3 := bc.AddCoinbaseTransaction(charlieWallet.Address, 30.0)

	if err := bc.AddBlock([]*Transaction{coinbase1, coinbase2, coinbase3}); err != nil {
		fmt.Printf("Error adding coinbase block: %v\n", err)
		return
	}

	fmt.Printf("   Alice received: 100.0 coins\n")
	fmt.Printf("   Bob received: 50.0 coins\n")
	fmt.Printf("   Charlie received: 30.0 coins\n")
	time.Sleep(1 * time.Second)

	// Display balances
	fmt.Println("\n4. Current balances:")
	fmt.Printf("   Alice: %.2f coins\n", bc.GetBalance(aliceWallet.Address))
	fmt.Printf("   Bob: %.2f coins\n", bc.GetBalance(bobWallet.Address))
	fmt.Printf("   Charlie: %.2f coins\n", bc.GetBalance(charlieWallet.Address))
	time.Sleep(1 * time.Second)

	// Create and sign transactions
	fmt.Println("\n5. Creating and signing transactions...")

	// Transaction 1: Alice sends 10 coins to Bob
	tx1 := NewTransaction(aliceWallet.Address, bobWallet.Address, 10.0)
	if err := bc.ValidateTransaction(tx1); err != nil {
		fmt.Printf("Error: Transaction 1 is invalid: %v\n", err)
		return
	}
	if err := aliceWallet.SignTransaction(tx1); err != nil {
		fmt.Printf("Error signing transaction 1: %v\n", err)
		return
	}
	fmt.Printf("   Transaction 1: %s\n", tx1.String())

	// Transaction 2: Bob sends 5 coins to Charlie
	tx2 := NewTransaction(bobWallet.Address, charlieWallet.Address, 5.0)
	if err := bc.ValidateTransaction(tx2); err != nil {
		fmt.Printf("Error: Transaction 2 is invalid: %v\n", err)
		return
	}
	if err := bobWallet.SignTransaction(tx2); err != nil {
		fmt.Printf("Error signing transaction 2: %v\n", err)
		return
	}
	fmt.Printf("   Transaction 2: %s\n", tx2.String())

	// Transaction 3: Charlie sends 3 coins to Alice
	tx3 := NewTransaction(charlieWallet.Address, aliceWallet.Address, 3.0)
	if err := bc.ValidateTransaction(tx3); err != nil {
		fmt.Printf("Error: Transaction 3 is invalid: %v\n", err)
		return
	}
	if err := charlieWallet.SignTransaction(tx3); err != nil {
		fmt.Printf("Error signing transaction 3: %v\n", err)
		return
	}
	fmt.Printf("   Transaction 3: %s\n", tx3.String())
	time.Sleep(1 * time.Second)

	// Add blocks with transactions
	fmt.Println("\n6. Adding blocks to the blockchain...")
	if err := bc.AddBlock([]*Transaction{tx1}); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	if err := bc.AddBlock([]*Transaction{tx2}); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	if err := bc.AddBlock([]*Transaction{tx3}); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	// Display balances after transactions
	fmt.Println("\n7. Balances after transactions:")
	fmt.Printf("   Alice: %.2f coins\n", bc.GetBalance(aliceWallet.Address))
	fmt.Printf("   Bob: %.2f coins\n", bc.GetBalance(bobWallet.Address))
	fmt.Printf("   Charlie: %.2f coins\n", bc.GetBalance(charlieWallet.Address))
	time.Sleep(1 * time.Second)

	// Display the blockchain
	fmt.Println("\n8. Displaying the blockchain:")
	fmt.Println("==========================================")
	bc.Print()

	// Validate the blockchain (normal case - no tampering)
	fmt.Println("\n9. Validating the blockchain (normal case - no tampering)...")
	fmt.Println("   Expected: IsValid() should return TRUE (blockchain is valid)")
	if bc.IsValid() {
		fmt.Println("   GOOD: Blockchain is valid! No issues detected.")
	} else {
		fmt.Println("   PROBLEM: Blockchain is invalid! Something is wrong.")
	}

	// Test insufficient balance
	fmt.Println("\n10. Testing insufficient balance scenario...")
	invalidTx := NewTransaction(aliceWallet.Address, bobWallet.Address, 1000.0)
	if err := bc.ValidateTransaction(invalidTx); err != nil {
		fmt.Printf("   Transaction rejected: %v\n", err)
	} else {
		fmt.Println("   ERROR: Transaction should have been rejected!")
	}

	// Test tampering detection - Scenario 1: Modify transaction without recalculating Merkle root
	fmt.Println("\n11. Testing tampering detection...")

	// Use block 2 (first transaction block after coinbase)
	// Block 0: Genesis, Block 1: Coinbase, Block 2: tx1, Block 3: tx2, Block 4: tx3
	tamperBlockIndex := 2
	if len(bc.Blocks) <= tamperBlockIndex || len(bc.Blocks[tamperBlockIndex].Transactions) == 0 {
		fmt.Println("   ERROR: Cannot find block for tampering test")
		return
	}

	// Save original transaction and hash for restoration
	originalTx := bc.Blocks[tamperBlockIndex].Transactions[0]
	originalMerkleRoot := bc.Blocks[tamperBlockIndex].MerkleRoot
	originalHash := bc.Blocks[tamperBlockIndex].Hash

	fmt.Println("\n   Scenario 1: Hacker modifies transaction WITHOUT recalculating Merkle root")
	fmt.Printf("   Modifying transaction in Block #%d...\n", tamperBlockIndex)
	tamperedTx := NewTransaction(originalTx.From, originalTx.To, 1000.0) // Change amount
	bc.Blocks[tamperBlockIndex].Transactions[0] = tamperedTx
	// Note: Merkle root and hash are NOT recalculated

	fmt.Println("\n   Validating blockchain after tampering...")
	fmt.Println("   Context: Testing if system can DETECT tampering")
	fmt.Println("   Expected: IsValid() should return FALSE (to indicate tampering was detected)")
	if bc.IsValid() {
		fmt.Println("   PROBLEM: Tampering was NOT detected! This is a security breach!")
		fmt.Println("   The blockchain thinks it's valid even though data was tampered.")
	} else {
		fmt.Println("   GOOD: Tampering was successfully detected!")
		fmt.Println("   The security system worked correctly - it detected the tampering.")
		fmt.Println("   Reason: Merkle root mismatch - stored root doesn't match calculated root")
	}

	// Restore the blockchain for scenario 2
	fmt.Println("\n   Restoring blockchain...")
	bc.Blocks[tamperBlockIndex].Transactions[0] = originalTx
	bc.Blocks[tamperBlockIndex].MerkleRoot = originalMerkleRoot
	bc.Blocks[tamperBlockIndex].Hash = originalHash

	// Test tampering detection - Scenario 2: Modify transaction AND recalculate Merkle root
	fmt.Println("\n   Scenario 2: Hacker modifies transaction AND recalculates Merkle root")
	fmt.Printf("   Modifying transaction in Block #%d...\n", tamperBlockIndex)
	bc.Blocks[tamperBlockIndex].Transactions[0] = tamperedTx
	fmt.Println("   Recalculating Merkle root and hash (without mining)...")
	merkleTree := NewMerkleTree(bc.Blocks[tamperBlockIndex].Transactions)
	bc.Blocks[tamperBlockIndex].MerkleRoot = merkleTree.GetRootHash()
	bc.Blocks[tamperBlockIndex].Hash = bc.Blocks[tamperBlockIndex].CalculateHash()

	fmt.Println("\n   Validating blockchain after tampering...")
	fmt.Println("   Context: Testing if system can DETECT tampering")
	fmt.Println("   Expected: IsValid() should return FALSE (to indicate tampering was detected)")
	if bc.IsValid() {
		fmt.Println("   PROBLEM: Tampering was NOT detected! This is a security breach!")
		fmt.Println("   The blockchain thinks it's valid even though data was tampered.")
	} else {
		fmt.Println("   GOOD: Tampering was successfully detected!")
		fmt.Println("   The security system worked correctly - it detected the tampering.")
		fmt.Println("   Reason: Proof of work invalid - hash doesn't meet difficulty requirement")
	}

	fmt.Println("\n=== Demo Complete ===")
}
