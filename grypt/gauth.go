package main

// import (
// 	"bytes"
// 	"crypto/hmac"
// 	"crypto/rand"
// 	"crypto/sha1"
// 	"encoding/base32"
// 	"encoding/pem"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"math/big"
// 	"os"
// 	"os/user"
// 	"path"
// 	"syscall"

// 	"golang.org/x/crypto/ssh/terminal"
// )

// const (
// 	CONFIG_PLAIN = ".config/gauth.csv"
// 	CONFIG_PEM   = ".config/gauth.pem"

// 	HDR_PEM = "TOPT KEYFILE"
// )

// func AuthCode(sec string, ts int64) (string, error) {
// 	key, err := base32.StdEncoding.DecodeString(sec)
// 	if err != nil {
// 		return "", err
// 	}
// 	enc := hmac.New(sha1.New, key)
// 	msg := make([]byte, 8, 8)
// 	msg[0] = (byte)(ts >> (7 * 8) & 0xff)
// 	msg[1] = (byte)(ts >> (6 * 8) & 0xff)
// 	msg[2] = (byte)(ts >> (5 * 8) & 0xff)
// 	msg[3] = (byte)(ts >> (4 * 8) & 0xff)
// 	msg[4] = (byte)(ts >> (3 * 8) & 0xff)
// 	msg[5] = (byte)(ts >> (2 * 8) & 0xff)
// 	msg[6] = (byte)(ts >> (1 * 8) & 0xff)
// 	msg[7] = (byte)(ts >> (0 * 8) & 0xff)
// 	if _, err := enc.Write(msg); err != nil {
// 		return "", err
// 	}
// 	hash := enc.Sum(nil)
// 	offset := hash[19] & 0x0f
// 	trunc := hash[offset : offset+4]
// 	trunc[0] &= 0x7F
// 	res := new(big.Int).Mod(new(big.Int).SetBytes(trunc), big.NewInt(1000000))
// 	return fmt.Sprintf("%06d", res), nil
// }

// func askPassAndEncryptTotpFile(ofile, ifile string) (err error) {
// 	cfgPlainContent, err := ioutil.ReadFile(ifile)
// 	if err != nil {
// 		return
// 	}

// 	// XXX do the encryption
// 	fmt.Printf("password: ")
// 	passwd, err := terminal.ReadPassword(syscall.Stdin)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Printf("\n")

// 	fmt.Printf("retype password: ")
// 	rpasswd, err := terminal.ReadPassword(syscall.Stdin)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Printf("\n")

// 	if bytes.Compare(passwd, rpasswd) != 0 {
// 		err = fmt.Errorf("Passwords don't match\n")
// 		return
// 	}

// 	// write the new file
// 	cfgContentBlock, err := AEADEncryptPEMBlock(rand.Reader, HDR_PEM, cfgPlainContent, passwd)
// 	if err != nil {
// 		err = fmt.Errorf("Encryption failure (%s).\n", err.Error())
// 		//return fmt.Errorf("encryption problem\n")
// 		return
// 	}

// 	cfgPemContent := pem.EncodeToMemory(cfgContentBlock)
// 	err = ioutil.WriteFile(ofile, cfgPemContent, 0600)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.Remove(ifile)
// 	if err != nil {
// 		fmt.Printf("warning could not remove %s.\n", ifile)
// 	}

// 	fmt.Printf("encrypted to %s\n", ofile)
// 	return nil
// }

// func askPassAndDecryptTotpFile(ofile, ifile string) (err error) {

// 	cfgContent, err := ioutil.ReadFile(ifile)
// 	if err != nil || IsEncryptedPemFile(ifile) == false {
// 		err = fmt.Errorf("Non-existent/Invalid encrypted TOTP keyfile (%s).", err.Error())
// 		return
// 	}

// 	cfgPemBlock, _ := pem.Decode(cfgContent)
// 	if cfgPemBlock == nil || cfgPemBlock.Type != HDR_PEM {
// 		err = fmt.Errorf("Invalid TOTP keyfile PEM Block\n")
// 		return
// 	}

// 	fmt.Printf("password: ")
// 	passwd, err := terminal.ReadPassword(syscall.Stdin)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Printf("\n")

