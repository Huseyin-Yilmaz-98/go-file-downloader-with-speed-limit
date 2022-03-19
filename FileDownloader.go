package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FileDownloader struct {
	Downloader
	url                string
	totalSize          int64
	downloadedBytes    int64
	filePath           string
	timeTookToDownload time.Duration
	speedLimit         int64
}

const bufferSize = 16 * 1024

//Checks whether the http response has content-length headers and returns its value if found.
func parseContentLength(resp *http.Response) (int64, error) {
	var result int64
	contentLength := resp.Header.Get("content-length")
	if contentLength == "" {
		return result, errors.New("no content-length header found")
	}

	totalSize, err := strconv.ParseInt(contentLength, 10, 64)

	if err != nil {
		return result, errors.New("content-length not an integer")
	}

	return totalSize, nil
}

func (d *FileDownloader) printProgress() {
	fmt.Printf("\r %d/%d kb downloaded", int64(bytesToKB(d.downloadedBytes)), int64(bytesToKB(d.totalSize)))
	if d.totalSize > 0 {
		fmt.Printf(" (%d%% complete)", int64(float64(d.downloadedBytes)/float64(d.totalSize)*100))
	}
}

func (d *FileDownloader) Download() error {
	resp, err := http.Get(d.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	totalSize, err := parseContentLength(resp)
	if err == nil {
		fmt.Printf("Server reported total size as %d kb\n", int64(bytesToKB(totalSize)))
		d.totalSize = totalSize
	} else {
		fmt.Println(err)
		d.totalSize = 0
	}

	tempFilePath := d.filePath + ".tmp"
	file, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, bufferSize)       //byte array that will hold downloaded chunks
	speedLimit := kbToBytes(d.speedLimit) //convert speed limit to bytes for comparison
	idealBytesPerMilliSecond := int64(float64(speedLimit) / float64(time.Second.Milliseconds()))
	startTime := time.Now()
	lastStartMark := time.Now()
	quota := speedLimit

	fmt.Println("Starting to download " + d.filePath)
	for {
		chunkSize, err := resp.Body.Read(buf)
		if chunkSize > 0 {
			_, err := file.Write(buf[0:chunkSize])
			if err != nil {
				return err
			}
			quota -= int64(chunkSize)
			d.downloadedBytes += int64(chunkSize)

			elapsed := time.Since(lastStartMark)

			if elapsed >= time.Second || quota <= 0 {
				if quota <= 0 {
					overflowDuration := time.Duration(idealBytesPerMilliSecond * -quota)
					time.Sleep(time.Second - elapsed + overflowDuration)
				}

				d.printProgress()
				lastStartMark = time.Now()
				quota = speedLimit
			}
		}
		if err != nil {
			if err != io.EOF { //EOF means successful download
				return err
			}
			break
		}
	}

	d.printProgress()
	d.timeTookToDownload = time.Since(startTime)

	file.Close()

	err = os.Rename(tempFilePath, d.filePath)
	if err != nil {
		return err
	}

	return nil
}
