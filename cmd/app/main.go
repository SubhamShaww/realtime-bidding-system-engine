package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func init() {
	testPort := 6060
	go func() {
		fmt.Printf("pprof server running on : %d\n", testPort)
		http.ListenAndServe("localhost:"+strconv.Itoa(testPort), nil)
	}()
}
