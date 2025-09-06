package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func bidHandler(w http.ResponseWriter, _ *http.Request) {
	// simulate variable latency (eg. 200-900ms)
	latency := 200 + rand.Intn(700)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	// simulate random failure (10% chance)
	if rand.Float32() < 0.1 {
		http.Error(w, "timeout", http.StatusGatewayTimeout)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "bid response after %dms\n", latency)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/bid", bidHandler)
	port := 8081
	fmt.Printf("Mock bidder running on port %d\n", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
