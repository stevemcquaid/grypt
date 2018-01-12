package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    // "time" // or "runtime"
)

func setup(){
    fmt.Println("...")
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

func wait(){
    fmt.Println("Sleeping...")
    select{}
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
    fmt.Println("Running setup...")
    setup()
    fmt.Println("Done with setup.")
    wait()
}