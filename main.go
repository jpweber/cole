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

	"github.com/jpweber/cole/dmtimer"
)

const (
	version = "v0.1.0"
)

var (
	ns = notifier.NotificationSet{}
)

func init() {
	// Log as text. Color with tty attached
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.WarnLevel)
}

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
	ns = notifier.NotificationSet{
		Config: conf,
		Timers: dmtimer.DmTimers{},
	}

	// HTTP Handlers
	http.HandleFunc("/ping", logger(ping))
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
