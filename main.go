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
	fmt.Println("          Network/P2P, Smart Contracts, Web3 Integration, PBFT Consensus")
	fmt.Println("          Raft Consensus, Layer 2 Payment Channels")
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

	// Demo: Web3 Integration
	fmt.Println("\n19. Demonstrating Web3 Integration...")

	fmt.Println("\n   Starting Web3 JSON-RPC server...")
	web3Server := NewWeb3Server(bc, "localhost", 8545)
	if err := web3Server.Start(); err != nil {
		fmt.Printf("Error starting Web3 server: %v\n", err)
	} else {
		fmt.Println("   Web3 server started on http://localhost:8545")
		fmt.Println("\n   Available Web3 API endpoints:")
		fmt.Println("   - web3_clientVersion - Get client version")
		fmt.Println("   - eth_blockNumber - Get latest block number")
		fmt.Println("   - eth_getBalance - Get account balance")
		fmt.Println("   - eth_getBlockByNumber - Get block by number")
		fmt.Println("   - eth_getTransactionCount - Get transaction count")
		fmt.Println("   - eth_sendTransaction - Send new transaction")
		fmt.Println("   - eth_call - Execute contract call (read-only)")
		fmt.Println("   - eth_getCode - Get contract code")

		fmt.Println("\n   Example curl commands:")
		fmt.Println("   curl -X POST http://localhost:8545 \\")
		fmt.Println("     -H \"Content-Type: application/json\" \\")
		fmt.Println("     -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\":1}'")

		fmt.Println("\n   Note: Web3 server is running in background")
		fmt.Println("   You can test the endpoints using curl or Postman")

		// Give server time to start
		time.Sleep(1 * time.Second)
	}

	// Demo: PBFT Consensus
	fmt.Println("\n20. Demonstrating PBFT (Practical Byzantine Fault Tolerance) Consensus...")

	fmt.Println("\n   Setting up PBFT network...")
	// Create list of nodes (addresses)
	pbftNodes := []string{
		aliceWallet.Address,
		bobWallet.Address,
		charlieWallet.Address,
		minerWallet.Address,
	}

	fmt.Printf("   Total nodes: %d\n", len(pbftNodes))
	fmt.Printf("   Byzantine fault tolerance: Can tolerate %d faulty nodes\n", (len(pbftNodes)-1)/3)
	fmt.Printf("   Required votes (quorum): %d (2f+1)\n", 2*((len(pbftNodes)-1)/3)+1)

	// Create transactions for PBFT block
	fmt.Println("\n   Creating transactions for PBFT block...")
	pbftTx1 := NewTransaction(aliceWallet.Address, bobWallet.Address, 2.0)
	if err := aliceWallet.SignTransaction(pbftTx1); err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
	} else {
		fmt.Printf("   Transaction 1: Alice -> Bob (2.0 coins)\n")

		// Create block using PBFT consensus
		fmt.Println("\n   Creating block using PBFT consensus...")
		if err := bc.CreateBlockWithPBFT([]*Transaction{pbftTx1}, pbftNodes, aliceWallet.Address); err != nil {
			fmt.Printf("Error creating PBFT block: %v\n", err)
		}
	}

	fmt.Println("\n   PBFT Consensus Features:")
	fmt.Println("   Three-phase protocol (Pre-Prepare, Prepare, Commit)")
	fmt.Println("   Byzantine fault tolerance (tolerates malicious nodes)")
	fmt.Println("   Deterministic finality (no forks)")
	fmt.Println("   Quorum-based consensus (2f+1 votes required)")

	fmt.Println("\n=== Demo Complete ===")

	// Demo: Raft Consensus
	fmt.Println("\n21. Demonstrating Raft Consensus...")

	fmt.Println("\n   Setting up Raft network...")
	// Create list of nodes (addresses) for Raft cluster
	raftNodes := []string{
		aliceWallet.Address,
		bobWallet.Address,
		charlieWallet.Address,
		minerWallet.Address,
	}

	fmt.Printf("   Total nodes: %d\n", len(raftNodes))
	fmt.Printf("   Majority required: %d nodes\n", len(raftNodes)/2+1)
	fmt.Printf("   Fault tolerance: %d nodes can fail\n", (len(raftNodes)-1)/2)

	// Create transactions for Raft block
	fmt.Println("\n   Creating transactions for Raft block...")
	raftTx1 := NewTransaction(bobWallet.Address, charlieWallet.Address, 1.5)
	if err := bobWallet.SignTransaction(raftTx1); err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
	} else {
		fmt.Printf("   Transaction: Bob -> Charlie (1.5 coins)\n")

		raftTx2 := NewTransaction(charlieWallet.Address, aliceWallet.Address, 1.0)
		if err := charlieWallet.SignTransaction(raftTx2); err != nil {
			fmt.Printf("Error signing transaction: %v\n", err)
		} else {
			fmt.Printf("   Transaction: Charlie -> Alice (1.0 coins)\n")

			// Create block using Raft consensus
			fmt.Println("\n   Creating block using Raft consensus...")
			if err := bc.CreateBlockWithRaft([]*Transaction{raftTx1, raftTx2}, aliceWallet.Address, raftNodes); err != nil {
				fmt.Printf("Error creating Raft block: %v\n", err)
			}
		}
	}

	fmt.Println("\n   Raft Consensus Features:")
	fmt.Println("   Leader election (democratic leader selection)")
	fmt.Println("   Log replication (strong consistency)")
	fmt.Println("   Safety guarantees (at most one leader per term)")
	fmt.Println("   Heartbeat mechanism (leader maintains authority)")
	fmt.Println("   Fault tolerance (tolerates (N-1)/2 node failures)")
	fmt.Println("\n   Comparison with PBFT:")
	fmt.Println("   - Raft: Simpler to understand, leader-based")
	fmt.Println("   - PBFT: More complex, no fixed leader, BFT-resistant")
	fmt.Println("   - Raft: Crash fault tolerance (CFT)")
	fmt.Println("   - PBFT: Byzantine fault tolerance (BFT)")

	fmt.Println("\n=== Demo Complete ===")

	// Demo: Layer 2 Payment Channels
	fmt.Println("\n22. Demonstrating Layer 2 Payment Channels (State Channels)...")

	fmt.Println("\n   Setting up payment channel...")
	// Alice and Bob want to create a payment channel
	channel, err := bc.ChannelManager.CreateChannel(
		aliceWallet.Address,
		bobWallet.Address,
		20.0, // Alice deposits 20 coins
		10.0, // Bob deposits 10 coins
		24*time.Hour, // 24 hour timeout
	)
	if err != nil {
		fmt.Printf("Error creating channel: %v\n", err)
	} else {
		fmt.Printf("\n   ✓ Payment channel created successfully\n")
		fmt.Printf("   Channel ID: %s\n", channel.State.ChannelID[:16]+"...")
		time.Sleep(1 * time.Second)

		// Perform micropayments through the channel
		fmt.Println("\n   Processing micropayments through channel...")

		// Micropayment 1: Alice pays Bob 0.5 coins
		fmt.Println("\n   Transaction 1: Alice → Bob (0.5 coins)")
		newState1, err := channel.MicroPayment(aliceWallet.Address, 0.5)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			// Sign the new state
			sig1, _ := channel.SignState(newState1, aliceWallet.Address)
			sig2, _ := channel.SignState(newState1, bobWallet.Address)

			// Commit the signed state
			signedState := &ChannelSignature{
				State:      newState1,
				Signature1: sig1,
				Signature2: sig2,
			}
			channel.CommitState(signedState)
		}
		time.Sleep(500 * time.Millisecond)

		// Micropayment 2: Bob pays Alice 0.3 coins
		fmt.Println("\n   Transaction 2: Bob → Alice (0.3 coins)")
		newState2, err := channel.MicroPayment(bobWallet.Address, 0.3)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			sig1, _ := channel.SignState(newState2, aliceWallet.Address)
			sig2, _ := channel.SignState(newState2, bobWallet.Address)

			signedState := &ChannelSignature{
				State:      newState2,
				Signature1: sig1,
				Signature2: sig2,
			}
			channel.CommitState(signedState)
		}
		time.Sleep(500 * time.Millisecond)

		// Micropayment 3: Alice pays Bob 1.2 coins
		fmt.Println("\n   Transaction 3: Alice → Bob (1.2 coins)")
		newState3, err := channel.MicroPayment(aliceWallet.Address, 1.2)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			sig1, _ := channel.SignState(newState3, aliceWallet.Address)
			sig2, _ := channel.SignState(newState3, bobWallet.Address)

			signedState := &ChannelSignature{
				State:      newState3,
				Signature1: sig1,
				Signature2: sig2,
			}
			channel.CommitState(signedState)
		}
		time.Sleep(500 * time.Millisecond)

		// Perform more rapid micropayments (simulating high-frequency transactions)
		fmt.Println("\n   Processing rapid micropayments...")
		for i := 0; i < 5; i++ {
			amount := 0.1
			var sender string
			if i%2 == 0 {
				sender = aliceWallet.Address
			} else {
				sender = bobWallet.Address
			}

			newState, err := channel.MicroPayment(sender, amount)
			if err == nil {
				sig1, _ := channel.SignState(newState, aliceWallet.Address)
				sig2, _ := channel.SignState(newState, bobWallet.Address)

				signedState := &ChannelSignature{
					State:      newState,
					Signature1: sig1,
					Signature2: sig2,
				}
				channel.CommitState(signedState)
			}
		}
		time.Sleep(500 * time.Millisecond)

		// Display channel status
		fmt.Println("\n   Final Channel Status:")
		fmt.Printf("   %s\n", channel.GetStatus())

		// Close the channel
		fmt.Println("\n   Closing payment channel...")
		finalSig := &ChannelSignature{
			State:      channel.State,
			Signature1: "final_signature_1",
			Signature2: "final_signature_2",
		}
		if err := channel.CloseChannel(finalSig); err != nil {
			fmt.Printf("Error closing channel: %v\n", err)
		}

		// Display channel statistics
		fmt.Println("\n   Channel Manager Statistics:")
		stats := bc.ChannelManager.GetChannelStatistics()
		fmt.Printf("   Total Channels: %d\n", stats["total_channels"])
		fmt.Printf("   Open Channels: %d\n", stats["open_channels"])
		fmt.Printf("   Closed Channels: %d\n", stats["closed_channels"])
		fmt.Printf("   Total Off-Chain Transactions: %d\n", stats["total_transactions"])
		fmt.Printf("   Total Volume: %.2f coins\n", stats["total_volume"])

		fmt.Println("\n   Layer 2 Payment Channel Benefits:")
		fmt.Println("   ✓ Instant transactions (no block confirmation wait)")
		fmt.Println("   ✓ Near-zero transaction fees")
		fmt.Println("   ✓ High throughput (100+ TPS in channel)")
		fmt.Println("   ✓ Privacy (transactions not on blockchain)")
		fmt.Println("   ✓ Final security (settlement on blockchain)")

		fmt.Println("\n   Comparison with On-Chain Transactions:")
		fmt.Println("   - On-Chain: ~10s/block, requires mining, fees")
		fmt.Println("   - Payment Channel: Instant, no mining, minimal fees")
		fmt.Println("   - 8 transactions processed off-chain instantly")
		fmt.Println("   - Only 2 on-chain transactions: open + close")
		fmt.Println("   - Gas savings: ~75% (8 vs ~32 on-chain transactions)")
	}

	fmt.Println("\n=== Demo Complete ===")
}
