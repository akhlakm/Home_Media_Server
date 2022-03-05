package indexer

import (
	"encoding/json"
	"sync"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileItem struct {
	Path   string  `json:"Path"`
	Desc   string  `json:"Desc"`
	SizeMB float32 `json:"SizeMB"`
	URL    string  `json:"URL"`
	Likes  		int  	`json:"Likes"`
	Dislikes  	int  	`json:"Disikes"`
}

var HASHLEN int = 16
var list map[string]FileItem
var jsonfile string
var jsfile string

func AddLike(hash string) {
	if fItem, exists := list[hash]; exists {
		fmt.Println("Like:", hash)
		fItem.Likes++
		list[hash] = fItem
		SaveItems()
	}
}

func AddDislike(hash string) {
	if fItem, exists := list[hash]; exists {
		fmt.Println("Dislike:", hash)
		fItem.Dislikes++
		list[hash] = fItem
		SaveItems()
	}
}


func convertMP4(root, www string) {
	// convert the videos to mp4
	fmt.Println("Checking for non MP4 videos ...")
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() && IsFileType(s, Videos2Convert) {
			if ffmpeg_convert_mp4(s) {
				os.Remove(s)
			}
		}
		return nil
	})
}

func segmentVideo(root, www string) {	
	// segmentize videos
	fmt.Println("Checking for long MP4 videos ...")
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		// ignore items if we are inside the www directory
		wwwmp4 := www + "mp4/"
		if strings.HasPrefix(s, wwwmp4) {
			return nil
		}
		// already segmented
		if strings.HasPrefix(filepath.Base(s), "part.") {
			return nil
		}
		if !d.IsDir() && IsFileType(s, SupportedVideos) {
			// get video info
			probe := ffprobe_video(s)
			duration := extract_duration(probe)
			// segmentize video if longer than 20 secs
			if duration > 20 {
				hash, err := chunk_md5(s, HASHLEN)
				if err != nil {
					return err
				}
				if ffmpeg_segment(hash, s, 20) {
					// remove the original file
					os.Remove(s)
					// sleep after segmentation
					time.Sleep(time.Duration(100) * time.Millisecond)
				}
			}
		}
		return nil
	})
}

