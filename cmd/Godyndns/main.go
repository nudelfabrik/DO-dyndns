package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/nudelfabrik/GOdyndns"
	"github.com/nudelfabrik/GOdyndns/settings"
)

func main() {

	config := flag.String("f", "/usr/local/etc/DO-dyndns.json", "JSON Config file")
	pid := flag.String("p", "", "Pid File")

	flag.Parse()

	if *pid != "" {
		writePidFile(*pid)
	}

	setting, err := settings.LoadSettings(*config)
	if err != nil {
		fmt.Println(err)
		return
	}

	client, err := GOdyndns.CreateClient(setting)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = GOdyndns.Update(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	if setting.StartServer {
		GOdyndns.Server(client, setting.Port)
	}

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
