package main

import (
	"fmt"
	"time"
)

func main() {
	loadConfig()
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
