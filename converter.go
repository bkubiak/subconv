package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dwbuiten/go-mediainfo/mediainfo"
)

var mpl2Regex *regexp.Regexp = regexp.MustCompile(`^\[(\d+)\]\[(\d+)\](.+)`)

func convert(subtitles string, vPath string) string {
	var converted string
	format := format(subtitles)
	fmt.Printf("format: %s\n", format)
	switch format {
	case "microDVD":
		converted = microDVD(subtitles, vPath)
	case "mpl2":
		converted = mpl2(subtitles)
	}

	return converted
}

func format(subtitles string) string {
	var format string

	if mpl2Regex.MatchString(subtitles) {
		format = "mpl2"
	}

	return format
}

func mpl2(subtitles string) string {
	var result string

	getTime := func(dsecStr string) (string, error) {
		dsec, err := strconv.ParseUint(dsecStr, 10, 64)
		if err != nil {
			return "", err
		}

		t := time.Time{}
		displayTime := t.Add(time.Duration(dsec*100) * time.Millisecond).Format("15:04:05.000")
		displayTime = strings.Replace(displayTime, ".", ",", 1)

		return displayTime, nil
	}

	lines := strings.Split(subtitles, "\n")

	for i, line := range lines {
		matches := mpl2Regex.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		startTime, err := getTime(matches[1])
		if err != nil {
			continue
		}

		stopTime, err := getTime(matches[2])
		if err != nil {
			continue
		}

		text := strings.Replace(matches[3], "\r", "", -1)
		text = strings.Replace(text, "|", "\n", -1)

		result += fmt.Sprintf("%d \n%s --> %s \n%s \n\n", i+1, startTime, stopTime, text)
	}

	return result
}

func microDVD(subtitles string, vPath string) string {
	fps, err := getFPS(vPath)
	if err != nil {
		fmt.Println("FPS rate: 23.976 (default)")
		fps = 23.976
	}
	fmt.Printf("FPS rate: %.3f\n", fps)

	return ""
}

func getFPS(path string) (float64, error) {
	mediainfo.Init()
	info, err := mediainfo.Open(path)
	if err != nil {
		return 0, err
	}
	defer info.Close()

	val, err := info.Get("FrameRate", 0, mediainfo.Video)
	if err != nil {
		return 0, err
	}

	fps, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	return fps, nil
}
