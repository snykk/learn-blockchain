package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// Web3Server represents a Web3 JSON-RPC server
type Web3Server struct {
	blockchain *Blockchain
	address    string
	port       int
	server     *http.Server
	mu         sync.RWMutex
	running    bool
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewWeb3Server creates a new Web3 server
func NewWeb3Server(blockchain *Blockchain, address string, port int) *Web3Server {
	return &Web3Server{
		blockchain: blockchain,
		address:    address,
		port:       port,
		running:    false,
	}
}

// Start starts the Web3 server
func (w *Web3Server) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return fmt.Errorf("server already running")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", w.handleRequest)

	addr := fmt.Sprintf("%s:%d", w.address, w.port)
	w.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	w.running = true
	fmt.Printf("Web3 server started on %s\n", addr)

	go func() {
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Web3 server error: %v\n", err)
		}
	}()

	return nil
}

// Stop stops the Web3 server
func (w *Web3Server) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return fmt.Errorf("server not running")
	}

	w.running = false
	return w.server.Close()
}

// handleRequest handles incoming JSON-RPC requests
func (w *Web3Server) handleRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.sendError(rw, -32700, "Parse error", 0)
		return
	}

	// Set response headers
	rw.Header().Set("Content-Type", "application/json")

	// Route to appropriate handler
	var result interface{}
	var err error

	switch req.Method {
	case "web3_clientVersion":
		result = w.clientVersion()
	case "eth_blockNumber":
		result = w.blockNumber()
	case "eth_getBalance":
		result, err = w.getBalance(req.Params)
	case "eth_getBlockByNumber":
		result, err = w.getBlockByNumber(req.Params)
	case "eth_getTransactionCount":
		result, err = w.getTransactionCount(req.Params)
	case "eth_sendTransaction":
		result, err = w.sendTransaction(req.Params)
	case "eth_call":
		result, err = w.call(req.Params)
	case "eth_getCode":
		result, err = w.getCode(req.Params)
	default:
		w.sendError(rw, -32601, "Method not found", req.ID)
		return
	}

	if err != nil {
		w.sendError(rw, -32000, err.Error(), req.ID)
		return
	}

	w.sendResponse(rw, result, req.ID)
}

// sendResponse sends a successful JSON-RPC response
func (w *Web3Server) sendResponse(rw http.ResponseWriter, result interface{}, id int) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
	json.NewEncoder(rw).Encode(resp)
}

// sendError sends an error JSON-RPC response
func (w *Web3Server) sendError(rw http.ResponseWriter, code int, message string, id int) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
		ID: id,
	}
	json.NewEncoder(rw).Encode(resp)
}

// clientVersion returns the client version
func (w *Web3Server) clientVersion() string {
	return "learn-blockchain/v1.0.0/go"
}

// blockNumber returns the latest block number
func (w *Web3Server) blockNumber() string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	blockNum := len(w.blockchain.Blocks) - 1
	return fmt.Sprintf("0x%x", blockNum)
}

// getBalance returns the balance of an address
func (w *Web3Server) getBalance(params []interface{}) (string, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	address, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	// Remove 0x prefix if present
	if len(address) > 2 && address[:2] == "0x" {
		address = address[2:]
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	balance := w.blockchain.GetBalance(address)

	// Convert to Wei (1 coin = 1e18 Wei for compatibility)
	weiBalance := int64(balance * 1e18)
	return fmt.Sprintf("0x%x", weiBalance), nil
}

// getBlockByNumber returns a block by number
func (w *Web3Server) getBlockByNumber(params []interface{}) (interface{}, error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("missing block number parameter")
	}

	blockNumStr, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid block number parameter")
	}

	var blockNum int
	if blockNumStr == "latest" {
		w.mu.RLock()
		blockNum = len(w.blockchain.Blocks) - 1
		w.mu.RUnlock()
	} else {
		// Parse hex number
		if len(blockNumStr) > 2 && blockNumStr[:2] == "0x" {
			blockNumStr = blockNumStr[2:]
		}
		num, err := strconv.ParseInt(blockNumStr, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid block number format")
		}
		blockNum = int(num)
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	if blockNum < 0 || blockNum >= len(w.blockchain.Blocks) {
		return nil, nil // Block not found, return null
	}

	block := w.blockchain.Blocks[blockNum]

	// Format block for Web3 response
	return map[string]interface{}{
		"number":           fmt.Sprintf("0x%x", block.Index),
		"hash":             "0x" + block.Hash,
		"parentHash":       "0x" + block.PreviousHash,
		"timestamp":        fmt.Sprintf("0x%x", block.Timestamp.Unix()),
		"transactions":     formatTransactions(block.Transactions),
		"transactionsRoot": "0x" + block.MerkleRoot,
	}, nil
}

