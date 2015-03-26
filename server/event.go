package server

import (
	"strings"
	"time"
)

type Event struct {
	Server  CheckableServer
	Time    time.Time
	Type    string
	Latency time.Duration
}

func NewRedAlert(checkableServer CheckableServer, latency time.Duration) *Event {
	return &Event{Server: checkableServer, Time: time.Now(), Type: "redalert", Latency: latency}
}

func NewGreenAlert(checkableServer CheckableServer, latency time.Duration) *Event {
	return &Event{Server: checkableServer, Time: time.Now(), Type: "greenalert", Latency: latency}
}

func (e *Event) isRedAlert() bool {
	return e.Type == "redalert"
}

func (e *Event) isGreenAlert() bool {
	return e.Type == "greenalert"
}

func (e *Event) ShortMessage() string {

	if e.isRedAlert() {
		return strings.Join([]string{"Uhoh,", e.Server.GetServerDetails().Name, "not responding. Failed ping to", e.Server.GetServerDetails().Address}, " ")
	}

	if e.isGreenAlert() {
		return strings.Join([]string{"Woo-hoo,", e.Server.GetServerDetails().Name, "is now reachable. Successful ping to", e.Server.GetServerDetails().Address}, " ")
	}

	return ""
}

func (e *Event) PrintLatency() int64 {
	return e.Latency.Nanoseconds() / 1e6
}
