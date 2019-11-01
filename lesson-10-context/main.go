package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type MyContext interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

func main() {
	fmt.Println("App started")
	cancelExample()
	timeoutExample()
	deadlineExample()
	valueExample()
	childExample()
	serverExample()
	fmt.Println("App exited")
}

func cancelExample() {
	fmt.Println("cancelExample")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	select {
	case <-ctx.Done():
		fmt.Println("cancelExample: done")
	}
}

func timeoutExample() {
	fmt.Println("timeoutExample")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel() // needed to avoid a context leak
	select {
	case <-ctx.Done():
		fmt.Println("timeoutExample: done")
	}
}

func deadlineExample() {
	fmt.Println("deadlineExample")
	ctx, cancel := context.WithDeadline(context.Background(), time.Date(2009, time.November, 1, 0, 0, 0, 0, time.UTC))
	defer cancel() // needed to avoid a context leak
	select {
	case <-ctx.Done():
		fmt.Println("deadlineExample: done")
	}
}

func valueExample() {
	fmt.Println("valueExample")
	ctx := context.WithValue(context.Background(), "name", "Eric")
	fmt.Println(ctx.Value("name"))
	fmt.Println(ctx.Value("age"))
}

func childExample() {
	fmt.Println("childExample")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	go doSomething(ctx)
	select {
	case <-ctx.Done():
		fmt.Println("childExample: parent done")
	}
}

func doSomething(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	select {
	case <-ctx.Done():
		fmt.Println("childExample: child context closed")
	}
}

func serverExample() {
	fmt.Println("serverExample")
	go func() {
		time.Sleep(5 * time.Second)
		clientExample()
	}()
	http.HandleFunc("/status", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serverExample: request received")
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println("serverExample: done")
	}
}

func clientExample() {
	fmt.Println("clientExample")

	// create client
	client := http.Client{
		Timeout: time.Second * 30,
	}

	// create request
	var err error
	var request *http.Request
	//request, err = http.NewRequest("GET", "http://localhost:8080/status", nil)
	ctx, cancel := context.WithCancel(context.Background())
	request, err = http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/status", nil)
	if err != nil {
		return
	}

	// cancel the request in 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("clientExample: client cancelling request")
		cancel()
	}()

	// handle the response
	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// read the body
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Println(string(body))

}
