package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// CalculateHash calculates SHA-256 hash of the given data
func CalculateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// BytesToHex converts byte slice to hexadecimal string
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}
