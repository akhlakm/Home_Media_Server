package indexer

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func CheckFfmpeg() {
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		log.Fatal(err)
	}

	_, err = exec.LookPath("ffmpeg")
	if err != nil {
		log.Fatal(err)
	}
}

//
// Get compact information of a video
//
func ffprobe_video(vid string) string {

	// ffprobe -v quiet -print_format compact -show_format file.mp4
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "compact",
		"-show_format", vid)

	// run
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("\nPROBE ERROR:", vid)
		return ""
	}

	return string(out)
}

func extract_duration(probe string) float32 {
	sections := strings.Split(probe, "|")
	if len(sections) >= 8 {
		duration := sections[7]
		ts := strings.Split(duration, "=")[1]
		t, _ := strconv.ParseFloat(ts, 32)
		return float32(t)
	}
	return -1.0
}

//
// Segment video
//
func ffmpeg_convert_mp4(vid string) bool {
	// output mp4
	fmtout := filepath.Join(filepath.Dir(vid), filepath.Base(vid)+".mp4")

	// previously converted
	if fileExists(fmtout) {
		return true
	}

	// ffmpeg segmentation command, abort on file exist
	cmd := exec.Command("ffmpeg",
		"-nostats", "-loglevel", "1",
		"-i", vid,
		fmtout)

	// run
	// fmt.Println(cmd.String())
	fmt.Printf("\nConverting %s ... ", vid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR\n", string(out))
		return false
	}
	fmt.Printf("OK %s\n", out)
	return true
}

//
// Segment video
//
func ffmpeg_segment(hash string, vid string, seglength int) bool {
	dir := filepath.Join(filepath.Dir(vid), "part."+hash)
	// output mp4 name format
	fmtout := dir + ".%02d.mp4"

	// ffmpeg segmentation command, abort on file exist
	cmd := exec.Command("ffmpeg",
		"-nostats", "-loglevel", "1",
		"-i", vid,
		"-f", "segment",
		"-segment_time", strconv.Itoa(seglength),
		"-c", "copy",
		"-reset_timestamps", "1",
		"-map", "0",
		"-n", fmtout)

	// run
	// fmt.Println(cmd.String())
	fmt.Printf("\nSegmenting %s ... ", vid)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR\n", string(out))
		return false
	}
	fmt.Printf("OK %s\n", out)
	return true
}

//
// Create poster for videos
//
func ffmpeg_create_poster(hash string, vid string, duration float32) []string {
	N := 5
	dir := filepath.Join(filepath.Dir(vid), "poster."+hash)

	posters := []string{}

	if duration < 1 {
		return posters
	}

	for n := 1; n <= N; n++ {
		t := int((float32(n) - 0.5) * duration)
		fmtout := dir + "." + strconv.Itoa(n) + ".jpg"

		// ffmpeg image command
		cmd := exec.Command("ffmpeg",
			"-nostats", "-loglevel", "1",
			"-ss", strconv.Itoa(t),
			"-i", vid,
			"-vf", "select='eq(pict_type\\,I)'",
			"-vframes", "1",
			"-n", fmtout)

		// run
		// fmt.Println(cmd.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error", vid, string(out))
		}
		fmt.Println(fmtout)
		fmt.Printf("%s\n", out)
		posters = append(posters, fmtout)
	}

	return posters
}
