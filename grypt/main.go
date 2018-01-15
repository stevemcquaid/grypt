package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	// "time" // or "runtime"
)

var file string
var encryptOrNot bool

func argParse() {

	// @TODO no default allowed
	flag.StringVar(&file, "file", "", "Absolute path to file to encrypt")

	decryptFlag := flag.Bool("d", false, "decrypt file")
	encryptFlag := flag.Bool("e", false, "encrypt file")
	// tempDecrypt := flag.Bool("tempD", false, "Decrypt only for duration of program")

	flag.Parse()

	if *encryptFlag == true && *decryptFlag == true {
		// trying to encrypt and decrypt at the same time!
		fmt.Printf("-e and -d options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == false && *decryptFlag == false {
		// You need to do something!
		fmt.Printf("-e or -d option must be set")
		os.Exit(1)
	} else if *encryptFlag == true {
		encryptOrNot = true
	} else if *decryptFlag == true {
		encryptOrNot = false
	}
}

// Blocking function to ask user for password
func askPass(prompt1 string, prompt2 string) ([]byte, error) {
	// XXX do the encryption
	fmt.Printf(prompt1)
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	// if err != nil {
	// 	return
	// }
	fmt.Printf("\n")

	fmt.Printf(prompt2)
	rpasswd, err := terminal.ReadPassword(syscall.Stdin)
	// if err != nil {
	// 	return
	// }
	fmt.Printf("\n")

	if bytes.Compare(passwd, rpasswd) != 0 {
		err = fmt.Errorf("Passwords don't match\n")
		return nil, err
	}

	return passwd, nil
}

func run() {
	// #Possible CLI flags
	// file / path
	// password
	// keyfile
	// crypto - method
	// strength

	// Usage examples:
	// grcypt --encrypt -p password -f /path/to/file
	// gcrypt --encrypt -f /path/to/file
	// > Please enter the password to encrypt:
	// > Re-type password:

	argParse()

	if encryptOrNot == true {
		// Encrypt
		passwd, _ := askPass("Please enter the password to encrypt: ", "Re-type password: ")
		fmt.Println("Encrypting...")
		encrypt(file, passwd)
	} else {
		// Decrypt
		passwd, _ := askPass("Please enter the password to decrypt: ", "Re-type password: ")
		fmt.Println("Decrypting...")
		decrypt(file, passwd)
	}
}

// <--- Control Loop --->

func setupSigtermHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()
}

func wait() {
	fmt.Println("Sleeping...")
	select {}
	// for {
	//     fmt.Println("Sleeping...")
	//     time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	// }
}

func cleanup() {
	fmt.Println("Cleanup")
}

func main() {
	setupSigtermHandler()
	run()
	wait()
}
