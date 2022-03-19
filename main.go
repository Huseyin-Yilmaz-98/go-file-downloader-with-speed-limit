package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

//Parses url and speed limit from command line arguments and the file name from the url string.
func parseArgs() (string, int64, string) {
	var speedLimit int64 = 100000
	args := os.Args[1:]
	if len(args) == 0 {
		panic("No url passed!")
	} else if len(args) == 1 {
		fmt.Println("Setting speed limit to default value.")
	} else {
		parsedSpeedLimit, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			fmt.Println("Invalid speed limit input. Setting speed limit to default value.")
		} else {
			speedLimit = parsedSpeedLimit
		}
	}

	urlString := args[0]
	fileName := parseFileNameFromURLString(urlString)
	if fileName == "" {
		fmt.Println("Failed to parse file name from url. Giving random name.")
		fileName = fmt.Sprintf("unknown_%d", time.Now().Unix())
	}

	return urlString, speedLimit, fileName
}

func main() {
	urlString, speedLimit, fileName := parseArgs()
	d := FileDownloader{url: urlString,
		filePath:   fileName,
		speedLimit: speedLimit}

	err := d.Download()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\nDownloaded %d kb in %d seconds with an average speed of %d kb/s...",
			int64(bytesToKB(d.downloadedBytes)),
			int64(d.timeTookToDownload.Seconds()),
			int64(bytesToKB(d.downloadedBytes)/d.timeTookToDownload.Seconds()))
	}
}
