package searcher

import (
	"encoding/xml"
	"time"

	"redalert/common"
	"redalert/model"
	"redalert/zumata"
)

const (
	PACKAGE_PATH string = "searcher/lctg"
)

type LCTGSearchObject struct {
	Login       string
	Password    string
	CheckInDate string
	RegionId    string
}

type LCTGSearchResults struct {
	XMLName         xml.Name             `xml:"SearchResponse"`
	ReturnStatus    LCTGStatus           `xml:"ReturnStatus"`
	PropertyResults []LCTGPropertyResult `xml:"PropertyResults>PropertyResult"`
	Latency         time.Duration
	// SearchURL       string                 `xml:"SearchURL"`
}

type LCTGStatus struct {
	Success   bool   `xml:"Success"`
	Exception string `xml:"Exception"`
}

type LCTGPropertyResult struct {
	PropertyID int `xml:"PropertyID"`
}

func (l LCTGSearcher) generateSearchParams() *LCTGSearchObject {
	// Generate a date a month later from now
	oneMonthLater := time.Now().Add(time.Duration(30 * 24 * time.Hour))
	oneMonthLaterString := oneMonthLater.Format(common.YYYY_MM_DD_FORMAT)

	// Generate a random ZUMATA mapped location
	randomLocationMapping := zumata.RandomLocationMapping()

	return &LCTGSearchObject{
		Login:       l.Config.Login,
		Password:    l.Config.Password,
		CheckInDate: oneMonthLaterString,
		RegionId:    randomLocationMapping.LCTG.RegionId,
	}
}

func unwrapSearchResponse(responseBody []byte) (unwrappedSearchResults LCTGSearchResults, err error) {
	err = xml.Unmarshal(responseBody, &unwrappedSearchResults)
	return unwrappedSearchResults, err
}

func (l LCTGSearchResults) processAndSave(latency time.Duration) (err error) {
	l.Latency = latency
	if l.ReturnStatus.Success {
		return model.SaveLastSuccess(l)
	} else {
		return model.SaveLastFailure(l)
	}
}

func (l LCTGSearchResults) GetCollectionName() string {
	return string(common.LCTG)
}

func (l LCTGSearchResults) GetLatency() time.Duration {
	return l.Latency
}

func (l LCTGSearchResults) GetMessage() string {
	if l.ReturnStatus.Success {
		return ""
	}
	return l.ReturnStatus.Exception
}

func (l LCTGSearchResults) GetNumberOfResults() int {
	return len(l.PropertyResults)
}