func calcDestPath(s, hash, www string) string {
	// create the storage directories
	var err error
	if err = os.MkdirAll(www+"img", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	wwwimg := www + "img/"
	if err = os.MkdirAll(www+"gif", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	wwwgif := www + "gif/"
	if err = os.MkdirAll(www+"mp4", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	wwwmp4 := www + "mp4/"

	// ignore items if we are inside the www directory
	// if strings.HasPrefix(s, wwwimg) || strings.HasPrefix(s, wwwgif) || strings.HasPrefix(s, wwwmp4) {
	// 	return ""
	// }

	// path to move to
	parentdir := filepath.Base(filepath.Dir(s))
	var destpath string = ""
	ext := strings.ToLower(filepath.Ext(s))

	if IsFileType(s, RenamedJPG) {
		ext = ".jpg"
	}

	if IsFileType(s, SupportedImages) {
		destpath = filepath.Join(wwwimg, parentdir, strings.ToLower(hash + ext))
		fmt.Printf("i")
	} else if IsFileType(s, SupportedAnimations) {
		destpath = filepath.Join(wwwgif, parentdir, strings.ToLower(hash + ext))
		fmt.Printf("a")
	} else if IsFileType(s, SupportedVideos) {
		// get video info
		probe := ffprobe_video(s)
		duration := extract_duration(probe)
		// do not add video if longer than 30 secs
		if duration > 30 {
			// fmt.Println("\nToo big:", s)
			return ""
		}
		// if its a segmented video, do not rename
		if strings.HasPrefix(filepath.Base(s), "part.") {
			destpath = filepath.Join(wwwmp4, parentdir, strings.ToLower(filepath.Base(s)))
		} else {
			destpath = filepath.Join(wwwmp4, parentdir, strings.ToLower(hash + ext))
		}
		fmt.Printf("v")
	} else {
		// fmt.Println("\nIgnore:", s)
		return ""
	}

	return destpath
	
}

func addFiles(retChan chan int, root string, www string) {
	tot := 0
	num := 0
	
	// move all files to their www folder
	fmt.Println("Moving all files to", www, " ...")
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if !d.IsDir() {
			// if not supported, ignore
			if !IsFileType(s, SupportedAnimations) && !IsFileType(s, SupportedImages) && !IsFileType(s, SupportedVideos) {
				// log.Println("\nIgnore:", s)
				return nil
			}

			// calc file hash
			hash, err := chunk_md5(s, HASHLEN)
			if err != nil {
				return err
			}

			destpath := calcDestPath(s, hash, www)
			if destpath == "" {
				return nil
			}

			// create the parent dir
			if err = os.MkdirAll(filepath.Dir(destpath), os.ModePerm); err != nil {
				log.Fatal(err)
				return err
			}

			// if file is recorded in the db
			if _, exists := list[hash]; exists {
				// if it actually exists in www
				if fileExists(destpath) {
					fmt.Printf("=")
				} else {
					fmt.Printf("+")
					MoveFile(s, destpath)
				}
			} else {
				MoveFile(s, destpath)
			}

			// new file
			fi, err := os.Stat(destpath)
			var size float32 = -1.0
			if err == nil {
				// size in MB
				size = float32(fi.Size()) / 1024 / 1024
			}

			// add item
			list[hash] = FileItem{
				Path:   destpath,
				Desc:   "",
				SizeMB: size,
				URL:    "",
			}
			num += 1
		} else {
			if num > 0 {
				SaveItems()
			}
			fmt.Printf(" %6d\n", num)
			tot += num
			num = 0
			fmt.Println("\nListing", s, "...")
		}
		return nil
	})

	fmt.Println(" OK")
	SaveItems()
	retChan <- tot
}

func SaveItems() {
	// save as json file
	jsonString, _ := json.MarshalIndent(list, "", "   ")
	if err := ioutil.WriteFile(jsonfile, jsonString, 0644); err != nil {
		panic(err)
	}

	fmt.Println("Save:", jsonfile)

	// save as javascript file
	content := append([]byte("var items = "), jsonString...)
	content = append(content, []byte(";\n")...)
	if err := ioutil.WriteFile(jsfile, content, 0644); err != nil {
		panic(err)
	}

	fmt.Println("Save:", jsfile)
}

func Run(root, www string, walk bool) {
	now := time.Now()
	fmt.Println("\n\t   Shindook Media Indexer")
	fmt.Println("\t", now.Format(time.UnixDate), "\n")

	jsfile = filepath.Join(www, "list.js")
	jsonfile = filepath.Join(www, "list.json")

	fmt.Println("INBOX DIR:\t", root)
	fmt.Println("SAVE DIR:\t", www)
	fmt.Println("JS FILE:\t", jsfile)
	fmt.Println("JSON FILE:\t", jsonfile)

	fmt.Println("----------------------\n")
	start := time.Now()

	// check root exists
	if !DirExists(root) {
		log.Printf("Root directory does not exist.")
		return
	}

	// file list
	list = make(map[string]FileItem)

	// parse the database file
	if fileExists(jsonfile) {
		jsonString, _ := ioutil.ReadFile(jsonfile)
		err := json.Unmarshal(jsonString, &list)
		if err != nil {
			fmt.Println("JSON DB Parse Failed:", jsonfile, err)
			return
		}

		// make sure the added files are still okay
		for hash, item := range list {
			if !fileExists(item.Path) {
				delete(list, hash)
			}
		}
	}

	if walk {
		tot := 0
		totChan := make(chan int)
		var wg sync.WaitGroup
	
		// walk the directory, populate the file list
		CheckFfmpeg()
		fmt.Println("FFMPEG OK")
	
		wg.Add(1)
		// spawn bg processor for ffmpeg
		go func() {
			defer wg.Done()
			convertMP4(root, www)
			segmentVideo(root, www)
			fmt.Println("FFMPeg processes finished!")
		}()
	
		// start adding, wait until data on channel
		go addFiles(totChan, root, www)
		tot += <- totChan
	
		// wait for the ffmpeg process
		fmt.Println("Waiting for the ffmpeg processes to finish ... ")
		wg.Wait()
	
		// add the newly processed ffmpeg files
		// wait until data on channel
		go addFiles(totChan, root, www)
		tot += <- totChan
	
		elapsed := time.Since(start)
		fmt.Println("\n----------------------")
		fmt.Println("TOTAL:\t", tot)
		fmt.Println("SAVE:\t", jsfile)
		fmt.Println("UPDATE:\t", jsonfile)
		fmt.Println("ELAPSED:", elapsed)
	}
}
