package searcher

import (
	"time"

	"redalert/common"
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
