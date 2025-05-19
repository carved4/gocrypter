# gocrypter - an implementation of https://github.com/Amaop/Rust-Crypter in golang + cpp with ChaCha20-Poly1305 authenticated encryption + a custom PE loader to execute in memory 

## Notice
feel free to open issues or DM with feedback or use cases!

This project contains two components:

1. **crypt** - A tool to encrypt a preferably golang executable file using ChaCha20-Poly1305 authenticated encryption, and turn it into two parts
    - a - the encrypted_input.bin (which are the encrypted bytes of the target executable)
    - b - the config.txt which contains information needed for secure key derivation and decryption
2. **stub** - a golang binary that embeds, decrypts, and executes the encrypted file directly in memory using a custom implementation of PE exec with self process hollowing and a custom PE loader, which is also embedded in the stub binary to perform the
ACTUAL in memory execution of the encrypted file. (no go memexec bullshit dropping to /temp)


## How to Use
1. Clone the repository 
2. cd gocrypter
3. Put your target executable in gocrypter/crypt
4. In a shell, cd into gocrypter/crypt
5. To use crypter: `go run main.go example.exe` (whatever the name of the binary you put in crypt, obviously)
6. Then, the files encrypted_Input.bin and config.txt will be created in the gocrypter/stub directory
7. cd into gocrypter/stub
8. `bash build-all.sh` 
9. Done!

## Features

- Encryption/Decryption: ChaCha20-Poly1305 AEAD (Authenticated Encryption with Associated Data)
- Secure Key Derivation: Argon2id, a memory-hard KDF resistant to brute-force attacks
- Protection Against Bit-Flip Attacks: Authentication prevents tampering with encrypted data
- Memory Execution: stub/cpp and internal/runpe/runpe.go (relocation, self process hollowing, will compile for any EXE im pretty sure, but if not let me know)
- The encrypted file and configuration data are embedded in the stub binary using Go's embed package
- Sensitive files are automatically cleaned up after build is complete

## Security Improvements

The project now uses several security best practices:

1. **Authenticated Encryption**: ChaCha20-Poly1305 provides both confidentiality and integrity/authenticity
2. **Secure Key Derivation**: Instead of storing encryption keys directly, we use Argon2id to derive keys
3. **Memory-Hard Key Derivation**: Makes brute force attacks significantly more expensive
4. **Clean-up Process**: Sensitive files are removed after successful build

## Security Considerations

While significantly improved, this is still a basic implementation and may not be suitable for production use without additional OPSEC measures.
(obviously do not use this outside of a lab environment or an authorized redteam engagement, I am not responsible for your actions)

- No anti-debugging or VM detection mechanisms are implemented
(have a package for this but may not share until further testing)
## Example Usage with Go Binary that writes "it worked" to current dir





https://github.com/user-attachments/assets/029eee04-0d72-4a69-8b37-940f96f9aac5

