package searcher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"redalert/common"
)

type LCTGSearcher struct {
	common.ServerDetails
	*common.ServerWatcher
	Config *LCTGSearchConfig
}

var lctgSearchClient = http.Client{
	Timeout: time.Duration(30 * time.Second),
}

func NewLCTGSearcher() *LCTGSearcher {
	config, err := ReadConfigFile()
	if err != nil {
		panic(fmt.Sprintf("Missing or invalid config %v", err))
	}

	return &LCTGSearcher{
		ServerDetails: common.ServerDetails{
			Name:     config.Name,
			Address:  config.SearchPath,
			Interval: config.Interval,
		},
		ServerWatcher: &common.ServerWatcher{
			Log: log.New(os.Stdout, config.Name+" ", log.Ldate|log.Ltime),
		},
		Config: config,
	}
}

func (l LCTGSearcher) GetServerDetails() common.ServerDetails {
	return l.ServerDetails
}

func (l LCTGSearcher) GetServerWatcher() *common.ServerWatcher {
	return l.ServerWatcher
}

func (l LCTGSearcher) IncrFailCount() {
	l.ServerWatcher.FailCount++
}

func (l *LCTGSearcher) Healthcheck() (time.Duration, error) {

	startTime := time.Now()
	l.GetServerWatcher().Log.Println("Searching: ", l.GetServerDetails().Name)

	searchObject := l.generateSearchParams()

	// Build request from template
	requestBody := buildRequestFromTemplate(searchObject)

	request, err := http.NewRequest("POST", l.ServerDetails.Address, bytes.NewBufferString(requestBody.String()))
	if err != nil {
		return 0, errors.New("redalert search: failed parsing url in http.NewRequest " + err.Error())
	}

	// panic(fmt.Sprintf("LCTG Request: %+v", requestBody.String()))

	// Add appropriate headers
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := lctgSearchClient.Do(request)

	endTime := time.Now()
	latency := endTime.Sub(startTime)
	l.GetServerWatcher().Log.Println(common.White, "Analytics: ", latency, common.Reset)

	if err != nil {
		return latency, errors.New("redalert search: failed client.Do " + err.Error())
	}
	defer response.Body.Close()

	// Ensure no http errors were returned
	if response.StatusCode != http.StatusOK {
		return latency, errors.New("redalert search: non-200 status code. status code was " + strconv.Itoa(response.StatusCode))
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return latency, errors.New("redalert search: unable to read search response " + err.Error())
	}

	// Validate response body is well formed

	l.GetServerWatcher().Log.Printf("Status: %v", response.StatusCode)
	l.GetServerWatcher().Log.Printf("Response: %v", string(responseBody))

	l.GetServerWatcher().Log.Println(common.Green, "OK", common.Reset, l.GetServerDetails().Name)

	// panic(fmt.Sprintf("LCTG Response: %+v", string(responseBody)))

	return latency, nil
}

func buildRequestFromTemplate(searchObject *LCTGSearchObject) (parsedTemplateBytes bytes.Buffer) {
	parsedTemplate, err := template.ParseFiles(PACKAGE_PATH + "/request.tmpl")
	if err != nil {
		panic(err.Error())
	}

	err = parsedTemplate.Execute(&parsedTemplateBytes, searchObject)
	if err != nil {
		panic(err.Error())
	}

	return parsedTemplateBytes
}
