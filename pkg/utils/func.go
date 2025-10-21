package utils

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GenerateBids(n int) []Bid {
	bids := make([]Bid, n)
	for i := range n {
		bids[i] = Bid{
			ID:       i + 1,
			Priority: rand.Intn(100), // Random priority between 0-99
		}
	}

	return bids
}

func ReadBidsFromRedisStream(ctx context.Context, rdb *redis.Client, stream string, out chan<- Bid) {
	lastID := "0"
	for {
		xs, err := rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{stream, lastID},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			fmt.Println("Redis read error:", err)
			continue
		}
		for _, x := range xs[0].Messages {
			id, idErr := strconv.Atoi(x.Values["id"].(string))
			priority, priorityErr := strconv.Atoi(x.Values["priority"].(string))
			// convert id and priority to int as needed
			// ...
			if idErr != nil || priorityErr != nil {
				fmt.Println("Data conversion error:", idErr, priorityErr)
				continue
			}
			out <- Bid{
				ID:       id,
				Priority: priority,
			}
			lastID = x.ID
		}
	}
}

func SendBidResultToKafka(writer *kafka.Writer, bidID int, status string) error {
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", bidID)),
		Value: []byte(status),
	}
	return writer.WriteMessages(context.Background(), msg)
}

var (
	BidLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "bid_latency_seconds",
		Help: "Bid latency in seconds.",
	})
	BidSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bid_success_total",
		Help: "Total successful bids.",
	})
	BidTimeout = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "bid_timeout_total",
		Help: "Total timed out bids.",
	})
)

func InitPrometheus() {
	prometheus.MustRegister(BidLatency, BidSuccess, BidTimeout)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
}
