package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func server(c *DoClient, port string) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		update(c)
		w.WriteHeader(200)
	})

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
		// Added Timeouts to prevent resource exhaustion
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// Shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		httpServer.Close()
	}()

	httpServer.ListenAndServe()

}
