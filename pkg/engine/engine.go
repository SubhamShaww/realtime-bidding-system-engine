package engine

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
)

type Bid = utils.Bid
type PriorityQueue = utils.PriorityQueue

// simulate bid processing with timeout
func processBid(bid Bid, sem chan struct{}, wg *sync.WaitGroup, metrics chan time.Duration) {
	defer wg.Done()
	sem <- struct{}{}        // acquire slot
	defer func() { <-sem }() // release slot

	start := time.Now()
	fmt.Printf("Bid %d with (Priority %d) started at %s\n", bid.ID, bid.Priority, start.Format("15:04:05.000"))

	done := make(chan struct{})
	go func() {
		// simulate variable processing time
		time.Sleep(time.Duration(400+bid.ID*100) * time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
		end := time.Now()
		duration := end.Sub(start)
		fmt.Printf("Bid %d ended at %s (Duration: %v)\n", bid.ID, end.Format("15:04:05.000"), duration)
		metrics <- duration
	case <-time.After(1 * time.Second):
		fmt.Printf("Bid %d timed out!\n", bid.ID)
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
