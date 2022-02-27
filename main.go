package main

import (
	"mediaindexer/indexer"
	"os/signal"
	"os"
	"syscall"
)

func main() {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        indexer.SaveItems()
        os.Exit(1)
    }()

	indexer.Run()
}
