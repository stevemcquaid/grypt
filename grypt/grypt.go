package main

import (
	"bytes"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	HDR_PEM = "TOPT KEYFILE"
)

func encrypt(file string, password string) {
	fmt.Println("Encrypting...")
}

func askPassAndEncryptTotpFile(ofile, ifile string) (err error) {
	cfgPlainContent, err := ioutil.ReadFile(ifile)
	if err != nil {
		return
	}

	// XXX do the encryption
	fmt.Printf("password: ")
	passwd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return
	}
	fmt.Printf("\n")

	fmt.Printf("retype password: ")
	rpasswd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return
	}
	fmt.Printf("\n")

	if bytes.Compare(passwd, rpasswd) != 0 {
		err = fmt.Errorf("Passwords don't match\n")
		return
	}

	// write the new file
	cfgContentBlock, err := AEADEncryptPEMBlock(rand.Reader, HDR_PEM, cfgPlainContent, passwd)
	if err != nil {
		err = fmt.Errorf("Encryption failure (%s).\n", err.Error())
		//return fmt.Errorf("encryption problem\n")
		return
	}

	cfgPemContent := pem.EncodeToMemory(cfgContentBlock)
	err = ioutil.WriteFile(ofile, cfgPemContent, 0600)
	if err != nil {
		return err
	}

	err = os.Remove(ifile)
	if err != nil {
		fmt.Printf("warning could not remove %s.\n", ifile)
	}

	fmt.Printf("encrypted to %s\n", ofile)
	return nil
}
