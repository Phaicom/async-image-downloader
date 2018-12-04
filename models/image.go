package models

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Image struct {
	AlbumID      int64  `json:"albumId"`
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

func createFile(imageURL string) (*os.File, error) {
	fileURL, err := url.Parse(imageURL)
	if err != nil {
		return nil, err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")

	fileName := segments[len(segments)-1]

	// Create directory
	// os.Mkdir("./images", os.ModePerm)
	// Create file relate to file name
	file, err := os.Create("./images/" + fileName + ".png")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Download image file from url
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Copy response file to target file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("Just Downloaded a file %v with size %d\n", file.Name, size)

	return file, nil
}
