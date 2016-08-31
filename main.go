package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Not enough arguments.")
		return
	}

	path, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	subtitles, err := getSubtitles(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := convert(subtitles, path); err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(subtitles)
	// fmt.Println(converted)
}
