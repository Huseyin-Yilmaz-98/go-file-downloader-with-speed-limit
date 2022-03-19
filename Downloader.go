package main

type Downloader interface {
	Download() error
}
