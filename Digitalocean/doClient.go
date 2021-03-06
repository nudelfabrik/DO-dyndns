package Digitalocean

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/digitalocean/godo"
	"github.com/nudelfabrik/GOdyndns/settings"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type DoClient struct {
	domain    string
	subdomain string
	recordID  int
	lastIP    string
	client    *godo.Client
}

func NewDoClient(setting *settings.Settings) (*DoClient, error) {
	doClient := &DoClient{}

	pat := setting.Token
	setting.Token = ""

	tokenSource := &tokenSource{
		AccessToken: pat,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	doClient.client = godo.NewClient(oauthClient)

	err := doClient.setDomain(setting.Domain)
	if err != nil {
		return nil, err
	}
	err = doClient.setSubdomain(setting.Subdomain)
	if err != nil {
		return nil, err
	}
	return doClient, err
}

func (c *DoClient) Update(ip string) error {
	if c.lastIP == ip {
		// Record is up to date
		fmt.Println(time.Now().Format(time.RFC1123), " Record is up to date")
		return nil
	}
	request := &godo.DomainRecordEditRequest{}
	request.Name = c.subdomain
	request.Type = "A"
	request.Data = ip

	var err error
	var rec *godo.DomainRecord
	if c.recordID == 0 {
		rec, _, err = c.client.Domains.CreateRecord(context.Background(), c.domain, request)
	} else {
		rec, _, err = c.client.Domains.EditRecord(context.Background(), c.domain, c.recordID, request)
	}
	if rec != nil {
		c.lastIP = rec.Data
		fmt.Println(time.Now().Format(time.RFC1123), " Changed IP: ", ip)
	}
	return err
}

func (c *DoClient) setDomain(domain string) error {
	c.domain = domain
	opt := &godo.ListOptions{}
	for {
		domains, resp, err := c.client.Domains.List(context.Background(), opt)
		if err != nil {
			return err
		}

		for _, d := range domains {
			if d.Name == domain {
				return nil
			}
		}
		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return errors.New("Domain not found: " + domain)
}

func (c *DoClient) setSubdomain(subdomain string) error {
	c.subdomain = subdomain
	opt := &godo.ListOptions{}
	for {
		records, resp, err := c.client.Domains.Records(context.Background(), c.domain, opt)
		if err != nil {
			return err
		}

		for _, r := range records {
			if r.Name == subdomain {
				c.recordID = r.ID
				c.lastIP = r.Data
				return nil
			}
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return nil

}
