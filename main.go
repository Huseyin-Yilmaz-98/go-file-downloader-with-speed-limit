package main

import (
	"fmt"
)

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
