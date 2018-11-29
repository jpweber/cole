package notifications

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jpweber/cole/alertmanager"

	"github.com/jpweber/cole/slack"
	log "github.com/sirupsen/logrus"
)

// Notification - the body of the alert message from cole
// source: would be the prometheus instance that we did not hear from
// message: whatever error message you want
// timestamp: time in the format of "2006-01-02 15:04:05"
// remoteEndpoint: URL of where to send the notification
// method: http method to use POST,GET,PUT etc.
type Notification struct {
	TimerID        string
	Message        alertmanager.AlertMessage
	Timestamp      time.Time
	RemoteEndpoint string
	Method         string
}

func (n *Notification) constructBody() ([]byte, error) {
	n.Timestamp = time.Now()
	jsonBody, err := json.Marshal(n)
	if err != nil {
		return jsonBody, err
	}
	return jsonBody, nil

}

func (n *Notification) Alert() {
	log.Println("Sending Alert. Missed deadman switch notification.")
	// set up for future specific notification types
	// right now only have a generic webhook
	// send a notification
	// n.genericWebHook()
	n.slack()
}

// genericWebHook - takes url as and expects json to be the payload
func (n *Notification) genericWebHook() {

	jsonBody, err := n.constructBody()
	if err != nil {
		log.Error("Error marshalling new data", err)
	}
	req, err := http.NewRequest(n.Method, n.RemoteEndpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error:", err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error readong reponse boddy", err)
	}

	log.Info(string(respData))

}

// TODO:
// temporary hack just to test the slack part before I figure
// out the rest of the proper integration
func genericWebHook(endpoint string, body slack.Payload) {

	jsonBody, err := json.Marshal(body)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error:", err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error readong reponse boddy", err)
	}

	log.Info(string(respData))

}

func (n *Notification) slack() {

	// my personal slack for testing
	// TODO need to figure out how this is going to be passed to us.
	webhookURL := "https://hooks.slack.com/services/..."

	payload := slack.Payload{
		Text:      n.Message.CommonAnnotations["summary"] + " - " + n.Message.CommonAnnotations["description"],
		Username:  "Cole - DeadManSwitch Monitor",
		Channel:   "#general",
		IconEmoji: ":monkey_face:",
	}

	genericWebHook(webhookURL, payload)
}
