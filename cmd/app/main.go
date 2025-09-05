package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/engine"
	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

func init() {
	testPort := 6060
	go func() {
		fmt.Printf("pprof server running on : %d\n", testPort)
		http.ListenAndServe("localhost:"+strconv.Itoa(testPort), nil)
	}()
}

func main() {
	fmt.Println("Realtime Bidding System Main Application")
	bids := utils.GenerateBids(10)
	maxConcurrentBidders := 5
	metricsSize := 20
	engine.RunEngineWithBids(bids, maxConcurrentBidders, metricsSize)
}
