package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ContractType represents the type of smart contract
type ContractType string

const (
	ContractTypeSimple ContractType = "simple"
	ContractTypeToken  ContractType = "token"
	ContractTypeEscrow ContractType = "escrow"
	ContractTypeVoting ContractType = "voting"
)

// ContractContext holds execution context for contract calls
type ContractContext struct {
	Caller string
	Value  float64
	Args   []string
}

// SmartContract represents a smart contract deployed on the blockchain
type SmartContract struct {
	Address   string                 // Contract address (derived from deployer and nonce)
	Deployer  string                 // Address of the contract deployer
	Type      ContractType           // Type of contract
	Bytecode  string                 // Contract bytecode (simplified as string instructions)
	State     map[string]interface{} // Contract state storage
	CreatedAt int64                  // Block index when contract was created
	mu        sync.RWMutex           // Mutex for thread-safe state access
}

// ContractCall represents a call to a smart contract function
type ContractCall struct {
	ContractAddress string   // Address of the contract being called
	Function        string   // Function name to call
	Args            []string // Function arguments
	Value           float64  // Value sent with the call (for payable functions)
}

// NewSmartContract creates a new smart contract instance
func NewSmartContract(deployer string, contractType ContractType, bytecode string, blockIndex int64) *SmartContract {
	// Generate contract address from deployer address and block index
	addressData := fmt.Sprintf("%s%d", deployer, blockIndex)
	hash := sha256.Sum256([]byte(addressData))
	address := "0x" + hex.EncodeToString(hash[:])[:40]

	return &SmartContract{
		Address:   address,
		Deployer:  deployer,
		Type:      contractType,
		Bytecode:  bytecode,
		State:     make(map[string]interface{}),
		CreatedAt: blockIndex,
	}
}

// Execute executes a contract call and returns the result
func (sc *SmartContract) Execute(function string, args []string, caller string, value float64) (interface{}, error) {
	ctx := &ContractContext{
		Caller: caller,
		Value:  value,
		Args:   args,
	}

	switch sc.Type {
	case ContractTypeSimple:
		return sc.executeSimple(function, ctx)
	case ContractTypeToken:
		return sc.executeToken(function, ctx)
	case ContractTypeEscrow:
		return sc.executeEscrow(function, ctx)
	case ContractTypeVoting:
		return sc.executeVoting(function, ctx)
	default:
		return nil, fmt.Errorf("unknown contract type: %s", sc.Type)
	}
}

// Helper methods for safe state access

func (sc *SmartContract) getStateString(key string) (string, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	val, exists := sc.State[key]
	if !exists {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

func (sc *SmartContract) getStateFloat(key string) (float64, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	val, exists := sc.State[key]
	if !exists {
		return 0, false
	}
	f, ok := val.(float64)
	return f, ok
}

func (sc *SmartContract) getStateBool(key string) (bool, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	val, exists := sc.State[key]
	if !exists {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

func (sc *SmartContract) setState(key string, value interface{}) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.State[key] = value
}

func (sc *SmartContract) getBalances() map[string]float64 {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if _, exists := sc.State["balances"]; !exists {
		sc.State["balances"] = make(map[string]float64)
	}
	return sc.State["balances"].(map[string]float64)
}

func (sc *SmartContract) getProposals() map[string]int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if _, exists := sc.State["proposals"]; !exists {
		sc.State["proposals"] = make(map[string]int)
	}
	return sc.State["proposals"].(map[string]int)
}

func (sc *SmartContract) getVoters() map[string]bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if _, exists := sc.State["voters"]; !exists {
		sc.State["voters"] = make(map[string]bool)
	}
	return sc.State["voters"].(map[string]bool)
}

// Validation helpers

func validateArgsCount(args []string, required int, funcName string) error {
	if len(args) < required {
		return fmt.Errorf("%s requires %d argument(s)", funcName, required)
	}
	return nil
}

func parseAmount(amountStr string) (float64, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %s", amountStr)
	}
	if amount < 0 {
		return 0, fmt.Errorf("amount cannot be negative: %.2f", amount)
	}
	return amount, nil
}

func truncateAddress(addr string) string {
	if len(addr) > 16 {
		return addr[:16] + "..."
	}
	return addr
}

