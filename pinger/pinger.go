package pinger

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"redalert/common"
	"redalert/server"
)

type Pinger struct {
	server.ServerDetails
	server.ServerWatcher
}

var pingerClient = http.Client{
	Timeout: time.Duration(10 * time.Second),
}

var Online []Pinger

// Setup servers to be collected by main
func init() {
	config, err := ReadConfigFile()
	if err != nil {
		panic(fmt.Sprintf("Missing or invalid config %v", err))
	}

	for _, individualPinger := range config.TargetServers {
		Online = append(Online, Pinger{
			ServerDetails: server.ServerDetails{
				Name:     individualPinger.Name,
				Address:  individualPinger.Address,
				Interval: individualPinger.Interval,
			},
			ServerWatcher: server.ServerWatcher{
				Log: log.New(os.Stdout, individualPinger.Name+" ", log.Ldate|log.Ltime),
			},
		})
	}
}

func (p Pinger) GetServerDetails() server.ServerDetails {
	return p.ServerDetails
}

func (p Pinger) GetServerWatcher() server.ServerWatcher {
	return p.ServerWatcher
}

func (p Pinger) Healthcheck() (time.Duration, error) {
	startTime := time.Now()
	p.GetServerWatcher().Log.Println("Pinging: ", p.GetServerDetails().Name)

	req, err := http.NewRequest("GET", p.GetServerDetails().Address, nil)
	if err != nil {
		return 0, errors.New("redalert ping: failed parsing url in http.NewRequest " + err.Error())
	}

	req.Header.Add("User-Agent", "Redalert/1.0")
	resp, err := pingerClient.Do(req)

	endTime := time.Now()
	latency := endTime.Sub(startTime)
	p.GetServerWatcher().Log.Println(common.White, "Analytics: ", latency, common.Reset)

	if resp != nil {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
	if err != nil {
		return latency, errors.New("redalert ping: failed client.Do " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return latency, errors.New("redalert ping: non-200 status code. status code was " + strconv.Itoa(resp.StatusCode))
	}
	p.GetServerWatcher().Log.Println(common.Green, "OK", common.Reset, p.GetServerDetails().Name)

	return latency, nil
}
