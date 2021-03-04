package main

import (
	"io"
	"math/rand"
	"path/filepath"

	"github.com/rwcarlsen/goexif/exif"
)

// 10e239c4167f.jpg => 10e239c4167f
func fileIdFromFilename(filename string) string {
	return filename[0 : len(filename)-len(filepath.Ext(filename))]
}

func randBetween(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func findAttributionFromExifArtist(id string) (string, error) {
	f, err := images.Open("images/" + id + ".jpg")
	if err != nil {
		return "", err
	}
	defer f.Close()

	exifData, err := exif.Decode(f)
	if err != nil {
		if err == io.EOF { // no exif whatsoever
			return "", nil
		} else {
			return "", err
		}
	}

	artist, err := exifData.Get(exif.Artist)
	if err != nil {
		if _, is := err.(exif.TagNotPresentError); is {
			return "", nil
		} else {
			return "", err
		}
	}

	return artist.StringVal()
}
