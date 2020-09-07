package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

// ParseHTML is parse html format
func ParseHTML(url string) string {
	resp, err := http.Get(url)
	CheckError(err)
	body, err := ioutil.ReadAll(resp.Body)
	CheckError(err)
	resp.Body.Close()
	return string(body)
}

// UpdateLog is auto update status to log
func UpdateLog(wrStatus string) {
	file, err := os.OpenFile("DDNS.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)
	defer file.Close()
	log.SetOutput(file)
	log.Println(wrStatus)
}

// UpdateDDNS is call google ddns service
func UpdateDDNS(mode int) {
	ddnsURL := "https://vDcK9v1XLjvVKA2Q:iumEDItXAksy8oPw@domains.google.com/nic/update?hostname=www.chliu.dev"
	context := ParseHTML(ddnsURL)
	status, _ := regexp.Compile("[a-z0-9A-Z.]{1,16}")
	analysisCode := status.FindAllString(context, -1)
	if analysisCode[0] == "good" {
		UpdateLog("Setup DDNS is sucessful")
		data := []byte(analysisCode[1])
		err := ioutil.WriteFile("./IP.txt", data, 0644)
		CheckError(err)
	} else if analysisCode[0] == "nochg" && mode == 1 {
		UpdateLog("Resetup DDNS is sucessful")
		data := []byte(analysisCode[1])
		err := ioutil.WriteFile("./IP.txt", data, 0644)
		CheckError(err)
	}
}

func main() {
	context := ParseHTML("https://www.myip.com/")
	ipRules, _ := regexp.Compile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
	ipArr := ipRules.FindAllString(context, 1)
	nowIP := ipArr[0]
	file, _ := os.Open("./IP.txt")
	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		UpdateDDNS(1)
	}
	logFile, err := os.OpenFile("DDNS.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	CheckError(err)
	logInfo, err := logFile.Stat()
	CheckError(err)
	if logInfo.Size() > 5000 {
		os.Remove("DDNS.log")
	}
	CheckError(err)
	storedIP := string(fileBody)
	if nowIP != storedIP {
		UpdateDDNS(0)
	} else {
		UpdateLog("IP isn't changed")
	}
}

// CheckError is check error whatever happen
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
