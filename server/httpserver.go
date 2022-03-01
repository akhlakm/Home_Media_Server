package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"strings"
	"mediaindexer/indexer"
)

var httpServer *http.Server = nil

func serveAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		fmt.Println("Req:", r.URL.Path)
		// get the parts
		ss := strings.Split(r.URL.Path, "/")
		if len(ss) == 4 {
			hash := ss[3]
			if len(hash) == 0 {
				w.Write([]byte("400"))
			} else if ss[2] == "like" {
				indexer.AddLike(hash)
				w.Write([]byte("200"))
			} else if ss[2] == "dislike" {
				indexer.AddDislike(hash)
				w.Write([]byte("200"))
			} else {
				w.Write([]byte("404"))
			}
		} else {
			w.Write([]byte("Hello world from " + r.URL.Path))
		}
	}
}

func RunServer(wg *sync.WaitGroup, embedFs fs.FS, port int) {
	defer wg.Done()

	httpServer = &http.Server{Addr: fmt.Sprintf(":%d", port)}

	// register the websocket url
	http.HandleFunc("/api/", serveAPI)

	// register the embeded public folder to /
	fsys := http.FileServer(http.FS(embedFs))
	http.Handle("/", http.StripPrefix("/", fsys))

	fmt.Printf("Server started at port %d, api url %s ... \n", port, "/api")
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
	fmt.Println("Server exit!")
}

//
// Gracefully stop the http server.
//
func Shutdown() {
	if httpServer != nil {
		httpServer.Close()
		fmt.Println("Server shutdown!")
	}
}
