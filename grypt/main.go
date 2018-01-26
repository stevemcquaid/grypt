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
var doCleanup bool
var password []byte

func argParse() {
	// @TODO no default allowed
	flag.StringVar(&file, "file", "", "Absolute path to file to encrypt")

	tempDecryptFlag := flag.Bool("d", false, "Temporarily decrypt file, re-encrypt on exit")
	forceDecryptFlag := flag.Bool("D", false, "Force decrypt file, show plaintext even after exit")
	encryptFlag := flag.Bool("e", false, "Encrypt file")

	flag.Parse()

	if *forceDecryptFlag == true && *tempDecryptFlag == true {
		// Trying to tempDecrypt and forceDecrypt at the same time!
		fmt.Printf("-d and -D options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == true && *tempDecryptFlag == true {
		// Trying to encrypt and tempDecrypt at the same time!
		fmt.Printf("-e and -d options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == true && *forceDecryptFlag == true {
		// Trying to encrypt and forceDecrypt at the same time!
		fmt.Printf("-e and -D options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == false && *tempDecryptFlag == false && *forceDecryptFlag == false {
		// You need to do something!
		fmt.Printf("-e or -d or -D option must be set")
		os.Exit(1)
	} else if *encryptFlag == true {
		encryptOrNot = true
	} else if *forceDecryptFlag == true {
		encryptOrNot = false
		doCleanup = false
	} else if *tempDecryptFlag == true {
		encryptOrNot = false
		doCleanup = true
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

	argParse()

	if encryptOrNot == true {
		// Encrypt
		passwd, _ := askPass("Please enter the password to encrypt: ", "Re-type password: ")
		fmt.Println("Encrypting...")
		encrypt(file, passwd)
	} else {
		// Decrypt
		passwd, _ := askPass("Please enter the password to decrypt: ", "Re-type password: ")
		// Only save password in global if needed
		if doCleanup {
			password = passwd
		}
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
	if doCleanup {
		fmt.Println("Starting cleanup...")
		encrypt(file, password)
		fmt.Println("Completed cleanup.")
	} else {
		fmt.Println("No cleanup needed.")
	}
}

func main() {
	setupSigtermHandler()
	run()
	wait()
}
