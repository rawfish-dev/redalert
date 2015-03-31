package common

type PartnerType string

const (
	ST     PartnerType = "st"
	LCTG   PartnerType = "lctg"
	ZUMATA PartnerType = "zumata"

	YYYY_MM_DD_FORMAT string = "2006-01-02"
)

type PartnerMapping struct {
	ZumataId string
	STId     string
	LCTG     LCTGTuple
}

type LCTGTuple struct {
	RegionId string
	ResortId string
}
