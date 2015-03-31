package common

import (
	"log"
	"time"
)

type Server struct {
	Name      string
	LastEvent *Event
}

type ServerDetails struct {
	Name     string
	Address  string
	Interval int
}

type ServerWatcher struct {
	Log       *log.Logger
	FailCount int
	LastEvent *Event
}

type CheckableServer interface {
	GetServerDetails() ServerDetails
	GetServerWatcher() *ServerWatcher
	IncrFailCount()
	Healthcheck() (time.Duration, error)
}

func SchedulePing(checkableServer CheckableServer, eventChan chan *Event, stopChan chan bool) {

	go func() {

		serverDetails := checkableServer.GetServerDetails()
		serverWatcher := checkableServer.GetServerWatcher()

		var err error
		var event *Event
		var latency time.Duration

		originalDelay := time.Second * time.Duration(serverDetails.Interval)
		delay := time.Second * time.Duration(serverDetails.Interval)

		for {

			latency, err = checkableServer.Healthcheck()

			if err != nil {

				serverWatcher.Log.Println(Red, "ERROR: ", err, Reset)

				// before sending an alert, pause 5 seconds & retry
				// prevent alerts from occaisional errors ('no such host' / 'i/o timeout') on cloud providers
				// todo: adjust sleep to fit with interval
				time.Sleep(5 * time.Second)
				_, rePingErr := checkableServer.Healthcheck()
				if rePingErr != nil {

					// re-ping fails (confirms error)

					event = NewRedAlert(checkableServer, latency)
					serverWatcher.LastEvent = event

					eventChan <- event

					checkableServer.IncrFailCount()
					if serverWatcher.FailCount > 0 {
						delay = time.Second * time.Duration(serverWatcher.FailCount*serverDetails.Interval)
					}

				} else {

					// re-ping succeeds (likely false positive)

					delay = originalDelay
					serverWatcher.FailCount = 0
				}

			} else {

				event = NewGreenAlert(checkableServer, latency)
				isRedalertRecovery := serverWatcher.LastEvent != nil && serverWatcher.LastEvent.isRedAlert()
				serverWatcher.LastEvent = event
				if isRedalertRecovery {
					serverWatcher.Log.Println(Green, "RECOVERY: ", Reset, serverDetails.Name)
					eventChan <- event
				}

				delay = originalDelay
				serverWatcher.FailCount = 0

			}

			select {
			case <-time.After(delay):
			case <-stopChan:
				return
			}
		}
	}()

}