// executeSimple executes a simple contract (basic storage)
func (sc *SmartContract) executeSimple(function string, ctx *ContractContext) (interface{}, error) {
	switch function {
	case "set":
		if err := validateArgsCount(ctx.Args, 2, "set"); err != nil {
			return nil, err
		}
		key, val := ctx.Args[0], ctx.Args[1]
		sc.setState(key, val)
		return fmt.Sprintf("Set %s = %s", key, val), nil

	case "get":
		if err := validateArgsCount(ctx.Args, 1, "get"); err != nil {
			return nil, err
		}
		key := ctx.Args[0]
		sc.mu.RLock()
		val, exists := sc.State[key]
		sc.mu.RUnlock()
		if !exists {
			return nil, fmt.Errorf("key '%s' not found", key)
		}
		return val, nil

	case "delete":
		if err := validateArgsCount(ctx.Args, 1, "delete"); err != nil {
			return nil, err
		}
		key := ctx.Args[0]
		sc.mu.Lock()
		_, exists := sc.State[key]
		if exists {
			delete(sc.State, key)
		}
		sc.mu.Unlock()
		if !exists {
			return nil, fmt.Errorf("key '%s' not found", key)
		}
		return fmt.Sprintf("Deleted key '%s'", key), nil

	case "exists":
		if err := validateArgsCount(ctx.Args, 1, "exists"); err != nil {
			return nil, err
		}
		key := ctx.Args[0]
		sc.mu.RLock()
		_, exists := sc.State[key]
		sc.mu.RUnlock()
		return exists, nil

	default:
		return nil, fmt.Errorf("unknown function: %s", function)
	}
}

// executeToken executes a token contract (ERC-20 like)
func (sc *SmartContract) executeToken(function string, ctx *ContractContext) (interface{}, error) {
	balances := sc.getBalances()

	// Initialize total supply if not exists
	if _, exists := sc.getStateFloat("totalSupply"); !exists {
		sc.setState("totalSupply", 0.0)
	}

	switch function {
	case "transfer":
		if err := validateArgsCount(ctx.Args, 2, "transfer"); err != nil {
			return nil, err
		}
		to := ctx.Args[0]
		amount, err := parseAmount(ctx.Args[1])
		if err != nil {
			return nil, err
		}
		if amount == 0 {
			return nil, fmt.Errorf("transfer amount must be greater than zero")
		}

		sc.mu.Lock()
		callerBalance := balances[ctx.Caller]
		if callerBalance < amount {
			sc.mu.Unlock()
			return nil, fmt.Errorf("insufficient balance: %.2f < %.2f", callerBalance, amount)
		}
		balances[ctx.Caller] -= amount
		balances[to] += amount
		sc.mu.Unlock()

		return fmt.Sprintf("Transferred %.2f tokens from %s to %s",
			amount, truncateAddress(ctx.Caller), truncateAddress(to)), nil

	case "balanceOf":
		if err := validateArgsCount(ctx.Args, 1, "balanceOf"); err != nil {
			return nil, err
		}
		address := ctx.Args[0]
		sc.mu.RLock()
		balance := balances[address]
		sc.mu.RUnlock()
		return balance, nil

	case "totalSupply":
		supply, _ := sc.getStateFloat("totalSupply")
		return supply, nil

	case "mint":
		if ctx.Caller != sc.Deployer {
			return nil, fmt.Errorf("only deployer can mint tokens")
		}
		if err := validateArgsCount(ctx.Args, 2, "mint"); err != nil {
			return nil, err
		}
		to := ctx.Args[0]
		amount, err := parseAmount(ctx.Args[1])
		if err != nil {
			return nil, err
		}

		sc.mu.Lock()
		totalSupply, _ := sc.State["totalSupply"].(float64)
		totalSupply += amount
		sc.State["totalSupply"] = totalSupply
		balances[to] += amount
		sc.mu.Unlock()

		return fmt.Sprintf("Minted %.2f tokens to %s (Total supply: %.2f)",
			amount, truncateAddress(to), totalSupply), nil

	case "burn":
		if err := validateArgsCount(ctx.Args, 1, "burn"); err != nil {
			return nil, err
		}
		amount, err := parseAmount(ctx.Args[0])
		if err != nil {
			return nil, err
		}

		sc.mu.Lock()
		callerBalance := balances[ctx.Caller]
		if callerBalance < amount {
			sc.mu.Unlock()
			return nil, fmt.Errorf("insufficient balance to burn: %.2f < %.2f", callerBalance, amount)
		}
		balances[ctx.Caller] -= amount
		totalSupply, _ := sc.State["totalSupply"].(float64)
		totalSupply -= amount
		sc.State["totalSupply"] = totalSupply
		sc.mu.Unlock()

		return fmt.Sprintf("Burned %.2f tokens from %s (Total supply: %.2f)",
			amount, truncateAddress(ctx.Caller), totalSupply), nil

	default:
		return nil, fmt.Errorf("unknown function: %s", function)
	}
}

