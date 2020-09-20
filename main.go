package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const ovhUpdateURL = "https://www.ovh.com/nic/update"

var host string
var username string
var password string
var params url.Values

var lastIP string

func init() {
	flag.StringVar(&host, "h", "", "host")
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.Parse()

	params = url.Values{}
	params.Add("system", "dyndns")
	params.Add("hostname", host)
}

func main() {
	log.Println("Starting DNS update...")

	t := time.NewTicker(5 * time.Minute)
	for {
		ip, err := retrievePublicIP()
		if err != nil {
			log.Println("An error occured retrieving public IP address")
		} else {
			if ip != lastIP {
				log.Println("Old IP: " + lastIP)
				log.Println("New IP: " + ip)

				updateDynHost(ip)
			}

			lastIP = ip
		}

		<-t.C
	}
}

func retrievePublicIP() (string, error) {
	log.Println("Retrieving public IP")
	const IPAPI = "https://api-ipv4.ip.sb/ip"

	resp, err := http.Get(IPAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}

	return strings.TrimSpace(string(jsonData)), nil
}

func updateDynHost(ip string) {
	updateURL, _ := url.Parse(ovhUpdateURL)
	params.Add("myip", ip)

	updateURL.RawQuery = params.Encode()

	log.Println(updateURL.String())
	client := &http.Client{}
	req, err := http.NewRequest("GET", updateURL.String(), nil)
	if err != nil {
		log.Fatal("Error creating request", err)
	}
	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}

	responseString := string(response)

	if resp.StatusCode != 200 {
		log.Println("Error updating the IP")
		log.Println(responseString)
		return
	}

	if strings.Contains(responseString, "good") {
		log.Println("New IP: " + ip)
	} else if strings.Contains(responseString, "nochg") {
		log.Println("Same IP, no change (" + ip + ")")
	}
	log.Println("--------------------------------------")
}
