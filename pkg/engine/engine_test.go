package engine

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

func TestRunEngine(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// start a mock bidder using httptest so test is hermetic
	handler := func(w http.ResponseWriter, _ *http.Request) {
		// simulate variable latency (eg. 50-200ms) (faster for tests)
		latency := 50 + rand.Intn(150)
		time.Sleep(time.Duration(latency) * time.Millisecond)

		// simulate random failure (10% chance)
		if rand.Float32() < 0.1 {
			http.Error(w, "timeout", http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "bid response after %dms\n", latency)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	fmt.Println("Testing RunEngine function with httptest server...")
	bids := utils.GenerateBids(10)
	maxConcurrentBidders := 5
	metricsSize := 20

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// pass server.URL as bidder endpoint
	RunEngineWithBids(ctx, bids, maxConcurrentBidders, metricsSize, server.URL+"/bid")
}
