package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadConfig() bool {
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	parsed := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")
	cf_Token = parsed[0]
	domain = parsed[1]
	ttl, err = strconv.Atoi(parsed[2])
	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	if ttl < 60 {
		log.Fatalln("TTL must be 60 or higher!")
		return false
	}
	updateInverval, err = strconv.Atoi(parsed[3])
	if err != nil {
		log.Fatalln(err.Error())
		return false
	}
	return true
}
