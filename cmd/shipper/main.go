package main

import (
	"log"
	"time"
)

const SecondsForRecheck = 60

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	for {
		log.Print("Testing :)")
		time.Sleep(time.Second * SecondsForRecheck)
	}
}
