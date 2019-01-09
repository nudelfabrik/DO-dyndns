package GOdyndns

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Client interface {
	Update(string) error
}

func Update(c Client) error {
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
