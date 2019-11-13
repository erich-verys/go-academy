package throttle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Throttle(delay time.Duration, client *http.Client, reqQ chan struct{}, resQ chan []byte) {
	for {
		time.Sleep(delay)
		select {
		case <-reqQ:
			go call(client, resQ)
		default:
			fmt.Println("No requests in the queue")
		}
	}
}

func call(client *http.Client, resQ chan []byte) {
	// perform request
	response, err := client.Get("http://localhost:8888/status")
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	// read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	// enqueue the response body into the request queue
	resQ <- body
}
