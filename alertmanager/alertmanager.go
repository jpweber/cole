package alertmanager

import (
	"encoding/json"
	"net/http"
)

// AlertMessage is the message we receive from Alertmanager
type AlertMessage struct {
	Version           string            `json:"version"`
	GroupKey          int               `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

// Alert is a single alert.
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt,omitempty"`
	EndsAt      string            `json:"EndsAt,omitempty"`
}

func DecodeAlertMessage(r *http.Request) (AlertMessage, error) {
	body := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var alertMessage AlertMessage
	if err := body.Decode(&alertMessage); err != nil {
		return alertMessage, err
	}
	return alertMessage, nil
}
