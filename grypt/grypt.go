package main

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	HDR_PEM = "TOPT KEYFILE"
)

func encrypt(file string, passwd []byte) (err error) {
	ifile := file
	ofile := file

	cfgPlainContent, err := ioutil.ReadFile(ifile)
	if err != nil {
		err = fmt.Errorf("Could not read file (%s).\n", ifile)
		return
	}

	// write the new file
	cfgContentBlock, err := AEADEncryptPEMBlock(rand.Reader, HDR_PEM, cfgPlainContent, passwd)
	if err != nil {
		err = fmt.Errorf("Encryption failure (%s).\n", err.Error())
		return
	}

	err = os.Remove(ifile)
	if err != nil {
		err = fmt.Errorf("Could not remove (%s).\n", ifile)
		return
	}

	cfgPemContent := pem.EncodeToMemory(cfgContentBlock)
	err = ioutil.WriteFile(ofile, cfgPemContent, 0600)
	if err != nil {
		err = fmt.Errorf("Could not write file (%s).\n", ofile)
		return
	}

	fmt.Printf("Encryption complete. (%s)\n", ofile)
	return nil
}

func decrypt(file string, passwd []byte) (err error) {
	ifile := file
	ofile := file

	cfgContent, err := ioutil.ReadFile(ifile)
	if err != nil || IsEncryptedPemFile(ifile) == false {
		err = fmt.Errorf("Non-existent/Invalid encrypted TOTP keyfile (%s).", err.Error())
		return
	}

	cfgPemBlock, _ := pem.Decode(cfgContent)
	if cfgPemBlock == nil || cfgPemBlock.Type != HDR_PEM {
		err = fmt.Errorf("Invalid TOTP keyfile PEM Block\n")
		return
	}

	cfgPlainContent, err := AEADDecryptPEMBlock(cfgPemBlock, passwd)
	if err != nil {
		err = fmt.Errorf("Invalid password/encrypted payload (%s)\n", err.Error())
		return //fmt.Errorf("invalid password\n")
	}

	err = ioutil.WriteFile(ofile, cfgPlainContent, 0600)
	if err != nil {
		return err
	}

	fmt.Printf("Decryption complete. (%s)\n", ofile)
	return nil
}
