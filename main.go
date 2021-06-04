package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type dns struct {
	URLConfig []urlSet `json:"URLConfig"`
}

type urlSet struct {
	URL      string `json:"URL"`
	UserName string `json:"UserName"`
	UserPwd  string `json:"UserPwd"`
}

func parseJSON() []string {
	var DNS dns
	urlSlice := make([]string, 0)
	jsonFile, _ := os.Open("setting.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &DNS)
	//fmt.Println(DNS)
	for _, config := range DNS.URLConfig {
		url := "https://" + config.UserName + ":" + config.UserPwd + "@domains.google.com/nic/update?hostname=" + config.URL
		urlSlice = append(urlSlice, url)
	}
	return urlSlice
}

// UpdateDDNS is call google ddns service
func UpdateDDNS() {
	logFile, err := os.OpenFile("./Record", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)
	defer logFile.Close()
	log.SetOutput((logFile))
	for _, url := range parseJSON() {
		cmd := exec.Command("curl", url)
		out, err := cmd.Output()
		CheckError(err)
		log.Println(string(out))
	}
}

func main() {
	UpdateDDNS()
}

// CheckError is check error whatever happen
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
