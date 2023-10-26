package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s CONFIG.txt [OUTPUT.txt]\n", os.Args[0])
		return
	}
	loadConfig(os.Args[1])
	if len(os.Args) == 3 {
		outputPath = os.Args[2]
		log.Println("Output information will be write to", os.Args[2])
	}
	for cf_getZoneId() != true {
		fmt.Println("Reattemping after 8 seconds...")
		time.Sleep(time.Second * 8)
	}
	for cf_ClearPrevRecords() != true {
		fmt.Println("Reattemping after 8 seconds...")
		time.Sleep(time.Second * 8)
	}

	// First setup
	cf_UpdateRecord()

	// Long loop..
	for {
		time.Sleep(time.Second * time.Duration(updateInverval))
		cf_UpdateRecord()
	}
}