// executeEscrow executes an escrow contract
func (sc *SmartContract) executeEscrow(function string, ctx *ContractContext) (interface{}, error) {
	// Initialize escrow state
	sc.mu.Lock()
	if _, exists := sc.State["deposited"]; !exists {
		sc.State["deposited"] = 0.0
	}
	if _, exists := sc.State["released"]; !exists {
		sc.State["released"] = false
	}
	if _, exists := sc.State["beneficiary"]; !exists {
		if len(ctx.Args) > 0 {
			sc.State["beneficiary"] = ctx.Args[0]
		} else {
			sc.State["beneficiary"] = ctx.Caller
		}
	}
	if _, exists := sc.State["arbiter"]; !exists {
		sc.State["arbiter"] = sc.Deployer
	}
	sc.mu.Unlock()

	deposited, _ := sc.getStateFloat("deposited")
	released, _ := sc.getStateBool("released")
	beneficiary, _ := sc.getStateString("beneficiary")
	arbiter, _ := sc.getStateString("arbiter")

	switch function {
	case "deposit":
		if released {
			return nil, fmt.Errorf("escrow already released")
		}
		if ctx.Value <= 0 {
			return nil, fmt.Errorf("deposit value must be greater than zero")
		}
		newTotal := deposited + ctx.Value
		sc.setState("deposited", newTotal)
		return fmt.Sprintf("Deposited %.2f coins to escrow. Total: %.2f", ctx.Value, newTotal), nil

	case "release":
		if ctx.Caller != arbiter && ctx.Caller != sc.Deployer {
			return nil, fmt.Errorf("only arbiter or deployer can release escrow")
		}
		if released {
			return nil, fmt.Errorf("escrow already released")
		}
		if deposited == 0 {
			return nil, fmt.Errorf("no funds in escrow")
		}
		sc.setState("released", true)
		return fmt.Sprintf("Released %.2f coins to beneficiary %s",
			deposited, truncateAddress(beneficiary)), nil

	case "refund":
		if ctx.Caller != arbiter && ctx.Caller != sc.Deployer {
			return nil, fmt.Errorf("only arbiter or deployer can refund escrow")
		}
		if released {
			return nil, fmt.Errorf("escrow already released")
		}
		if deposited == 0 {
			return nil, fmt.Errorf("no funds in escrow")
		}
		sc.setState("released", true)
		sc.setState("refunded", true)
		return fmt.Sprintf("Refunded %.2f coins", deposited), nil

	case "getBalance":
		return deposited, nil

	case "getStatus":
		refunded, _ := sc.getStateBool("refunded")
		return map[string]interface{}{
			"deposited":   deposited,
			"released":    released,
			"refunded":    refunded,
			"beneficiary": beneficiary,
			"arbiter":     arbiter,
		}, nil

	default:
		return nil, fmt.Errorf("unknown function: %s", function)
	}
}

// executeVoting executes a voting contract
func (sc *SmartContract) executeVoting(function string, ctx *ContractContext) (interface{}, error) {
	proposals := sc.getProposals()
	voters := sc.getVoters()

	// Initialize voting state
	if _, exists := sc.getStateBool("votingEnded"); !exists {
		sc.setState("votingEnded", false)
	}

	votingEnded, _ := sc.getStateBool("votingEnded")

	switch function {
	case "propose":
		if votingEnded {
			return nil, fmt.Errorf("voting has ended")
		}
		if err := validateArgsCount(ctx.Args, 1, "propose"); err != nil {
			return nil, err
		}
		proposal := ctx.Args[0]
		sc.mu.Lock()
		if _, exists := proposals[proposal]; exists {
			sc.mu.Unlock()
			return nil, fmt.Errorf("proposal '%s' already exists", proposal)
		}
		proposals[proposal] = 0
		sc.mu.Unlock()
		return fmt.Sprintf("Proposal '%s' added", proposal), nil

	case "vote":
		if votingEnded {
			return nil, fmt.Errorf("voting has ended")
		}
		if err := validateArgsCount(ctx.Args, 1, "vote"); err != nil {
			return nil, err
		}
		proposal := ctx.Args[0]

		sc.mu.Lock()
		if voters[ctx.Caller] {
			sc.mu.Unlock()
			return nil, fmt.Errorf("address already voted")
		}
		if _, exists := proposals[proposal]; !exists {
			sc.mu.Unlock()
			return nil, fmt.Errorf("proposal '%s' not found", proposal)
		}
		proposals[proposal]++
		voters[ctx.Caller] = true
		sc.mu.Unlock()

		return fmt.Sprintf("Voted for '%s'", proposal), nil

	case "getResults":
		sc.mu.RLock()
		result := make(map[string]int)
		for k, v := range proposals {
			result[k] = v
		}
		sc.mu.RUnlock()
		return result, nil

	case "getWinner":
		sc.mu.RLock()
		var winner string
		maxVotes := -1
		for proposal, votes := range proposals {
			if votes > maxVotes {
				maxVotes = votes
				winner = proposal
			}
		}
		sc.mu.RUnlock()
		if winner == "" {
			return nil, fmt.Errorf("no proposals found")
		}
		return map[string]interface{}{
			"winner": winner,
			"votes":  maxVotes,
		}, nil

	case "endVoting":
		if ctx.Caller != sc.Deployer {
			return nil, fmt.Errorf("only deployer can end voting")
		}
		if votingEnded {
			return nil, fmt.Errorf("voting already ended")
		}
		sc.setState("votingEnded", true)
		return "Voting ended", nil

	default:
		return nil, fmt.Errorf("unknown function: %s", function)
	}
}

