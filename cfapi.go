package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func cf_sendRequest(method, url string, body map[string]any) (map[string]any, error) {
	client := &http.Client{Timeout: 0}
	var data []byte
	var err error
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+cf_Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("Non-200 response")
	}
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var parsed map[string]any
	json.Unmarshal(data, &parsed)
	return parsed, nil
}

func cf_getZoneId() bool {
	response, err := cf_sendRequest("GET", "https://api.cloudflare.com/client/v4/zones", nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	cf_ZoneId = response["result"].([]any)[0].(map[string]any)["id"].(string)
	log.Println("Using zone", cf_ZoneId)
	return true
}

func cf_ClearPrevRecords() bool {
	response, err := cf_sendRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", cf_ZoneId), nil)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	for i := 0; i < len(response["result"].([]any)); i++ {
		if response["result"].([]any)[i].(map[string]any)["name"] == domain {
			_, err := cf_sendRequest("DELETE", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", cf_ZoneId, response["result"].([]any)[i].(map[string]any)["id"]), nil)
			if err != nil {
				log.Println("Failed to clear DNS record", response["result"].([]any)[i].(map[string]any)["id"])
			} else {
				log.Println("Removed record", response["result"].([]any)[i].(map[string]any)["id"])
			}
		}
	}
	return true
}

func cf_UpdateRecord() bool {
	method := "POST"
	urlAddition := ""
	if cf_RecordId != "" {
		method = "PUT"
		urlAddition = "/" + cf_RecordId
	}
	ipAddr, err := png_getMyIP()
	if err != nil {
		log.Println("Failed to get my IP:", err.Error())
		saveOutput("Failed to get my IP: " + err.Error())
		return false
	}
	if cf_RecordId == "" {
		lastIp = ipAddr
	}
	request := make(map[string]any)
	request["content"] = ipAddr
	request["name"] = domain
	request["proxied"] = false
	if isIPv6(ipAddr) {
		request["type"] = "AAAA"
	} else {
		request["type"] = "A"
	}
	request["comment"] = "Submitted by MPC dDNS at " + time.Now().Format("15:04:05")
	request["ttl"] = ttl
	if lastIp != ipAddr {
		cf_ClearPrevRecords()
		cf_RecordId = ""
		urlAddition = ""
		method = "POST"
	}
	response, err := cf_sendRequest(method, fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records%s", cf_ZoneId, urlAddition), request)
	if err != nil {
		log.Println("Failed to submit IP:", err.Error())
		saveOutput("Failed to submit IP: " + err.Error())
		return false
	}
	if cf_RecordId == "" {
		cf_RecordId = response["result"].(map[string]any)["id"].(string)
		log.Println("Record ID:", cf_RecordId)
	}
	lastIp = ipAddr
	output := fmt.Sprintf("IP %s submitted at %s", ipAddr, time.Now().Format("15:04:05"))
	fmt.Println(output)
	saveOutput(output)
	return true
}
