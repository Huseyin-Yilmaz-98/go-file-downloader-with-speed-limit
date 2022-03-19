package main

import (
	"fmt"
	"net/url"
	"strings"
)

func bytesToKB(bytes int64) float64 {
	return float64(bytes) / 1024
}

func kbToBytes(kb int64) int64 {
	return kb * 1024
}

func parseFileNameFromURLString(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	path := parsedURL.Path
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
