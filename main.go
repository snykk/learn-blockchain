package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Enhanced Blockchain Implementation ===")
	fmt.Println("Features: Transactions, Merkle Tree, Wallet & Signing, Balance System")
	fmt.Println("          Mempool, Full Signature Verification, Proof of Stake")
	fmt.Println("          Delegated Proof of Stake, Transaction Fees, Block Rewards")
	fmt.Println("          Network/P2P, Smart Contracts")
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

	// Transaction 1: Alice sends 10 coins to Bob (with fee)
	tx1 := NewTransactionWithFee(aliceWallet.Address, bobWallet.Address, 10.0, 0.5)
	if err := bc.ValidateTransaction(tx1); err != nil {
		fmt.Printf("Error: Transaction 1 is invalid: %v\n", err)
		return
	}
	if err := aliceWallet.SignTransaction(tx1); err != nil {
		fmt.Printf("Error signing transaction 1: %v\n", err)
		return
	}
	fmt.Printf("   Transaction 1: %s\n", tx1.String())

	// Transaction 2: Bob sends 5 coins to Charlie (with fee)
	tx2 := NewTransactionWithFee(bobWallet.Address, charlieWallet.Address, 5.0, 0.3)
	if err := bc.ValidateTransaction(tx2); err != nil {
		fmt.Printf("Error: Transaction 2 is invalid: %v\n", err)
		return
	}
	if err := bobWallet.SignTransaction(tx2); err != nil {
		fmt.Printf("Error signing transaction 2: %v\n", err)
		return
	}
	fmt.Printf("   Transaction 2: %s\n", tx2.String())

	// Transaction 3: Charlie sends 3 coins to Alice (no fee)
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

	// Add blocks with transactions and miner rewards
	fmt.Println("\n6. Adding blocks to the blockchain with miner rewards...")

	// Create a miner wallet
	minerWallet, err := NewWallet()
	if err != nil {
		fmt.Printf("Error creating miner wallet: %v\n", err)
		return
	}
	fmt.Printf("   Miner wallet: %s\n", minerWallet.Address)

	if err := bc.AddBlockWithReward([]*Transaction{tx1}, minerWallet.Address); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	if err := bc.AddBlockWithReward([]*Transaction{tx2}, minerWallet.Address); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	if err := bc.AddBlockWithReward([]*Transaction{tx3}, minerWallet.Address); err != nil {
		fmt.Printf("Error adding block: %v\n", err)
		return
	}
	time.Sleep(1 * time.Second)

	// Display balances after transactions
	fmt.Println("\n7. Balances after transactions (including fees):")
	fmt.Printf("   Alice: %.2f coins\n", bc.GetBalance(aliceWallet.Address))
	fmt.Printf("   Bob: %.2f coins\n", bc.GetBalance(bobWallet.Address))
	fmt.Printf("   Charlie: %.2f coins\n", bc.GetBalance(charlieWallet.Address))
	fmt.Printf("   Miner: %.2f coins (from rewards)\n", bc.GetMinerRewards(minerWallet.Address))
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

	// Demo: Mempool functionality
	fmt.Println("\n12. Demonstrating Mempool functionality...")

	// Create new transactions and add to mempool
	fmt.Println("\n   Creating new transactions and adding to mempool...")
	tx4 := NewTransaction(aliceWallet.Address, bobWallet.Address, 5.0)
	if err := aliceWallet.SignTransaction(tx4); err != nil {
		fmt.Printf("Error signing transaction 4: %v\n", err)
		return
	}
	if err := bc.AddTransactionToMempool(tx4); err != nil {
		fmt.Printf("Error adding to mempool: %v\n", err)
	} else {
		fmt.Printf("   Transaction 4 added to mempool: %s\n", tx4.String())
	}

	tx5 := NewTransaction(bobWallet.Address, charlieWallet.Address, 3.0)
	if err := bobWallet.SignTransaction(tx5); err != nil {
		fmt.Printf("Error signing transaction 5: %v\n", err)
		return
	}
	if err := bc.AddTransactionToMempool(tx5); err != nil {
		fmt.Printf("Error adding to mempool: %v\n", err)
	} else {
		fmt.Printf("   Transaction 5 added to mempool: %s\n", tx5.String())
	}

	fmt.Printf("\n   Mempool size: %d transactions\n", bc.Mempool.Size())

	// Create block from mempool
	fmt.Println("\n   Creating block from mempool transactions...")
	if err := bc.AddBlockFromMempool(10); err != nil {
		fmt.Printf("Error creating block from mempool: %v\n", err)
	} else {
		fmt.Printf("   Mempool size after block creation: %d transactions\n", bc.Mempool.Size())
	}

	// Demo: Full Signature Verification
	fmt.Println("\n13. Demonstrating Full Signature Verification...")
	fmt.Println("   All transactions in blockchain are now verified using stored public keys")
	fmt.Println("   Signature verification is performed automatically during validation")

	// Demo: Proof of Stake
	fmt.Println("\n14. Demonstrating Proof of Stake consensus...")
	stakeholders := bc.CalculateStakeFromBlockchain()
	fmt.Println("   Current stakeholders and their stakes:")
	for address, stake := range stakeholders {
		fmt.Printf("   - %s: %.2f coins\n", address[:16]+"...", stake)
	}

	// Select validator
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	pos := NewProofOfStake(lastBlock, stakeholders)
	validator := pos.SelectValidator()
	if validator != "" {
		fmt.Printf("\n   Selected validator: %s\n", validator[:16]+"...")
		fmt.Println("   (In PoS, validator is selected based on stake weight)")
		fmt.Println("   Note: PoS block creation requires validator to be selected")
	}

	// Demo: Transaction Fees and Block Rewards Summary
	fmt.Println("\n15. Transaction Fees and Block Rewards Summary...")
	fmt.Println("   Transaction fees are deducted from sender's balance")
	fmt.Println("   Block rewards are given to miners for creating blocks")
	if len(bc.Blocks) > 1 {
		// Find miner from block rewards
		for _, block := range bc.Blocks {
			for _, tx := range block.Transactions {
				if tx.From == "" && tx.To != "Genesis" {
					rewards := bc.GetMinerRewards(tx.To)
					if rewards > 0 {
						fmt.Printf("   Miner %s total rewards: %.2f coins\n", tx.To[:16]+"...", rewards)
						fmt.Println("   (Block rewards + transaction fees)")
						break
					}
				}
			}
		}
	}

	// Demo: Delegated Proof of Stake
	fmt.Println("\n16. Demonstrating Delegated Proof of Stake (DPoS)...")
	topDelegates := bc.GetTopDelegates(5)
	fmt.Println("   Top 5 delegates by votes:")
	for i, delegate := range topDelegates {
		fmt.Printf("   %d. %s - Votes: %.2f, Stake: %.2f\n",
			i+1, delegate.Address[:16]+"...", delegate.Votes, delegate.Stake)
	}

	// Select validator using DPoS
	if len(topDelegates) > 0 {
		lastBlock := bc.Blocks[len(bc.Blocks)-1]
		stakeholders := bc.CalculateStakeFromBlockchain()
		dpos := NewDelegatedProofOfStake(lastBlock, stakeholders)

		// Initialize votes from stakes
		for address, stake := range stakeholders {
			if stake > 0 {
				dpos.Vote(address, address, stake)
			}
		}

		validator := dpos.SelectValidator()
		if validator != "" {
			fmt.Printf("\n   Selected validator (round-robin): %s\n", validator[:16]+"...")
			fmt.Println("   (In DPoS, validators are selected in round-robin from top delegates)")
		}
	}

	// Demo: Network/P2P
	fmt.Println("\n17. Demonstrating Network/P2P functionality...")

	// Create nodes
	fmt.Println("\n   Creating network nodes...")
	node1 := NewNode("localhost", 3001)
	node2 := NewNode("localhost", 3002)
	node3 := NewNode("localhost", 3003)

	fmt.Printf("   Node 1: %s\n", node1.GetAddress())
	fmt.Printf("   Node 2: %s\n", node2.GetAddress())
	fmt.Printf("   Node 3: %s\n", node3.GetAddress())

	// Add peers
	fmt.Println("\n   Setting up peer connections...")
	node1.AddPeer(node2.GetAddress())
	node1.AddPeer(node3.GetAddress())
	node2.AddPeer(node1.GetAddress())
	node2.AddPeer(node3.GetAddress())
	node3.AddPeer(node1.GetAddress())
	node3.AddPeer(node2.GetAddress())

	fmt.Printf("   Node 1 peers: %d\n", len(node1.Peers))
	fmt.Printf("   Node 2 peers: %d\n", len(node2.Peers))
	fmt.Printf("   Node 3 peers: %d\n", len(node3.Peers))

	fmt.Println("\n   Note: In a real P2P network, nodes would:")
	fmt.Println("   - Start servers to accept connections")
	fmt.Println("   - Broadcast new blocks and transactions")
	fmt.Println("   - Synchronize blockchain with peers")
	fmt.Println("   - Handle network consensus")
	fmt.Println("   (Full network demo requires running multiple processes)")

	// Demo: Smart Contracts
	fmt.Println("\n18. Demonstrating Smart Contracts...")

	// Deploy a simple storage contract
	fmt.Println("\n   Deploying Simple Storage Contract...")
	simpleContract, err := bc.DeployContract(aliceWallet.Address, ContractTypeSimple, "simple_storage")
	if err != nil {
		fmt.Printf("Error deploying contract: %v\n", err)
	} else {
		fmt.Printf("   Contract deployed at: %s\n", simpleContract.GetAddress())
		fmt.Printf("   Deployer: %s\n", aliceWallet.Address[:16]+"...")

		// Call set function
		fmt.Println("\n   Calling set(key='name', value='Alice')...")
		tx1 := NewContractCallTransaction(aliceWallet.Address, simpleContract.GetAddress(), "set", []string{"name", "Alice"}, 0, 0.1)
		if err := aliceWallet.SignTransaction(tx1); err == nil {
			if err := bc.AddBlockWithReward([]*Transaction{tx1}, minerWallet.Address); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		}

		// Call get function
		fmt.Println("\n   Calling get(key='name')...")
		result, err := bc.CallContract(simpleContract.GetAddress(), "get", []string{"name"}, aliceWallet.Address, 0)
		if err != nil {
			fmt.Printf("Error calling contract: %v\n", err)
		} else {
			fmt.Printf("   Result: %v\n", result)
		}
	}

	// Deploy a token contract
	fmt.Println("\n   Deploying Token Contract...")
	tokenContract, err := bc.DeployContract(bobWallet.Address, ContractTypeToken, "token_contract")
	if err != nil {
		fmt.Printf("Error deploying contract: %v\n", err)
	} else {
		fmt.Printf("   Contract deployed at: %s\n", tokenContract.GetAddress())

		// Mint tokens
		fmt.Println("\n   Minting 100 tokens to Bob...")
		tx2 := NewContractCallTransaction(bobWallet.Address, tokenContract.GetAddress(), "mint", []string{bobWallet.Address, "100"}, 0, 0.1)
		if err := bobWallet.SignTransaction(tx2); err == nil {
			if err := bc.AddBlockWithReward([]*Transaction{tx2}, minerWallet.Address); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		}

		// Transfer tokens
		fmt.Println("\n   Transferring 20 tokens from Bob to Charlie...")
		tx3 := NewContractCallTransaction(bobWallet.Address, tokenContract.GetAddress(), "transfer", []string{charlieWallet.Address, "20"}, 0, 0.1)
		if err := bobWallet.SignTransaction(tx3); err == nil {
			if err := bc.AddBlockWithReward([]*Transaction{tx3}, minerWallet.Address); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		}

		// Check balance
		fmt.Println("\n   Checking Charlie's token balance...")
		balance, err := bc.CallContract(tokenContract.GetAddress(), "balanceOf", []string{charlieWallet.Address}, charlieWallet.Address, 0)
		if err != nil {
			fmt.Printf("Error calling contract: %v\n", err)
		} else {
			fmt.Printf("   Charlie's balance: %.2f tokens\n", balance)
		}
	}

	// Deploy a voting contract
	fmt.Println("\n   Deploying Voting Contract...")
	votingContract, err := bc.DeployContract(charlieWallet.Address, ContractTypeVoting, "voting_contract")
	if err != nil {
		fmt.Printf("Error deploying contract: %v\n", err)
	} else {
		fmt.Printf("   Contract deployed at: %s\n", votingContract.GetAddress())

		// Add proposals
		fmt.Println("\n   Adding proposals...")
		tx4 := NewContractCallTransaction(charlieWallet.Address, votingContract.GetAddress(), "propose", []string{"Option A"}, 0, 0.1)
		if err := charlieWallet.SignTransaction(tx4); err == nil {
			if err := bc.AddBlockWithReward([]*Transaction{tx4}, minerWallet.Address); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		}

		tx5 := NewContractCallTransaction(charlieWallet.Address, votingContract.GetAddress(), "propose", []string{"Option B"}, 0, 0.1)
		if err := charlieWallet.SignTransaction(tx5); err == nil {
			if err := bc.AddBlockWithReward([]*Transaction{tx5}, minerWallet.Address); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		}

		// Vote
		fmt.Println("\n   Alice voting for Option A...")
		tx6 := NewContractCallTransaction(aliceWallet.Address, votingContract.GetAddress(), "vote", []string{"Option A"}, 0, 0.1)
		if err := aliceWallet.SignTransaction(tx6); err == nil {
			if err := bc.AddBlockWithReward([]*Transaction{tx6}, minerWallet.Address); err != nil {
				fmt.Printf("Error adding block: %v\n", err)
			}
		}

		// Get results
		fmt.Println("\n   Getting voting results...")
		results, err := bc.CallContract(votingContract.GetAddress(), "getResults", []string{}, charlieWallet.Address, 0)
		if err != nil {
			fmt.Printf("Error calling contract: %v\n", err)
		} else {
			fmt.Printf("   Voting results: %v\n", results)
		}
	}

	fmt.Println("\n=== Demo Complete ===")
}
