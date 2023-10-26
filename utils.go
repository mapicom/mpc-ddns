package main

import (
	"log"
	"os"
)

func saveOutput(data string) bool {
	if outputPath == "" {
		return true
	}
	file, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println("Unable to write output:", err.Error())
		return false
	}
	defer file.Close()
	file.Write([]byte(data))
	return true
}
