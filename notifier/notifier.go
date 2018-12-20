package notifier

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/jpweber/cole/configuration"
	"github.com/jpweber/cole/dmtimer"
	"github.com/jpweber/cole/slack"
	"github.com/prometheus/alertmanager/template"

	log "github.com/sirupsen/logrus"
)

// NotificationSet - the body of the alert message from cole
type NotificationSet struct {
	Message template.Data
	Config  configuration.Conf
	Timers  dmtimer.DmTimers
}

func (n *NotificationSet) constructBody() ([]byte, error) {
	jsonBody, err := json.Marshal(n)
	if err != nil {
		return jsonBody, err
	}
	return jsonBody, nil

}

func (n *NotificationSet) Alert() {
	log.Println("Sending Alert. Missed deadman switch notification.")
	// set up for future specific notification types
	// switch on n.Config.SenderType
	switch n.Config.SenderType {
	case "slack":
		n.slack()
	case "pagerduty":
		n.pagerDuty()
	default:
		// thinking I should just pass the whole alert message here
		// n.genericWebHook()
	}
	//

}

// genericWebHook - takes url as and expects json to be the payload
func (n *NotificationSet) genericWebHook(jsonBody []byte) {

	req, err := http.NewRequest(
		n.Config.HTTPMethod,
		n.Config.HTTPEndpoint,
		bytes.NewBuffer(jsonBody),
	)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error:", err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading reponse boddy", err)
	}

	log.Info(string(respData))

}

func (n *NotificationSet) slack() {
	// DEBUG
	log.Println("slack method")
	payload := slack.Payload{
		Text:      "Missed DeadManSwitch Alert  - " + n.Message.CommonAnnotations["message"],
		Username:  "Cole - DeadManSwitch Monitor",
		Channel:   "#general",
		IconEmoji: ":monkey_face:",
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error marshalling new data", err)
	}
	n.genericWebHook(jsonBody)
}

func (n *NotificationSet) pagerDuty() {

	pdPayload := pagerduty.V2Payload{
		Summary:  "Missed DeadManSwitch Check in.",
		Source:   "cole",
		Severity: n.Message.CommonLabels["severity"],

		Timestamp: time.Now().Format(time.RFC3339),
		Group:     n.Message.CommonLabels["job"],
		Class:     n.Message.CommonLabels["alertname"],
		Details:   n.Message.CommonAnnotations["message"],
	}
	event := pagerduty.V2Event{
		RoutingKey: n.Config.PDIntegrationKey,
		Action:     "trigger",
		DedupKey:   n.Message.CommonLabels["alertname"],
		Client:     "Cole - Dead Man Switch Monitor",
		Payload:    &pdPayload,
	}

	resp, err := pagerduty.ManageEvent(event)
	if err != nil {
		log.Errorf("Error Created Event in Pager Duty: %s", err)
		return
	}
	log.Printf("%+v", resp)

}
