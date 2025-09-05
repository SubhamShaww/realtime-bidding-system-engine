package engine

import (
	"fmt"
	"testing"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

func TestRunEngine(t *testing.T) {
	fmt.Println("Testing RunEngine function...")
	bids := utils.GenerateBids(10)
	maxConcurrentBidders := 5
	metricsSize := 20
	engine.RunEngineWithBids(bids, maxConcurrentBidders, metricsSize)
}
