package engine

import (
	"fmt"
	"testing"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

// func init() {
// 	testPort := 6060
// 	go func() {
// 		fmt.Printf("pprof server running on : %d\n", testPort)
// 		http.ListenAndServe("localhost:"+strconv.Itoa(testPort), nil)
// 	}()
// }

func TestRunEngine(t *testing.T) {
	fmt.Println("Testing RunEngine function...")
	bids := utils.GenerateBids(10)
	maxConcurrentBidders := 5
	metricsSize := 20
	engine.RunEngineWithBids(bids, maxConcurrentBidders, metricsSize)
}
