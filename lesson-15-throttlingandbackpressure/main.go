package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-academy/lesson-15-throttlingandbackpressure/server"
	"github.com/go-academy/lesson-15-throttlingandbackpressure/throttle"
)

func main() {
	fmt.Println("service started")
	exit := make(chan struct{})
	go server.Run()
	go clientTest()
	<-exit
}

type MainData struct {
	client *http.Client
	reqQ   chan struct{}
	resQ   chan []byte
}

func clientTest() {
	// create custom client to include a timeout. note that there is no
	// defaut Go timeout. this must be done for production
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// create an instance of MainData and save the custom client
	var md MainData
	md.client = client
	md.reqQ = make(chan struct{}, 1) // request queue
	md.resQ = make(chan []byte, 100) // response queue
	// start throttler
	go throttle.Throttle(time.Second*1, md.client, md.reqQ, md.resQ)
	// push requests to the reqQ
	go enqueue(md.reqQ)
	go dequeue(md.resQ)
}

// enqueue to the request queue
func enqueue(reqQ chan struct{}) {
	for {
		time.Sleep(250 * time.Millisecond)
		select {
		case reqQ <- struct{}{}:
			fmt.Println("request queued")
		default:
			fmt.Println("error: request not queued: queue backed up")
		}
	}
}

// dequeue from the response queue
func dequeue(resQ chan []byte) {
	for {
		select {
		case data := <-resQ:
			fmt.Println(string(data))
		}
	}
}
