# grypt - Usage

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

# Task
  * Inputs:
    * keys || stdin password
  * Output:
    * decrypted secrets in volume
  * Cleanup:
    * re-encrypt secrets using inputs

# Goals
  * Used in place of veracrypt
  * Other projects should be able to use this easily
    * otp project migration to grypt should be nearly drop in
        * (docker?)
        * golang hook maybe?
    * Should be able to invoke a single line to allow this to work in other projects?

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



# Encryption Mechanisms

## AES-GCM-256/PEM
Encrypted using AES-GCM-256/PEM with additionnal datas (to protect PEM headers) instead of Salted CBC-128

`gauth -e` take the current ~/.config/gauth.csv and encrypts it to ~/.config/gauth.pem and remove the plaintext version.
`gauth -d` if you need to peek/poke in your token file, then `gauth -e` again.

gauth TOTP keyfile encryption uses:
  - AEAD Authenticated Encryption Additionnal Data modes (protect the plaintext PEM headers)
  - AES-GCM-256 authenticated encryption mode.
  - 16K rounds PBKDF2 key derivation function with SHA3-256
  - Crypto PRNG.

# TODO
  * Implement Cleanup
    - `-f` force flag ignores cleanup
      - `-e -f`
      - `-d -f`

  * #Possible CLI flags
    * file / path
    * password
    * keyfile
    * crypto - method
    * strength

	* Usage examples:
    * grcypt --encrypt -p password -f /path/to/file
    * gcrypt --encrypt -f /path/to/file
    * > Please enter the password to encrypt:
    * > Re-type password: