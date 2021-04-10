package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	domain   = "YOUR DOMAIN" //Your registered domain name eg. abcd.com
	hostname = "<YOUR SUBDOMAIN NAME>" // Use the Name given in A record 
	key      = "<YOUR GODADDY API KEY>" //API Key
	secret   = "<YOUR GODADDY SECRET"  //Secrete Text 
)

var AuthKey string
var CurrIP string //Current Domain IP

type GDdata struct {
	Data string `json:data`
}

func main() {
	AuthKey = fmt.Sprintf("sso-key %s:%s", key, secret)
	for {
		//Checks and updates IP every 10 minutes
		CheckandUpdateIP()
		time.Sleep(10 * time.Minute)
	}
}

func CheckandUpdateIP() {
	fmt.Println(time.Now(), "Checking IP...")
	ExtIP := GetExternalIP()
	if ExtIP != CurrIP {
		CurrIP = GetIP()
	}

	fmt.Println("Current Domain IP:", CurrIP,
		"Current External IP:", ExtIP)
	if CurrIP != ExtIP {
		UpdateIP(ExtIP)
	}
}

func GetIP() string {
	// Gets Current IP address stored in godaddy 'A' record
	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, hostname)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Authorization", AuthKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var bodyObj []GDdata
	json.Unmarshal(body, &bodyObj)

	return fmt.Sprintf("%v", bodyObj[0].Data)
}

func GetExternalIP() string {
	//Gets Current Public IP of the machine
	url := "https://api.ipify.org"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

func UpdateIP(NewIP string) {
	//Updates the New IP address to Godaddy 'A' Record
	url := fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/%s", domain, hostname)
	method := "PUT"

	payload := strings.NewReader(`[{"data":"` + NewIP + `"}]`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", AuthKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sucessfully Updated IP", NewIP, string(body))
}
