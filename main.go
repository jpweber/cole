package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jpweber/cole/notifications"
)

const (
	version = "v0.1.0"
)

var (
	interval       *int
	source         *string
	message        *string
	remoteEndpoint *string
	method         *string
)

func main() {

	versionPtr := flag.Bool("v", false, "Version")
	interval = flag.Int("t", 60, "Time interval, in seconds, to wait before sending an alert \nif a ping is not received")
	source = flag.String("s", "", "name of the prometheus server we are watching")
	message = flag.String("b", "Did not recieve a deadman switch alert.", "Body of the notification")
	remoteEndpoint = flag.String("e", "", "URL of the endpoint to send messages to. Include the scheme http|https")
	method = flag.String("m", "POST", "HTTP method to use when talking to the remote endpoint. Default is POST")
	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()
	if *versionPtr == true {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Println("Starting application...")

	// TODO:
	// read from config file

	// create notification
	n := notifications.Notification{
		Source:         *source,
		Message:        *message,
		RemoteEndpoint: *remoteEndpoint,
		Method:         *method,
	}

	// init first timer at launch of service
	// TODO:
	// figure out a way to start another timer after this alert fires.
	// we want this to continue to go off as long as the dead man
	// switch is not being tripped.
	dmsTimer := time.AfterFunc(time.Duration(*interval)*time.Second, n.Alert)

	// HTTP Handlers
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Pong")
		// stop any existing timer channels
		dmsTimer.Stop()
		// start a new timer
		dmsTimer = time.AfterFunc(time.Duration(*interval)*time.Second, n.Alert)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	})

	// Server Lifecycle
	s := http.Server{Addr: ":8080"}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
