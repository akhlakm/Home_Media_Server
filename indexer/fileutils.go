package indexer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

var SupportedImages = []string{
	".png", ".jpg", ".jpeg", ".jiff", ".jif",
}

var RenamedJPG = []string{
	".jpeg", ".jiff", ".jif",
}

var SupportedAnimations = []string{
	".gif", ".webp",
}

var SupportedVideos = []string{
	".mp4",
}

var Videos2Convert = []string{
	".ts", ".wmv", ".mkv", ".avi", ".webm", ".mov",
}

func IsFileType(s string, supported []string) bool {
	for _, ext := range supported {
		s = strings.ToLower(s)
		if strings.HasSuffix(s, ext) {
			return true
		}
	}
	return false
}

//
// Return first 5 characters md5 hash of the first 1 KB of a file
//
func chunk_md5(filePath string, characters int) (string, error) {
	chunksize := 1024 * 1024 // 1 Mb
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Read 1 KB
	buf := make([]byte, chunksize)

	bytesRead, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return returnMD5String, err
	}

	hash.Write(buf[:bytesRead])
	// fmt.Println("Bytes read:", bytesRead)

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:characters]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func MoveFile(src, dst string) error {
	if src == dst {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("chmod error: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}
	return nil
}
