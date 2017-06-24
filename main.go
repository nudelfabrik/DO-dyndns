package main

import "github.com/digitalocean/godo"
import "golang.org/x/oauth2"
import "net"
import "strings"
import "context"
import "os"
import "fmt"
import "errors"
import "io/ioutil"
import "time"
import "net/http"

var DoClient *godo.Client
var noSubdomain = errors.New("subdomain not found")

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	initDoClient()
	domain := os.Getenv("DOMAIN")
	subdomain := os.Getenv("SUBDOMAIN")

	err := checkDomain(domain)
	if err != nil {
		fmt.Println("Domain not found: ", err)
		return
	}

	request := &godo.DomainRecordEditRequest{}
	request.Name = subdomain
	request.Type = "A"
	ip, err := getIP()
	fmt.Println(ip)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Data = ip

	id, err := checkSubdomain(domain, subdomain)
	fmt.Println(id)
	fmt.Println(err)
	if err == noSubdomain {
		err = createEntry(domain, subdomain, request)
		if err != nil {
			fmt.Println(err)
		}
	} else if err != nil {
		fmt.Println(err)
		return
	} else {
		updateEntry(domain, id, request)
		if err != nil {
			fmt.Println(err)
		}
	}

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

func updateEntry(domain string, id int, req *godo.DomainRecordEditRequest) error {
	_, _, err := DoClient.Domains.EditRecord(context.Background(), domain, id, req)
	return err
}

func createEntry(domain, subdomain string, req *godo.DomainRecordEditRequest) error {
	_, _, err := DoClient.Domains.CreateRecord(context.Background(), domain, req)
	return err
}

func checkSubdomain(domain, subdomain string) (int, error) {
	opt := &godo.ListOptions{}
	for {
		records, resp, err := DoClient.Domains.Records(context.Background(), domain, opt)
		if err != nil {
			return 0, err
		}

		// append the current page's droplets to our list
		for _, r := range records {
			if r.Name == subdomain {
				return r.ID, nil
			}
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return 0, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return 0, noSubdomain

}

func checkDomain(domain string) error {
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := DoClient.Domains.List(context.Background(), opt)
		if err != nil {
			return err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
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

func initDoClient() {

	pat := os.Getenv("TOKEN")

	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	DoClient = godo.NewClient(oauthClient)
}
