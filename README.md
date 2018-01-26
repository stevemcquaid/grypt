# grypt - Encryption Utility in Golang

# Goals
  * Protect secrets stored in public by providing easy-to-use strong encryption
    * Used in place of, or in combination with, [Veracrypt](https://en.wikipedia.org/wiki/VeraCrypt)
  * Other projects/utilities should be able to use this mechanism with minimal effort
    * Should be able to invoke a single line during run script to allow this to work for other projects
    
      ```bash
      grypt -d -p password -f /path/to/secret/file &    # Run grypt to decrypt the file and send it to background
      GRYPT_PID=$!                                      # Get pid to be able to easily kill it later
      myGreatProgram --config /path/to/secret/file      # Run my program using file now decrypted to plaintext
      kill $GRYPT_PID                                   # Kill grypt to automatically re-encrypt the file
      ```
    * Integration points:
        - [x] Bash/filesystem
        - [ ] Docker/filesystem
        - [ ] Golang hook
    

# Encryption Mechanisms

## AES-GCM-256/PEM
Encrypted using AES-GCM-256/PEM with additionnal datas (to protect PEM headers) instead of Salted CBC-128
  - `grypt -e -f /path/to/secret` takes the plaintext file and encrypts it
  - `grypt -d -f /path/to/secret` takes the encrypted file and decrypts it. When grypt is sigterm'd, the file will be read in, and re-encrypted using the same password
  - `grypt -D -f /path/to/secret` takes the encrypted file and converts it into plain text

Grypt encryption uses:
  - AEAD Authenticated Encryption Additionnal Data modes (protect the plaintext PEM headers)
  - AES-GCM-256 authenticated encryption mode.
  - 16K rounds PBKDF2 key derivation function with SHA3-256
  - Crypto PRNG.
  

# Overview of Task
docker decrypt file
> password
mount that file in container to do something
> input for program
docker exec cleanup
    - if plaintext contents were changed
        > Do you want to save the changed? [y]/[n]
        > password
    - reencrypt the plaintext contents
    - delete the plaintext contents from the mounted volume
    - save the encrypted contents in the other mounted volume.


# Notes
  * otp project will use grypt under the hood, but will have config files based off veracrypt providing the keys
    * CLI #1 - User does not need to interact with the CLI if keys are provided as arguments in docker invocation
        ```bash
        > vera
        > # Enters Passwords for single or hidden volume
        > otp # Uses keys from veracrypt volume to unlock encrypted files
        ```
    * CLI #2 -
        ```bash
        > vera
        > 1234 # to unlock veracrypt volume with the keys and encrypted files
        > otp # Uses keys from veracrypt volume to unlock encrypted files
        > password # to unlock the encrypted file with grypt
        ```





# TODO
  - [ ] Implement Cleanup
    * Should reencrypt if -d flag is given, not reencrypt if -D flag is given
  - [ ] Assure great UX:
    * `gcrypt --encrypt -f /path/to/file`
    * `> Please enter the password to encrypt:`
    * `> Re-type password:`
  - [ ] More CLI flags:
    - [ ] Add flag for inline password
      * `grcypt --encrypt -k /path/to/keyfile -f /path/to/file`
    - [ ] Add flag for inline keyfile
      * `grcypt --encrypt -p password -f /path/to/file`
    - [ ] Add flag for modifying crypto or method
    - [ ] Add flag to change crypto strength