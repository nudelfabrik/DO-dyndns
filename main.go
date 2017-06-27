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
	ip4, _, err := getIP(true, false)
	if err != nil {
		return err
	}

	err = c.Update(ip4)
	return err
}

func getIP(v4, v6 bool) (string, string, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	req := func(url string) (string, error) {
		response, err := netClient.Get(url)
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
		return ip.String(), nil
	}

	var ip4, ip6 string
	var err error
	if v4 {
		ip4, err = req("http://ipv4.icanhazip.com")
		if err != nil {
			return "", "", err
		}
	}
	if v6 {
		ip6, err = req("http://ipv6.icanhazip.com")
		if err != nil {
			return "", "", err
		}
	}

	return ip4, ip6, nil

}
