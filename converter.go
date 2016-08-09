package main

import (
	"fmt"
	"strconv"

	"github.com/dwbuiten/go-mediainfo/mediainfo"
)

func getFPS(path string) (*float64, error) {
	mediainfo.Init()
	info, err := mediainfo.Open(path)
	if err != nil {
		return nil, err
	}
	defer info.Close()

	val, err := info.Get("FrameRate", 0, mediainfo.Video)
	if err != nil {
		return nil, err
	}

	fps, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}

	return &fps, nil
}

func convert(vPath string) string {
	var fps *float64
	fps, err := getFPS(vPath)
	if err != nil {
		fmt.Println("FPS rate: 23.976 (default)")
		*fps = 23.976
	}
	fmt.Printf("FPS rate: %.3f\n", *fps)

	return ""
}