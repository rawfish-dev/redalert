package searcher

import (
	"redalert/common"
	lctg "redalert/searcher/lctg"
)

var Online []common.CheckableServer

// Setup servers to be collected by main
func init() {
	// TODO:: streamline start up to auto discover searchers
	Online = append(Online, lctg.NewLCTGSearcher())
}
