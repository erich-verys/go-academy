package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	fmt.Println("Started")
	logger, _ = zap.NewProduction()
	http.HandleFunc("/status", addTime(status))
	http.HandleFunc("/version", addTime(version))
	http.ListenAndServe(":8888", nil)
}

func addTime(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	fmt.Println("addTime added")
	h := http.HandlerFunc(handler)
	return func(w http.ResponseWriter, r *http.Request) {
		// created a uuid
		id := uuid.New()
		// add uuid to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "request-id", id.String())
		start := time.Now().UnixNano()
		ctx = context.WithValue(ctx, "request-start", start)
		// adding new context to the request
		r = r.WithContext(ctx)
		// serving
		fmt.Println("request started")
		h.ServeHTTP(w, r)
		diff := time.Now().UnixNano() - start
		fmt.Printf("request completed\nuuid: %s\nmicro: %d\n", id.String(), diff/1000)
		fmt.Println(r.Context().Value("request-start"))
		fmt.Println(r.Context().Value("request-id"))
	}
}

func addTimeWithArgs(t time.Duration, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	fmt.Println("addTime added")
	h := http.HandlerFunc(handler)
	return func(w http.ResponseWriter, r *http.Request) {
		// get request context
		ctx := r.Context()
		// add timeout to context
		ctx, _ = context.WithTimeout(ctx, t)
		// add request-id to context
		id := uuid.New()
		ctx = context.WithValue(ctx, "request-id", id.String())
		// add start time to context
		start := time.Now().UnixNano()
		ctx = context.WithValue(ctx, "request-start", start)
		// adding new context to the request
		r = r.WithContext(ctx)
		// serve
		fmt.Println("request started")
		h.ServeHTTP(w, r)
		diff := time.Now().UnixNano() - start
		fmt.Printf("request completed\nuuid: %s\nmicro: %d\n", id.String(), diff/1000)
		fmt.Println(r.Context().Value("request-start"))
		fmt.Println(r.Context().Value("request-id"))
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("request received")
	//logger.Info("request received")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}

func version(w http.ResponseWriter, r *http.Request) {
	// create new context with timeout
	var cancel context.CancelFunc
	ctx := r.Context()
	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	// set new *http.Request
	r = r.WithContext(ctx)
	// create a done signal
	done := make(chan struct{})
	// do some things
	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()
	// handle timeout or completion of some things
	select {
	case <-ctx.Done():
		// timeout
		w.WriteHeader(400)
		w.Write([]byte(`{"error":"time exceeded"}`))
	case <-done:
		// success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":"1.0.0"}`))
	}
}
