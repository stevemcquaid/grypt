package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
	"golang.org/x/crypto/ssh/terminal"
	// "time" // or "runtime"
)

var file string
var encryptOrNot bool
var doCleanup bool
var password []byte

func argParse() {
	fileFlag := flag.String("f", "", "Absolute path to file to encrypt")
	tempDecryptFlag := flag.Bool("d", false, "Temporarily decrypt file, re-encrypt on exit")
	forceDecryptFlag := flag.Bool("D", false, "Force decrypt file, show plaintext even after exit")
	encryptFlag := flag.Bool("e", false, "Encrypt file")

	flag.Parse()

	if *fileFlag == "" {
		// You need to have a file to do something
		fmt.Println("File flag is required: -f")
		os.Exit(1)
	} else {
		file = *fileFlag
	}

	if *forceDecryptFlag == true && *tempDecryptFlag == true {
		// Trying to tempDecrypt and forceDecrypt at the same time!
		fmt.Println("-d and -D options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == true && *tempDecryptFlag == true {
		// Trying to encrypt and tempDecrypt at the same time!
		fmt.Println("-e and -d options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == true && *forceDecryptFlag == true {
		// Trying to encrypt and forceDecrypt at the same time!
		fmt.Println("-e and -D options are mutually exclusive")
		os.Exit(1)
	} else if *encryptFlag == false && *tempDecryptFlag == false && *forceDecryptFlag == false {
		// You need to do something!
		fmt.Println("-e or -d or -D option must be set")
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
	fmt.Printf(prompt1)
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n")

	fmt.Printf(prompt2)
	rpasswd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n")

	if bytes.Compare(passwd, rpasswd) != 0 {
		err = fmt.Errorf("Passwords don't match")
		return nil, err
	}
	return passwd, nil
}

func run() {
	argParse()

	if encryptOrNot == true {
		// Encrypt
		passwd, _ := askPass("Please enter the password to encrypt: ", "Re-type password: ")
		glog.V(2).Infof("Encrypting...")
		encrypt(file, passwd)
		// No need to cleanup - exit now
		exit()
	} else {
		// Decrypt
		passwd, _ := askPass("Please enter the password to decrypt: ", "Re-type password: ")

		if doCleanup {
			// Setup handler to wait for exit to reencrypt
			setupSigtermHandler()

			// Only save password in global if needed
			password = passwd
		}
		glog.V(2).Infof("Decrypting...")
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
		exit()
	}()
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

func wait() {
	fmt.Println("Sleeping...")
	select {}
	// for {
	//     fmt.Println("Sleeping...")
	//     time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	// }
}

func exit() {
	// @TODO Could handle cleanup of sigtermHandler here for cleaner control loop
	glog.Flush()
	os.Exit(0)
}

// <--- End Control Loop --->

func main() {
	run()
	wait()
}
