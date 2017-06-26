package main

import "net"
import "strings"
import "fmt"
import "errors"
import "io/ioutil"
import "time"
import "net/http"
import "os"

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

	ip, err := getIP()
	fmt.Println(ip)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Update(ip)

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
