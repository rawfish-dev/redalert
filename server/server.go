package server

import (
	"log"
	"sync"
	"time"

	"redalert/common"
)

type CheckableServer interface {
	GetServerDetails() ServerDetails
	GetServerWatcher() ServerWatcher
	Healthcheck() (time.Duration, error)
}

type ServerDetails struct {
	Name     string
	Address  string
	Interval int
	wg       sync.WaitGroup
}

type ServerWatcher struct {
	Log       *log.Logger
	failCount int
	LastEvent *Event
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

				serverWatcher.Log.Println(common.Red, "ERROR: ", err, common.Reset)

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

					serverWatcher.IncrFailCount()
					if serverWatcher.failCount > 0 {
						delay = time.Second * time.Duration(serverWatcher.failCount*serverDetails.Interval)
					}

				} else {

					// re-ping succeeds (likely false positive)

					delay = originalDelay
					serverWatcher.failCount = 0
				}

			} else {

				event = NewGreenAlert(checkableServer, latency)
				isRedalertRecovery := serverWatcher.LastEvent != nil && serverWatcher.LastEvent.isRedAlert()
				serverWatcher.LastEvent = event
				if isRedalertRecovery {
					serverWatcher.Log.Println(common.Green, "RECOVERY: ", common.Reset, serverDetails.Name)
					eventChan <- event
				}

				delay = originalDelay
				serverWatcher.failCount = 0

			}

			select {
			case <-time.After(delay):
			case <-stopChan:
				return
			}
		}
	}()

}

func (s *ServerWatcher) IncrFailCount() {
	s.failCount++
}
