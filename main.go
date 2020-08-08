package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func UpdateDDNS() {
	ddnsURL := "https://a6ZaoSCdhV0Zg3kn:EaJJJgeG4UnQwtuM@domains.google.com/nic/update?hostname=blog.chliu.dev"
	resp, err := http.Get(ddnsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	context := string(body)
	status, _ := regexp.Compile("[a-z0-9A-Z.]{1,16}")
	analysisCode := status.FindAllString(context, -1)
	if analysisCode[0] == "good" {
		data := []byte(analysisCode[1])
		err := ioutil.WriteFile("./IP.txt", data, 0644)
		if err != nil {
			log.Panic(err)
		}
	}
}

func main() {
	resp, err := http.Get("https://www.myip.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	context := string(body)
	ipRules, _ := regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
	ipArr := ipRules.FindAllString(context, 1)
	ip := ipArr[0]
	file, _ := os.Open("./IP.txt")
	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	storedIP := string(fileBody)
	fmt.Println(ip, storedIP)
	if ip != storedIP {
		UpdateDDNS()
	}
}