// GetAddress returns the contract address
func (sc *SmartContract) GetAddress() string {
	return sc.Address
}

// GetState returns the contract state as JSON
func (sc *SmartContract) GetState() (string, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	stateJSON, err := json.MarshalIndent(sc.State, "", "  ")
	if err != nil {
		return "", err
	}
	return string(stateJSON), nil
}

// GetDeployer returns the contract deployer address
func (sc *SmartContract) GetDeployer() string {
	return sc.Deployer
}

// GetType returns the contract type
func (sc *SmartContract) GetType() ContractType {
	return sc.Type
}

// ContractRegistry manages deployed smart contracts
type ContractRegistry struct {
	Contracts map[string]*SmartContract // Map of contract address to contract
	mu        sync.RWMutex              // Mutex for thread-safe access
}

// NewContractRegistry creates a new contract registry
func NewContractRegistry() *ContractRegistry {
	return &ContractRegistry{
		Contracts: make(map[string]*SmartContract),
	}
}

// DeployContract deploys a new smart contract
func (cr *ContractRegistry) DeployContract(deployer string, contractType ContractType, bytecode string, blockIndex int64) (*SmartContract, error) {
	if deployer == "" {
		return nil, fmt.Errorf("deployer address cannot be empty")
	}
	contract := NewSmartContract(deployer, contractType, bytecode, blockIndex)
	cr.mu.Lock()
	cr.Contracts[contract.Address] = contract
	cr.mu.Unlock()
	return contract, nil
}

// GetContract retrieves a contract by address
func (cr *ContractRegistry) GetContract(address string) (*SmartContract, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	contract, exists := cr.Contracts[address]
	if !exists {
		return nil, fmt.Errorf("contract not found: %s", address)
	}
	return contract, nil
}

// CallContract calls a function on a smart contract
func (cr *ContractRegistry) CallContract(contractAddress, function string, args []string, caller string, value float64) (interface{}, error) {
	contract, err := cr.GetContract(contractAddress)
	if err != nil {
		return nil, err
	}
	return contract.Execute(function, args, caller, value)
}

// GetAllContracts returns all deployed contracts
func (cr *ContractRegistry) GetAllContracts() []*SmartContract {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	contracts := make([]*SmartContract, 0, len(cr.Contracts))
	for _, contract := range cr.Contracts {
		contracts = append(contracts, contract)
	}
	return contracts
}

// GetContractsByDeployer returns all contracts deployed by a specific address
func (cr *ContractRegistry) GetContractsByDeployer(deployer string) []*SmartContract {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	contracts := make([]*SmartContract, 0)
	for _, contract := range cr.Contracts {
		if contract.Deployer == deployer {
			contracts = append(contracts, contract)
		}
	}
	return contracts
}

// IsContractAddress checks if an address is a contract address
func IsContractAddress(address string) bool {
	return len(address) == 42 && strings.HasPrefix(address, "0x")
}

// ParseContractCall parses contract call data from transaction data
func ParseContractCall(data string) (*ContractCall, error) {
	if data == "" {
		return nil, fmt.Errorf("empty contract call data")
	}

	// Simple format: "function:arg1,arg2,arg3"
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid contract call format: expected 'function:args'")
	}

	function := strings.TrimSpace(parts[0])
	if function == "" {
		return nil, fmt.Errorf("function name cannot be empty")
	}

	argsStr := parts[1]
	var args []string
	if argsStr != "" {
		args = strings.Split(argsStr, ",")
		// Trim whitespace from each argument
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}
	} else {
		args = []string{}
	}

	return &ContractCall{
		Function: function,
		Args:     args,
	}, nil
}
