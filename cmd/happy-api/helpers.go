package main

import (
	"fmt"
	"io"
	"math/rand"
	"path/filepath"

	"github.com/function61/gokit/crypto/cryptoutil"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/spf13/cobra"
)

func newEntry() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Generate ID for new file",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println(cryptoutil.RandBase64UrlWithoutLeadingDash(3))
		},
	}
}

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
		//nolint:gosimple
		if _, is := err.(exif.TagNotPresentError); is {
			return "", nil
		} else {
			return "", err
		}
	}

	return artist.StringVal()
}
