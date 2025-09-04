# realtime-bidding-system-engine
Real-Time Bidding Engine with Concurrency Limit in golang

- Used goroutines to simulate bid processing
- Used a buffered channel as a semaphore
- Used WaitGroup to wait for all bids
- Logged start and end time for each bid
