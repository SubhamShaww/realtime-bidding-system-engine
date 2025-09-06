package engine

import (
	"container/heap"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

type Bid = utils.Bid
type PriorityQueue = utils.PriorityQueue

// simulate bid processing with HTTP call to external bidder
func processBid(bid Bid, sem chan struct{}, wg *sync.WaitGroup, metrics chan time.Duration) {
	defer wg.Done()
	sem <- struct{}{}        // acquire slot
	defer func() { <-sem }() // release slot

	start := time.Now()
	fmt.Printf("Bid %d with (Priority %d) started at %s\n", bid.ID, bid.Priority, start.Format("15:04:05.000"))

	// make HTTP request to mock bidder
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	req, _ := http.NewRequest("GET", "http://localhost:8081/bid", nil)

	resp, err := client.Do(req)
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("Bid %d failed: %v\n", bid.ID, err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Bid %d succeeded: %s (Duration: %v)\n", bid.ID, string(body), duration)
		metrics <- duration
	} else {
		fmt.Printf("Bid %d failed: %s\n", bid.ID, string(body))
	}
}

func RunEngineWithBids(inputBids []Bid, maxConcurrentBidders int, metricsSize int) {
	fmt.Println("Started realtime bidding system engine...")
	sem := make(chan struct{}, maxConcurrentBidders)
	var wg sync.WaitGroup
	metricsChan := make(chan time.Duration, metricsSize)

	// create and prioritize bids
	bids := &PriorityQueue{}
	heap.Init(bids)
	for _, bid := range inputBids {
		heap.Push(bids, bid)
	}
	// for i := 1; i <= 10; i++ {
	// 	heap.Push(bids, Bid{ID: i, Priority: 10 - i}) // Higher ID = Lower Priority
	// }

	// process bids by priority
	for bids.Len() > 0 {
		bid := heap.Pop(bids).(Bid)
		wg.Add(1)
		go processBid(bid, sem, &wg, metricsChan)
	}

	wg.Wait()
	close(metricsChan)

	// Aggregate metrics
	var metrics utils.Metrics
	for d := range metricsChan {
		metrics.Log(d, true) // true = successful (since only succesful durations are sent)
	}
	metrics.TotalBids = len(inputBids)
	metrics.TimedOut = metrics.TotalBids - metrics.Successful

	metrics.Report()
	fmt.Println("All bids processed.")
}
