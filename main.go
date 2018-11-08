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
)

func alert() {
	log.Println("I would send alert here")
}

const (
	version = "v0.1.0"
)

var interval int

func main() {

	versionPtr := flag.Bool("v", false, "Version")
	interval := flag.Int("t", 30, "Time interval, in seconds, to wait before sending an alert if a ping is not received")
	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()
	if *versionPtr == true {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Println("Starting application...")
	f := func() {
		alert()
	}

	dmsTimer := time.AfterFunc(time.Duration(*interval)*time.Second, f)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Pong")
		dmsTimer.Stop()
		dmsTimer = time.AfterFunc(time.Duration(*interval)*time.Second, f)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	})

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
