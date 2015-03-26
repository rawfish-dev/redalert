package searcher

// import (
// 	"net/http"
// 	"time"
// )

// // Setup servers to be collected by main
// func init() {

// }

// func (s *Searcher) Healthcheck(time.Duration, error) {

// 	startTime := time.Now()
// 	s.log.Println("Pinging: ", s.Name)

// 	req, err := http.NewRequest("GET", s.Address, nil)
// 	if err != nil {
// 		return 0, errors.New("redalert ping: failed parsing url in http.NewRequest " + err.Error())
// 	}

// 	req.Header.Add("User-Agent", "Redalert/1.0")
// 	resp, err := GlobalClient.Do(req)

// 	endTime := time.Now()
// 	latency := endTime.Sub(startTime)
// 	s.log.Println(white, "Analytics: ", latency, reset)

// 	if resp != nil {
// 		io.Copy(ioutil.Discard, resp.Body)
// 		resp.Body.Close()
// 	}
// 	if err != nil {
// 		return latency, errors.New("redalert ping: failed client.Do " + err.Error())
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return latency, errors.New("redalert ping: non-200 status code. status code was " + strconv.Itoa(resp.StatusCode))
// 	}
// 	s.log.Println(green, "OK", reset, s.Name)

// 	return latency, nil
// }
