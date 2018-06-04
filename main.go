package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	customPath := ""
	if len(os.Args) == 2 {
		customPath = os.Args[1]
	}
	setting, err := loadSettings(customPath)
	if err != nil {
		fmt.Println(err)
	}
	client, err := NewDoClient(setting)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = update(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	if setting.StartServer {
		server(client, setting.Port)
	}

}

func update(c *DoClient) error {
	ip, err := getIP()
	if err != nil {
		return err
	}

	err = c.Update(ip)
	return err
}

func getIP() (string, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Get("http://ipv4.icanhazip.com")
	if err != nil {
		return "", err
	}
	responseText, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return "", err
	}

	str := string(responseText)
	str = strings.TrimSpace(str)
	ip := net.ParseIP(str)
	if ip == nil {
		return "", errors.New("Cannot Parse IP: " + str)
	}

	return ip.String(), err
}
