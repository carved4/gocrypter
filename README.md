# gocrypter - an implementation of https://github.com/Amaop/Rust-Crypter in golang with ChaCha20 instead of AES

This project contains two components:

1. **crypt** - A tool to encrypt an executable file using ChaCha20, and turn it into two parts
    a - the encrypted_input.bin (which are the encrypted bytes of the target executable)
    b - the key.txt which is used in the stub build to decrypt the bytes and execute in memory
2. **stub** - a golang binary that embeds, decrypts, and executes the encrypted file directly in memory using go-memexec 

## how to Use
1. clone the repository 
2. cd gocrypter
3. put your target executable in gocrypter/crypt
4. in a shell, cd into gocrypter/crypt
5. to use crypter: go run main.go example.exe (whatever the name of the binary you put in crypt, obviously)
6. then, the files encrypted_Input.bin and key.txt will be created in the gocrypter/stub directory
7. cd into gocrypter/stub
8. go build -ldflags="-s -w" -trimpath -o stub.exe
10. done!



- Encryption/Decryption: ChaCha20 stream cipher
- Memory Execution: [github.com/wrwrabbit/go-memexec](https://pkg.go.dev/github.com/wrwrabbit/go-memexec)
- The encrypted file and key are embedded in the stub binary using Go's embed package

## Security Considerations

This is a basic implementation and may not be suitable for production use without additional OPSEC measures.
(obviously do not use this outside of a lab environment or an authorized redteam engagement, I am not responsible for your actions)

- The key is stored in plaintext within the stub binary
- No anti-debugging or VM detection mechanisms are implemented


