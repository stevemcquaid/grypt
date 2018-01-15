package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	// "time" // or "runtime"
)

var password string
var file string

func argParse() {
	// #Possible CLI flags
	// crypto - method
	// strength
	// keyfile
	// password
	// file / path

	// @TODO no default allowed
	flag.StringVar(&file, "file", "", "Absolute path to file to encrypt")
	// @TODO no default allowed
	flag.StringVar(&password, "password", "", "Password provided inline")

	flag.Parse()

	// convert file to path obj
	path.Join("", file)
}

func run() {

	// TODO password + encrypt single file using AES-GCM-256/PEM
	// Usage examples:
	// grcypt --encrypt -p password -f /path/to/file
	// gcrypt --encrypt -f /path/to/file
	// > Please enter the password to encrypt:
	// > Re-type password:

	encrypt("/Users/steve/Dropbox/Code/crypt/grypt/test", "testtest")
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
	fmt.Println("Parsing Args...")
	argParse()
	run()
	wait()
}
