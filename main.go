package main

import (
	"mediaindexer/indexer"
	"mediaindexer/server"
	"flag"
    "embed"
    "io/fs"
    "log"
    "fmt"
	"os/signal"
	"os"
	"syscall"
    "sync"
)

//go:embed www
var embedFs embed.FS

func main() {
    // Named args
	root := flag.String("d", "/media/i/", "The inbox directory path to index")
	walk := flag.Bool("walk", false, "Specify to walk the root directory")
	serve := flag.Bool("serve", false, "Specify to serve on HTTP")
	flag.Parse()

	www := "/media/i/"

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        if *walk {
            indexer.SaveItems()
        }
        if *serve {
            server.Shutdown()
        }
        fmt.Println("Goodbye!")
        os.Exit(0)
    }()

    
    // embed the public directory contents ...
	fsys, err := fs.Sub(embedFs, "www")
	if err != nil {
        log.Fatalf("Failed to create embeded filesystem for the www/ directory.")
	}

    var wg sync.WaitGroup

    
    if *serve {
        wg.Add(1)
        go server.RunServer(&wg, fsys, 9000)
    }

    indexer.Run(*root, www, *walk)
    wg.Wait()
}
