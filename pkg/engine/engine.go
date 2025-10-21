package engine

import (
	"container/heap"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/SubhamShaww/realtime-bidding-system-engine/pkg/utils"
	"github.com/segmentio/kafka-go"
)

type Bid = utils.Bid
type BidResult = utils.BidResult
type PriorityQueue = utils.PriorityQueue

// simulate bid processing with HTTP call to external bidder
func processBid(
	ctx context.Context,
	kfwriter *kafka.Writer,
	bidderURL string,
	bid Bid,
	sem chan struct{},
	wg *sync.WaitGroup,
	results chan<- BidResult,
) {
	defer wg.Done()

	// try to acquire semaphore slot or abort if context cancelled
	select {
	case sem <- struct{}{}: // acquired slot
	case <-ctx.Done():
		return
	}
	defer func() { <-sem }() // release slot when done

	start := time.Now()
	fmt.Printf("Bid %d with (Priority %d) started at %s\n", bid.ID, bid.Priority, start.Format("15:04:05.000"))

	// create request using ctx so it cancels on shutdown
	req, _ := http.NewRequestWithContext(ctx, "GET", bidderURL, nil)
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Do(req)
	duration := time.Since(start)
	if err != nil {
		fmt.Printf("Bid %d failed: %v\n", bid.ID, err)
		// report failure result
		select {
		case results <- BidResult{
			BidID:    bid.ID,
			Duration: duration,
			Success:  false,
		}:
		case <-ctx.Done():
		}
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Bid %d succeeded: %s (Duration: %v)\n", bid.ID, string(body), duration)
		if kfwriter != nil {
			_ = utils.SendBidResultToKafka(kfwriter, bid.ID, "success")
		}

		// metrics/prometheus updated by caller/aggregator or here
		select {
		case results <- BidResult{
			BidID:    bid.ID,
			Duration: duration,
			Success:  true,
		}:
		case <-ctx.Done():
		}
	} else {
		fmt.Printf("Bid %d failed: %s\n", bid.ID, string(body))
		if kfwriter != nil {
			_ = utils.SendBidResultToKafka(kfwriter, bid.ID, "failure")
		}
		select {
		case results <- BidResult{
			BidID:    bid.ID,
			Duration: duration,
			Success:  false,
		}:
		case <-ctx.Done():
		}
		utils.BidTimeout.Inc()
	}
}

func RunEngineWithBids(
	ctx context.Context,
	inputBids []Bid,
	maxConcurrentBidders int,
	metricsSize int,
	bidderURL string,
) {
	fmt.Println("Started realtime bidding system engine...")
	sem := make(chan struct{}, maxConcurrentBidders)
	var wg sync.WaitGroup
	results := make(chan BidResult, metricsSize)

	// create kafka writer with required config
	kfwriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "bids",
		Balancer: &kafka.LeastBytes{},
	})
	defer kfwriter.Close()

	// create and prioritize bids
	bids := &PriorityQueue{}
	heap.Init(bids)
	for _, bid := range inputBids {
		heap.Push(bids, bid)
	}

	// process bids by priority
	for bids.Len() > 0 {
		// if context cancelled stop accepting new bids
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled - stopping intake of new bids")
			bids = &PriorityQueue{} // clear to break
			break
		default:
		}

		bid := heap.Pop(bids).(Bid)
		wg.Add(1)
		go processBid(ctx, kfwriter, bidderURL, bid, sem, &wg, results)
	}

	// wait for workers and then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate metrics using utils.Metrics and Prometheus
	var metrics utils.Metrics
	for r := range results {
		metrics.Log(r.Duration, r.Success)
		if r.Success {
			utils.BidLatency.Observe(r.Duration.Seconds())
			utils.BidSuccess.Inc()
		} else {
			utils.BidTimeout.Inc()
		}
	}

	metrics.TotalBids = len(inputBids)
	metrics.TimedOut = metrics.TotalBids - metrics.Successful

	metrics.Report()
	fmt.Println("All bids processed.")
}
