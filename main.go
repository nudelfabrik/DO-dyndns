package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {

	config := flag.String("f", "/usr/local/etc/DO-dyndns.json", "JSON Config file")
	pid := flag.String("p", "", "Pid File")

	flag.Parse()

	if *pid != "" {
		writePidFile(*pid)
	}

	setting, err := loadSettings(*config)
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

// Write a pid file, but first make sure it doesn't exist with a running pid.
func writePidFile(pidFile string) error {
	// Read in the pid file as a slice of bytes.
	if piddata, err := ioutil.ReadFile(pidFile); err == nil {
		// Convert the file contents to an integer.
		if pid, err := strconv.Atoi(string(piddata)); err == nil {
			// Look for the pid in the process list.
			if process, err := os.FindProcess(pid); err == nil {
				// Send the process a signal zero kill.
				if err := process.Signal(syscall.Signal(0)); err == nil {
					// We only get an error if the pid isn't running, or it's not ours.
					return fmt.Errorf("pid already running: %d", pid)
				}
			}
		}
	}
	// If we get here, then the pidfile didn't exist,
	// or the pid in it doesn't belong to the user running this app.
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0664)
}
