package main

import (
	domainapp "go-ffmpeg-progress/src/internal/app/domain"
	"log"
)

func main() {
	dom, err := domainapp.New()
	if err != nil {
		log.Fatal(err)
	}

	err = dom.Run()
	if err != nil {
		log.Fatal(err)
	}
}
