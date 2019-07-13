package converter

import (
	"encoding/json"
	"fmt"
	"time"

	"git.circuitco.de/self/bellbox"
)

type Prometheus struct {
}

type Alert struct {
	Status       string
	Labels       map[string]string
	Annotations  map[string]string
	GeneratorURL string
	StartsAt     time.Time
	EndsAt       time.Time
}

type AlertmanagerJson struct {
	GroupKey     string
	Version      string
	Receiver     string
	Status       string
	ExternalURL  string
	CommonLabels map[string]string
	GroupLabels  map[string]string
	Alerts       []Alert
}

func (p Prometheus) Handle(input []byte) bellbox.Message {
	msg := bellbox.Message{}
	j := AlertmanagerJson{}
	err := json.Unmarshal(input, &j)
	if err != nil {
		return msg
	}
	pass := 0
	fail := 0
	summaryFmt := "service %s (%s) summary: %s "
	fullSummary := ""
	for _, alert := range j.Alerts {
		// set the title based on the number of + and - service notifications
		if alert.Status == "firing" {
			fail += 1
		} else {
			pass += 1
		}
		fullSummary += fmt.Sprintf(summaryFmt, alert.Labels["name"], alert.Status, alert.Labels["summary"])
	}
	if fail > 0 {
		msg.Title = fmt.Sprintf("%d services failing. ", fail)
	}
	if pass > 0 {
		msg.Title += fmt.Sprintf("%d services OK. ", pass)
	}
	msg.Message = fullSummary
	fmt.Printf("e: %s a: %+v\n", err, j)
	return msg
}
