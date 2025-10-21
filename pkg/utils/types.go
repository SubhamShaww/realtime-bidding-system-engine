package utils

import (
	"fmt"
	"time"
)

// Bid represents a bid request
type Bid struct {
	ID       int
	Priority int // Higher value = higher priority
}

type BidResult struct {
	BidID    int
	Duration time.Duration
	Success  bool
}

// PriorityQueue implements heap.Interface
type PriorityQueue []Bid

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Priority > pq[j].Priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x any)        { *pq = append(*pq, x.(Bid)) }
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	bid := old[n-1]
	*pq = old[0 : n-1]
	return bid
}

// Metrics
type Metrics struct {
	TotalBids    int
	Successful   int
	TimedOut     int
	TotalLatency time.Duration
}

func (m *Metrics) Log(duration time.Duration, success bool) {
	m.TotalBids++
	if success {
		m.Successful++
		m.TotalLatency += duration
	} else {
		m.TimedOut++
	}
}

func (m *Metrics) Report() {
	fmt.Printf("Successful: %d\n", m.Successful)
	fmt.Printf("Avg Latency: %v\n", m.TotalLatency/time.Duration(m.Successful))
	fmt.Printf("Timed Out: %d\n", m.TimedOut)
	fmt.Printf("Total Bids: %d\n", m.TotalBids)
}
