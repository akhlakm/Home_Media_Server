package server

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"net/http"
	"sync"
	"strings"
	"mediaindexer/indexer"
	"time"
)

var httpServer *http.Server = nil
// wait 1 min after user leaves, 10 secs grace period
const timerWait = 70 * time.Second;
var hbTimer *time.Timer = nil

func handleApiReq(cmd string, urlParts []string) string {
	if cmd == "walk" {
		go indexer.Run()
		return "200"
	} else if cmd == "heartbeat" {
		if hbTimer == nil {
			hbTimer = time.NewTimer(timerWait)
			go func() {
				<-hbTimer.C
				hbTimer.Stop()
				hbTimer = nil
				// things to do after 1 min of user absence
				indexer.SaveItems()
			}()
		} else {
			hbTimer.Reset(timerWait);
		}
		return "200"
	} else if cmd == "like" && len(urlParts) > 3 {
		// host/api/like/<hash>
		hash := urlParts[3]
		indexer.AddLike(hash)
		return "200"
	} else if cmd == "dislike" && len(urlParts) > 3 {
		// host/api/dislike/<hash>
		hash := urlParts[3]
		indexer.AddDislike(hash)
		return "200"
	}

	return "404"
}

func fileUpload(r *http.Request, hash string) string {
    // ParseMultipartForm parses a request body as multipart/form-data
    r.ParseMultipartForm(32 << 20)

    file, _, err := r.FormFile("file") // Retrieve the file from form data
	if err != nil {
        return "500"
    }
    defer file.Close()                       // Close the file when we finish

    // This is path which we want to store the file
    f, err := os.OpenFile("_uploaded.png", os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        return "500"
    }
	defer f.Close()

	// Copy the file to the destination path
    io.Copy(f, file)

	// move to place in background
	go indexer.AddCaptionFile("_uploaded.png", hash)

    return "200"
}

func serveAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Println(r.Method, r.URL.Path)

	// get the parts host/api/like/<hash>
	ss := strings.Split(r.URL.Path, "/")

	if r.Method == "GET" {		
		w.Write([]byte(handleApiReq(ss[2], ss)))
		return
	} else if r.Method == "POST" {
		if len(ss) > 3 {
			w.Write([]byte(fileUpload(r, ss[3])))
			return
		}
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

func RunServer(wg *sync.WaitGroup, embedFs fs.FS, port int) {
	defer wg.Done()

	httpServer = &http.Server{Addr: fmt.Sprintf(":%d", port)}

	// register the websocket url
	http.HandleFunc("/api/", serveAPI)

	// register the embeded www folder to /
	fsys := http.FileServer(http.FS(embedFs))
	http.Handle("/", http.StripPrefix("/", fsys))

	// register the currect directory as a static folder
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("./"))))

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
