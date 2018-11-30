package notifier

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/jpweber/cole/dmtimer"

	"github.com/jpweber/cole/alertmanager"
	"github.com/jpweber/cole/configuration"
	"github.com/jpweber/cole/slack"

	log "github.com/sirupsen/logrus"
)

// NSets - type for holding notification sets
// type NSets map[string]NotificationSet

// NotificationSet - the body of the alert message from cole
type NotificationSet struct {
	Message alertmanager.AlertMessage
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
	// my personal slack for testing
	// TODO need to figure out how this is going to be passed to us.
	// n.RemoteEndpoint = "https://hooks.slack.com/services/..."
	n.Config.HTTPEndpoint = "https://hooks.slack.com/services/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"
	n.Config.HTTPMethod = "POST"
	payload := slack.Payload{
		Text:      n.Message.CommonAnnotations["summary"] + " - " + n.Message.CommonAnnotations["description"],
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
