package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Not enough arguments.")
		return
	}
	vPath := args[0]

	_ = convert(vPath)

	subtitles, err := getSubtitles(vPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(subtitles)
}
