package main

import (
	"fmt"
	"net/http"
	"time"
)

type DoClient struct {
	domain    string
	subdomain string
	recordID  int
	lastIP    string
	token     string
}

type record struct {
	TTL    int      `json:rrset_ttl"`
	Values []string `json:rrset_values"`
}

func NewDoClient(setting *settings) (*DoClient, error) {
	doClient := &DoClient{}

	doClient.token = setting.Token
	doClient.domain = setting.Domain
	doClient.subdomain = setting.Subdomain
	setting.Token = ""

	return doClient, nil
}

func (c *DoClient) Update(ip string) error {
	if c.lastIP == ip {
		// Record is up to date
		fmt.Println(time.Now().Format(time.RFC1123), " Record is up to date")
		return nil
	}
	c.lastIP = ip
	request := record{}
	request.TTL = 1800
	request.Values = []string{ip}

	req, err := http.NewRequest("PUT", "https://dns.api.gandi.net/api/v5/domains/"+c.domain+"/records"+c.subdomain+"/A", nil)
	req.Header.Add("X-Api-Key", c.token)

	client := &http.Client{}

	_, err = client.Do(req)
	return err
}
