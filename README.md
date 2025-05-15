# gocrypter - an implementation of https://github.com/Amaop/Rust-Crypter in golang with ChaCha20 instead of AES

(completely cross compatible, rust implementation only accepts x86_64 windows executables)

This project contains two components:

1. **crypt** - A tool to encrypt an executable file using ChaCha20, and turn it into two parts
    - a - the encrypted_input.bin (which are the encrypted bytes of the target executable)
    - b - the key.txt which is used in the stub build to decrypt the bytes and execute in memory
2. **stub** - a golang binary that embeds, decrypts, and executes the encrypted file directly in memory using go-memexec 

## How to Use
1. Clone the repository 
2. cd gocrypter
3. Put your target executable in gocrypter/crypt
4. In a shell, cd into gocrypter/crypt
5. To use crypter: `go run main.go example.exe` (whatever the name of the binary you put in crypt, obviously)
6. Then, the files encrypted_Input.bin and key.txt will be created in the gocrypter/stub directory
7. cd into gocrypter/stub
8. `go build -ldflags="-s -w" -trimpath -o stub.exe` 
(or if you're building an ELF binary)
 `GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -trimpath -o stub`
9. Done!

## Features

- Encryption/Decryption: ChaCha20 stream cipher
- Memory Execution: [github.com/amenzhinsky/go-memexec](https://github.com/amenzhinsky/go-memexec/tree/main)
- The encrypted file and key are embedded in the stub binary using Go's embed package

## Security Considerations

This is a basic implementation and may not be suitable for production use without additional OPSEC measures.
(obviously do not use this outside of a lab environment or an authorized redteam engagement, I am not responsible for your actions)

- The key is stored in plaintext within the stub binary
- No anti-debugging or VM detection mechanisms are implemented

## Example Results

### Pre stub REvil windows ransomware sample
![revilmain](https://github.com/user-attachments/assets/f78aaf2c-e51f-40f1-b5a4-6d107523e995)
[VT results] (https://www.virustotal.com/gui/file/2ca64feaaf5ab6cf96677fbc2bc0e1995b3bc93472d7af884139aa757240e3f6)

### Post stub REvil windows ransomware sample
![revilstub2](https://github.com/user-attachments/assets/52746c08-1e3f-45eb-876d-53c6b942d592)
[VT results] (https://www.virustotal.com/gui/file/4d07aa5f4ad2fdaf11319b38afda9d0cf6b2db358de8a9d5440f962221dbab25?nocache=1)
