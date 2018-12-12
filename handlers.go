package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// wrapper function for http logging
func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf("%s - %s", r.Method, r.URL)
		fn(w, r)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	// start timer for http request duration metric
	timer := prometheus.NewTimer(httpDurations)
	defer timer.ObserveDuration()

	// timer testing
	// time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

	// init my error
	var err error
	if r.Method == "GET" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported", 405)
		return
	}

	defer r.Body.Close()
	data := template.Data{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error("Error decoding body ", err)
		http.Error(w, "Error decoding body", http.StatusBadRequest)
		return
	}
	ns.Message = data
	timerID := ns.Message.Alerts[0].Labels["alertname"]
	// DEBUG
	log.Println("timerID:", timerID)
	if err != nil {
		log.Println("Cannot register checkin", err)
	}

	// log metric of alert recieved
	dmAlertsRecieved.Inc()
	if ns.Timers.Get(timerID) != nil {
		// stop any existing timer channel
		ns.Timers.Get(timerID).Stop()
	}

	// start a new timer
	ns.Timers.Add(timerID, time.AfterFunc(time.Duration(ns.Config.Interval)*time.Second, ns.Alert))
	// DEBUG
	timerCount.Set(float64(ns.Timers.Len()))
	w.WriteHeader(200)
}
