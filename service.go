package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"redalert/alert"
	"redalert/common"
)

type Service struct {
	servers    []common.CheckableServer
	alerts     map[alert.AlertType]alert.Alert
	serviceLog *log.Logger
	wg         sync.WaitGroup
}

// Create channel to hold events to trigger alerts, buffering to one for testing
var eventsChannel chan *common.Event

func (s *Service) initialize() {
	eventsChannel = make(chan *common.Event)
	s.serviceLog = log.New(os.Stdout, "Service ", log.Ldate|log.Ltime)

	// Start event listening loop
	go func() {
		for {
			select {
			case newEvent := <-eventsChannel:
				s.serviceLog.Printf("Received event: %+v\n", newEvent.Server.GetServerWatcher())
				s.triggerAlert(newEvent)
			}
		}
	}()
}

func (s *Service) AddServer(checkableServer common.CheckableServer) {
	s.servers = append(s.servers, checkableServer)
}

func (s *Service) Start() {

	// use this to keep the service running, even if no monitoring is occuring
	s.wg.Add(1)

	for serverIndex, _ := range s.servers {
		go s.Monitor(s.servers[serverIndex])
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			s.wg.Done()
		}
	}()

}

func (s *Service) Monitor(checkableServer common.CheckableServer) {

	s.wg.Add(1)

	stopScheduler := make(chan bool)
	common.SchedulePing(checkableServer, eventsChannel, stopScheduler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		for _ = range sigChan {
			stopScheduler <- true
			s.wg.Done()
		}
	}()

	s.wg.Wait()

	s.wg.Done()

}

func (s *Service) triggerAlert(event *common.Event) {

	go func() {

		var err error
		for _, individualAlert := range s.alerts {
			// Convert event to alert package
			alertPackage := &alert.AlertPackage{
				Message:     event.ShortMessage(),
				AlertLogger: event.Server.GetServerWatcher().Log,
			}

			err = individualAlert.Trigger(alertPackage)
			if err != nil {
				s.serviceLog.Println(common.Red, "CRITICAL: Failure triggering alert ["+individualAlert.Name()+"]: ", err.Error())
			}
		}

	}()
}
