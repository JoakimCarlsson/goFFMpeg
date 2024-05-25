package ffmpeg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ProgressData struct {
	Frame   int
	Fps     float64
	Q       int
	Size    int
	Time    time.Duration
	Bitrate float64
	Dup     int
	Drop    int
	Speed   float64
}

func parseProgressData(line string) (ProgressData, bool) {
	var patterns = map[string]*regexp.Regexp{
		"frame":   regexp.MustCompile(`frame=\s*(\d+)`),
		"fps":     regexp.MustCompile(`fps=\s*(\d+)`),
		"q":       regexp.MustCompile(`q=\s*([\d.]+)`),
		"size":    regexp.MustCompile(`size=\s*(\d+)\s*kB`),
		"time":    regexp.MustCompile(`time=\s*([\d:.]+)`),
		"bitrate": regexp.MustCompile(`bitrate=\s*([\d.]+)kbits/s`),
		"dup":     regexp.MustCompile(`dup=\s*(\d+)`),
		"drop":    regexp.MustCompile(`drop=\s*(\d+)`),
		"speed":   regexp.MustCompile(`speed=\s*([\d.]+)x`),
	}

	progressData := ProgressData{}
	found := false

	fmt.Println("Line: ", line)
	for key, pattern := range patterns {
		match := pattern.FindStringSubmatch(line)
		if len(match) > 1 {
			found = true
			switch key {
			case "frame":
				progressData.Frame = parseInt(match[1])
			case "fps":
				progressData.Fps = parseFloat(match[1])
			case "q":
				progressData.Q = parseInt(match[1])
			case "size":
				progressData.Size = parseInt(match[1])
			case "time":
				progressData.Time = parseDuration(match[1])
			case "bitrate":
				progressData.Bitrate = parseFloat(match[1])
			case "dup":
				progressData.Dup = parseInt(match[1])
			case "drop":
				progressData.Drop = parseInt(match[1])
			case "speed":
				progressData.Speed = parseFloat(match[1])
			}
		}
	}

	return progressData, found
}

func parseInt(value string) int {
	v, _ := strconv.Atoi(value)
	return v
}

func parseFloat(value string) float64 {
	v, _ := strconv.ParseFloat(value, 64)
	return v
}

func parseDuration(value string) time.Duration {
	v, _ := time.ParseDuration(strings.Replace(value, ":", "h", 1))
	return v
}
