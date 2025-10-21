package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/bid", bidHandler)
	port := 8081
	addr := ":" + strconv.Itoa(port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// context that is cancelled on SIGINT/SIGTERM for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// run server
	go func() {
		fmt.Printf("Mock bidder running on port %d\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("mock bidder listen error: %v\n", err)
		}
	}()

	// wait for shutdown signal
	<-ctx.Done()
	fmt.Println("Shutdown signal recieved, shutting down mock bidder...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("mock bidder shutdown error: %v\n", err)
	} else {
		fmt.Println("Mock bidder stopped gracefully")
	}
}
