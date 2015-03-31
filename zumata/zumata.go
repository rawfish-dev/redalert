package zumata

import (
	"crypto/rand"
	"encoding/binary"

	"redalert/common"
)

var locationMappings = []common.PartnerMapping{
	common.PartnerMapping{
		ZumataId: "f75a8cff-c26e-4603-7b45-1b0f8a5aa100", // SINGAPORE
		STId:     "580899",
		LCTG: common.LCTGTuple{
			RegionId: "277",
			ResortId: "1253",
		},
	},
}

func RandomLocationMapping() common.PartnerMapping {
	// Generate a random number
	var randomNumber int
	binary.Read(rand.Reader, binary.LittleEndian, &randomNumber)

	// Mod the random to get a number within our range
	randomIndex := randomNumber % len(locationMappings)

	return locationMappings[randomIndex]
}

// SINGAPORE
// {
//   "simplytravel_id": "580899",
//   "lctg": {
//     "region_id": "277",
//     "resort_id": "1253"
//   },
//   "xtg": {
//     "hobtest": "SIN",
//     "hob": "SIN",
//     "dowtest": "22394",
//     "dow": "22394",
//     "gtatest": "cSIN",
//     "gta": "cSIN"
//   },
//   "country": "Singapore",
//   "region": "",
//   "expedia_scraper_id": "180027",
//   "state": "",
//   "location": "Singapore",
//   "latitude": 1.2800945,
//   "longitude": 103.8509491,
//   "hotel_count": 104,
//   "search_count": 1876,
//   "active": true
// }
