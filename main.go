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
	root := flag.String("d", "/media/i/_inbox", "The inbox directory path to index")
	walk := flag.Bool("walk", false, "Specify to walk the root directory")
	serve := flag.Bool("serve", false, "Specify to serve on HTTP")
	flag.Parse()

    // where the files will be moved after adding to database from the root inbox
	www := "/media/i/"

    indexer.Init(*root, www)

    // Catch Ctrl+C and sigterm events for graceful shutdown
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        indexer.SaveItems()
        if *serve {
            server.Shutdown()
        }
        fmt.Println("Goodbye!")
        os.Exit(0)
    }()

    
    // embed the www directory contents ...
	fsys, err := fs.Sub(embedFs, "www")
	if err != nil {
        log.Fatalf("Failed to create embeded filesystem for the www/ directory.")
	}

    var wg sync.WaitGroup

    
    if *serve {
        wg.Add(1)
        go server.RunServer(&wg, fsys, 9000)
    } else {
        fmt.Println("Specify -serve to start the media server at :9000")
    }

    if *walk {
        go indexer.Run()
    } else {
        fmt.Println("Specify -walk to crawl the root inbox directory. Or request /api/walk/")
    }

    wg.Wait()
}
