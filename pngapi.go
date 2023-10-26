package main

import (
	"io"
	"net/http"
	"strings"
)

func png_getMyIP() (string, error) {
	res, err := http.Get("https://pinig.in/utils/getMyIP?nojson")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	return string(data), err
}

func isIPv6(ipAddr string) bool {
	return strings.Contains(ipAddr, ":")
}
