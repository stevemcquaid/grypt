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