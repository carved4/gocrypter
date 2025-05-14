package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/chacha20"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Run with %s <inputfile.exe>\n", os.Args[0])
		os.Exit(1)
	}

	fname := os.Args[1]
	plaintextBytes, err := os.ReadFile(fname)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	stubDir := filepath.Join("..", "stub")
	encryptedFilePath := filepath.Join(stubDir, "encrypted_Input.bin")
	keyFilePath := filepath.Join(stubDir, "key.txt")

	encryptedFile, err := os.Create(encryptedFilePath)
	if err != nil {
		log.Fatalf("Failed to create encrypted file: %v", err)
	}
	defer encryptedFile.Close()

	keyFile, err := os.Create(keyFilePath)
	if err != nil {
		log.Fatalf("Failed to create key file: %v", err)
	}
	defer keyFile.Close()

	key := make([]byte, 32)
	nonce := make([]byte, 12)
	
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Fatalf("Failed to generate random key: %v", err)
	}
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("Failed to generate random nonce: %v", err)
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		log.Fatalf("Failed to create cipher: %v", err)
	}

	encryptedBytes := make([]byte, len(plaintextBytes))
	cipher.XORKeyStream(encryptedBytes, plaintextBytes)

	if _, err := encryptedFile.Write(encryptedBytes); err != nil {
		log.Fatalf("Failed to write encrypted data: %v", err)
	}

	keyData := append(key, nonce...)
	if _, err := keyFile.Write(keyData); err != nil {
		log.Fatalf("Failed to write key data: %v", err)
	}

	fmt.Printf("Encryption completed successfully!\nFiles saved to:\n- %s\n- %s\n", encryptedFilePath, keyFilePath)
} 