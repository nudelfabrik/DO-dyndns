package Gandi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nudelfabrik/GOdyndns/settings"
)

type GandiClient struct {
	domain    string
	subdomain string
	recordID  int
	lastIP    string
	token     string
}

type record struct {
	TTL    int      `json:"rrset_ttl"`
	Values []string `json:"rrset_values"`
}

func NewGandiClient(setting *settings.Settings) (*GandiClient, error) {
	client := &GandiClient{}

	client.token = setting.Token
	client.domain = setting.Domain
	client.subdomain = setting.Subdomain
	setting.Token = ""

	return client, nil
}

func (c *GandiClient) Update(ip string) error {
	if c.lastIP == ip {
		// Record is up to date
		fmt.Println(time.Now().Format(time.RFC1123), " Record is up to date")
		return nil
	}
	c.lastIP = ip
	request := record{}
	request.TTL = 1800
	request.Values = []string{ip}

	data, err := json.Marshal(request)

	req, err := http.NewRequest("PUT", "https://dns.api.gandi.net/api/v5/domains/"+c.domain+"/records/"+c.subdomain+"/A", bytes.NewBuffer(data))
	req.Header.Add("X-Api-Key", c.token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 201 {
		return errors.New(resp.Status)
	}

	fmt.Println(time.Now().Format(time.RFC1123), " Changed IP: ", ip)
	return nil
}
