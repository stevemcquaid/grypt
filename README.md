# Grypt: Golang Encryption Utility

An encryption utility that is a Work in Progress. Decryption is having issues at the moment

### Goal #1 - Protect secrets stored in public by providing easy-to-use strong encryption
  * Used in place of, or in combination with, [Veracrypt](https://en.wikipedia.org/wiki/VeraCrypt)
  * Provide a mechanism by which secrets are unencrypted during use, but encrypted at rest
### Goal #2 - Other projects/utilities should be able to use this mechanism with minimal effort
  * Provide a secure way for files containing secrets to processed by programs in plain text for simplicity, but be stored securely

# Grypt Encryption Mechanisms
  - AES-GCM-256/PEM with additionnal datas (to protect PEM headers) instead of Salted CBC-128
  - AEAD Authenticated Encryption Additionnal Data modes (protect the plaintext PEM headers)
  - AES-GCM-256 authenticated encryption mode.
  - 16K rounds PBKDF2 key derivation function with SHA3-256
  - Crypto PRNG.

# Using this repository
A few `make` targets have been provided for convienance in building the docker container and running it
  - `make help` to see all available targets [Source](Makefile#L22)
  - `make build` to build the docker container [Source](scripts/build.sh)
  - `make bash` to run a bash shell within the docker container [Source](scripts/bash.sh)
  - `make example` to run grypt in the most basic case [Source](scripts/example.sh)

# grypt CLI
  - `grypt -e -f /path/to/secret` takes the plaintext file and encrypts it
  - `grypt -d -f /path/to/secret` takes the encrypted file and decrypts it. When grypt is sigterm'd, the file will be read in, and re-encrypted using the same password
  - `grypt -D -f /path/to/secret` takes the encrypted file and converts it into plain text

# Docker Usage
  Run the grypt CLI from within a docker container to avoid complex dependency chains.  Utilize mounted volumes to transfer data in & out of the docker container.
  Variables for simplicity:
  ```bash
  LOCAL_SECRETS_PATH=/full/host/machine/path/to/secrets/
  DOCKER_SECRETS_PATH=/secrets
  ```
  - `docker run -it -v $LOCAL_SECRETS_PATH:$DOCKER_SECRETS_PATH --rm stevemcquaid/grypt:latest grypt -f $DOCKER_SECRETS_PATH/file -e` takes the plaintext file and encrypts it
  - `docker run -it -v $LOCAL_SECRETS_PATH:$DOCKER_SECRETS_PATH --rm stevemcquaid/grypt:latest grypt -f $DOCKER_SECRETS_PATH/file -d` takes the encrypted file and decrypts it temporarily. Leave the docker container running for as long as you want the file to remain in plain text. When you `ctrl-c`/kill the docker container, grypt will be sent SIGTERM, and before quitting, the file will be read in again, and re-encrypted using the same password
  - `docker run -it -v $LOCAL_SECRETS_PATH:$DOCKER_SECRETS_PATH --rm stevemcquaid/grypt:latest grypt -f $DOCKER_SECRETS_PATH/file -D` takes the encrypted file and converts it into plain text


# TODO
  - [x] Implement Cleanup
    - [x] Should reencrypt if -d flag is given, not reencrypt if -D flag is given
  - [x] Improve UX:
      ```bash
      gcrypt --encrypt -f /path/to/file
      > Please enter the password to encrypt:
      > Re-type password:
      ```
  - [ ] Add 1 unit test
  - [ ] More CLI flags:
    - [ ] Add verbose flag options
      * `--encrypt --decryptTemp --decryptForce --file --help`
    - [ ] Add flag for inline password
      * `grcypt --encrypt -k /path/to/keyfile -f /path/to/file`
    - [ ] Add flag for inline keyfile
      * `grcypt --encrypt -p password -f /path/to/file`
    - [ ] Add flag to change crypto or method
    - [ ] Add flag to change crypto strength
  - [ ] Should be able to invoke a single line during run script to allow this to work for other projects
      ```bash
      grypt -d -p password -f /path/to/secret/file &    # Run grypt to decrypt the file and send it to background
      GRYPT_PID=$!                                      # Get pid to be able to easily kill it later
      myGreatProgram --config /path/to/secret/file      # Run my program using file now decrypted to plaintext
      kill $GRYPT_PID                                   # Kill grypt to automatically re-encrypt the file
      ```
  - [x] Multiple integration points:
      - [x] Bash/filesystem
      - [x] Docker/filesystem
      - [x] Golang library

    