// getTransactionCount returns the number of transactions sent from an address
func (w *Web3Server) getTransactionCount(params []interface{}) (string, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	address, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	// Remove 0x prefix if present
	if len(address) > 2 && address[:2] == "0x" {
		address = address[2:]
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	count := 0
	for _, block := range w.blockchain.Blocks {
		for _, tx := range block.Transactions {
			if tx.From == address {
				count++
			}
		}
	}

	return fmt.Sprintf("0x%x", count), nil
}

// sendTransaction sends a new transaction
func (w *Web3Server) sendTransaction(params []interface{}) (string, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("missing transaction parameter")
	}

	txData, ok := params[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid transaction parameter")
	}

	// Extract transaction fields
	from, _ := txData["from"].(string)
	to, _ := txData["to"].(string)
	valueStr, _ := txData["value"].(string)

	// Parse value (hex)
	if len(valueStr) > 2 && valueStr[:2] == "0x" {
		valueStr = valueStr[2:]
	}
	value, err := strconv.ParseInt(valueStr, 16, 64)
	if err != nil {
		return "", fmt.Errorf("invalid value format")
	}

	// Convert from Wei to coins (1e18 Wei = 1 coin)
	amount := float64(value) / 1e18

	// Create transaction
	tx := NewTransaction(from, to, amount)

	w.mu.Lock()
	err = w.blockchain.AddTransactionToMempool(tx)
	w.mu.Unlock()

	if err != nil {
		return "", err
	}

	// Return transaction hash
	txHash := hex.EncodeToString(tx.Hash())
	return "0x" + txHash, nil
}

// call executes a contract call (read-only)
func (w *Web3Server) call(params []interface{}) (string, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("missing call parameter")
	}

	callData, ok := params[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid call parameter")
	}

	to, _ := callData["to"].(string)
	_ = callData["data"] // Contract data (not used in simplified implementation)

	// Remove 0x prefix
	if len(to) > 2 && to[:2] == "0x" {
		to = to[2:]
	}

	// This is a simplified implementation
	// In a real implementation, you would decode the data and execute the contract
	return "0x", nil
}

// getCode returns the code at a given address (for contracts)
func (w *Web3Server) getCode(params []interface{}) (string, error) {
	if len(params) < 1 {
		return "", fmt.Errorf("missing address parameter")
	}

	address, ok := params[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid address parameter")
	}

	// Remove 0x prefix
	if len(address) > 2 && address[:2] == "0x" {
		address = address[2:]
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	// Check if address is a contract
	if IsContractAddress(address) {
		contract, err := w.blockchain.GetContract(address)
		if err != nil {
			return "0x", nil // No code
		}
		return "0x" + contract.Bytecode, nil
	}

	return "0x", nil // No code (regular address)
}

// formatTransactions formats transactions for Web3 response
func formatTransactions(transactions []*Transaction) []map[string]interface{} {
	result := make([]map[string]interface{}, len(transactions))
	for i, tx := range transactions {
		result[i] = map[string]interface{}{
			"from":  tx.From,
			"to":    tx.To,
			"value": fmt.Sprintf("0x%x", int64(tx.Amount*1e18)),
			"hash":  "0x" + hex.EncodeToString(tx.Hash()),
		}
	}
	return result
}