// 	cfgPlainContent, err := AEADDecryptPEMBlock(cfgPemBlock, passwd)
// 	if err != nil {
// 		err = fmt.Errorf("Invalid password/encrypted payload (%s)\n", err.Error())
// 		return //fmt.Errorf("invalid password\n")
// 	}

// 	err = ioutil.WriteFile(ofile, cfgPlainContent, 0600)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("decrypted to %s\n", ofile)
// 	return nil
// }

// func authDecryptParseTotpFile(ifile string) (r [][]string, err error) {
// 	// decrypt
// 	cfgPemContent, err := ioutil.ReadFile(ifile)
// 	if err != nil || IsEncryptedPemFile(ifile) == false {
// 		err = fmt.Errorf("Non-existent/Invalid encrypted TOTP keyfile (%s).", err.Error())
// 		return
// 	}

// 	cfgPemBlock, _ := pem.Decode(cfgPemContent)
// 	if cfgPemBlock == nil || cfgPemBlock.Type != HDR_PEM {
// 		err = fmt.Errorf("Invalid TOTP keyfile PEM Block\n")
// 		return
// 	}

// 	fmt.Printf("password: ")
// 	passwd, err := terminal.ReadPassword(syscall.Stdin)
// 	if err != nil {
// 		return
// 	}
// 	fmt.Printf("\n")

// 	cfgContent, err := AEADDecryptPEMBlock(cfgPemBlock, passwd)
// 	if err != nil {
// 		err = fmt.Errorf("Invalid password/encrypted payload (%s)\n", err.Error())
// 		return
// 	}

// }

// //
// // default, try to open PEM file by default
// // then try the PLAIN file, WARN if plain text.
// //
// // -e : encrypt gauth.csv to gauth.pem (or file config)
// // if file already exist, ask for password, rencrypt/truncate with the same pass.
// // if file does not exist, ask for password, verify password, create the new file.
// //
// // -d : decrypt gauth.pem to gauth.csv (or file config)
// //

// func main() {

// 	user, err := user.Current()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// decrypt the file take the pem and generate a csv (truncate).
// 	decryptFlag := flag.Bool("d", false, "decrypt TOTP keyfile (~/.config/gauth.pem -> gauth.csv)")
// 	// if the timestamp of csv is > pem and you can decrypt pem, then
// 	// reencrypt the csv and replace pem, otherwise fail.
// 	encryptFlag := flag.Bool("e", false, "encrypt TOTP keyfile (~/.config/gauth.csv -> gauth.pem)")

// 	var cfgPem string
// 	flag.StringVar(&cfgPem, "CONFIG_PEM", path.Join(user.HomeDir, CONFIG_PEM), "Absolute path to PEM config file to use")

// 	var cfgPlain string
// 	flag.StringVar(&cfgPlain, "CONFIG_PLAIN", path.Join(user.HomeDir, CONFIG_PLAIN), "Absolute path to PLAIN CSV config file to use")

// 	flag.Parse()
// 	sArgs := flag.Args()

// 	// trying to encrypt and decrypt at the same time?!
// 	if *encryptFlag == true && *decryptFlag == true {
// 		fmt.Printf("-e and -d options are mutually exclusive")
// 		os.Exit(1)
// 	}

// 	if *decryptFlag == true {
// 		err := askPassAndDecryptTotpFile(cfgPlain, cfgPem)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		os.Exit(0)

// 	}

// 	if *encryptFlag == true {
// 		err := askPassAndEncryptTotpFile(cfgPem, cfgPlain)
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
// 		os.Exit(0)

// 	} // end of if encryptFlag

// 	// default behaviour
// 	cfg, err := authDecryptParseTotpFile(cfgPem)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	currentTS, progress := TimeStamp()
// 	prevTS := currentTS - 1
// 	nextTS := currentTS + 1

// 	for _, record := range cfg {
// 		name := record[0]
// 		secret := normalizeSecret(record[1])
// 		prevToken := authCodeOrDie(secret, prevTS)
// 		currentToken := authCodeOrDie(secret, currentTS)
// 		nextToken := authCodeOrDie(secret, nextTS)
// 		fmt.Printf(HDR_FMT, name, prevToken, currentToken, nextToken)
// 	}
// }
