# gocrypter - an implementation of https://github.com/Amaop/Rust-Crypter in golang + cpp with ChaCha20 instead of AES encryption + a custom PE loader to execute in memory 
## Important Notice

This project contains two components:

1. **crypt** - A tool to encrypt a preferably golang executable file using ChaCha20, and turn it into two parts
    - a - the encrypted_input.bin (which are the encrypted bytes of the target executable)
    - b - the key.txt which is used in the stub build to decrypt the bytes and execute in memory
2. **stub** - a golang binary that embeds, decrypts, and executes the encrypted file directly in memory using a custom implementation of PE exec with process hollowing and a custom PE loader, which is also embedded in the stub binary to perform the
ACTUAL in memory execution of the encrypted file. (no go memexec bullshit dropping to /temp)


## How to Use
1. Clone the repository 
2. cd gocrypter
3. Put your target executable in gocrypter/crypt
4. In a shell, cd into gocrypter/crypt
5. To use crypter: `go run main.go example.exe` (whatever the name of the binary you put in crypt, obviously)
6. Then, the files encrypted_Input.bin and key.txt will be created in the gocrypter/stub directory
7. cd into gocrypter/stub
8. `bash build-all.sh` 
9. Done!

## Features

- Encryption/Decryption: ChaCha20 stream cipher
- Memory Execution: stub/cpp and internal/runpe/runpe.go (relocation, process hollowing, will compile for any EXE im pretty sure, but if not let me know)
- The encrypted file and key are embedded in the stub binary using Go's embed package

## Security Considerations

This is a basic implementation and may not be suitable for production use without additional OPSEC measures.
(obviously do not use this outside of a lab environment or an authorized redteam engagement, I am not responsible for your actions)

- The key is stored in plaintext within the stub binary
(implementing secure key derivation soon)
- No anti-debugging or VM detection mechanisms are implemented
(have a package for this but may not share until further testing)
## Example Usage with Go Binary that writes "it worked" to current dir


https://github.com/user-attachments/assets/b0a5e1dd-6e37-41ac-8564-c578c754fdba

