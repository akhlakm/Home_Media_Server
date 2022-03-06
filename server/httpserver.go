package server

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"bytes"
	"os"
	"net/http"
	"sync"
	"strings"
	"mediaindexer/indexer"
	"encoding/base64"
	"image/png"
)

var httpServer *http.Server = nil

func handleApiReq(cmd string, urlParts []string) string {
	if cmd == "walk" {
		go indexer.Run()
		return "200"
	} else if cmd == "like" && len(urlParts) > 4 {
		// host/api/like/<hash>
		hash := urlParts[3]
		indexer.AddLike(hash)
		return "200"
	} else if cmd == "dislike" && len(urlParts) > 4 {
		// host/api/dislike/<hash>
		hash := urlParts[3]
		indexer.AddDislike(hash)
		return "200"
	}

	return "404"
}
// This function returns the filename(to save in database) of the saved file
// or an error if it occurs
func fileUpload(r *http.Request) (string, error) {
    // ParseMultipartForm parses a request body as multipart/form-data
    r.ParseMultipartForm(32 << 20)

    file, handler, err := r.FormFile("file") // Retrieve the file from form data

    if err != nil {
        return "", err
    }
    defer file.Close()                       // Close the file when we finish

    // This is path which we want to store the file
    f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

    if err != nil {
        return "", err
    }

    // Copy the file to the destination path
    io.Copy(f, file)

    return handler.Filename, nil
}

func uploadB64(req *http.Request) string {
	imgString := req.FormValue("image")
	fmt.Printf("Upload: %s, %d, %T\n", imgString, len(imgString), imgString)

	unbased, err := base64.StdEncoding.DecodeString(imgString)
	if err != nil {
		panic("Cannot decode b64")
	}
	
	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		panic("Bad png")
	}
	
	f, err := os.OpenFile("uploaded.png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}
	
	png.Encode(f, im)
	return "200"
}

func serveAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Println("Req:", r.URL.Path)

	if r.Method == "GET" {
		// get the parts host/api/like/<hash>
		ss := strings.Split(r.URL.Path, "/")
		w.Write([]byte(handleApiReq(ss[2], ss)))
	} else if r.Method == "POST" {
		w.Write([]byte(uploadB64(r)))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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
