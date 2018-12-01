package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jpweber/cole/configuration"
	"github.com/jpweber/cole/notifier"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/jpweber/cole/alertmanager"
	"github.com/jpweber/cole/dmtimer"
)

const (
	version = "v0.1.0"
)

func main() {

	versionPtr := flag.Bool("v", false, "Version")
	configFile := flag.String("c", "example.toml", "Path to Configuration File")

	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()
	if *versionPtr == true {
		fmt.Println(version)
		os.Exit(0)
	}

	log.Println("Starting application...")

	// read from config file
	conf := configuration.ReadConfig(*configFile)

	// init first timer at launch of service
	// TODO:
	// figure out a way to start another timer after this alert fires.
	// we want this to continue to go off as long as the dead man
	// switch is not being tripped.

	// init notificaiton set
	ns := notifier.NotificationSet{
		Config: conf,
		Timers: dmtimer.DmTimers{},
	}

	// HTTP Handlers
	http.HandleFunc("/ping/", func(w http.ResponseWriter, r *http.Request) {
		// init my error
		var err error

		if r.Method != "POST" {
			http.Error(w, "Only POST method is supported", 405)
			return
		}
		log.Info("Pong")
		ns.Message, err = alertmanager.DecodeAlertMessage(r)
		if err != nil {
			log.Error(err)
			return
		}
		timerID := ns.Message.GroupLabels["alertname"]
		log.Info(timerID)
		if err != nil {
			log.Println("Cannot register checkin", err)
		}

		if ns.Timers.Get(timerID) != nil {
			// stop any existing timer channel
			ns.Timers.Get(timerID).Stop()
		}

		// start a new timer
		ns.Timers.Add(timerID, time.AfterFunc(time.Duration(ns.Config.Interval)*time.Second, ns.Alert))
		w.WriteHeader(200)

	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, version)
	})

	http.Handle("/metrics", promhttp.Handler())

	// Server Lifecycle
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Info("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
