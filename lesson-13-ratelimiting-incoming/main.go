package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("Service started")
	l := initLimiter()
	http.HandleFunc(l.psiCheck("/status", status))
	http.ListenAndServe(":8888", nil)
}

type Limiter struct {
	requests int
	sync.Mutex
}

type Handler func(http.ResponseWriter, *http.Request)

func initLimiter() *Limiter {
	var l Limiter
	go func() {
		tick := time.Tick(time.Second * 5)
		for {
			select {
			case <-tick:
				l.Lock()
				l.requests = 0
				l.Unlock()
			}
		}
	}()
	return &l
}

func (l *Limiter) psiCheck(pattern string, handler Handler) (string, Handler) {
	h := http.HandlerFunc(handler)
	return pattern, func(w http.ResponseWriter, r *http.Request) {
		l.Lock()
		// check if requests exceeded
		if l.requests < 1 {
			l.requests++
			l.Unlock()
			// serve request
			h.ServeHTTP(w, r)
			return
		}
		l.Unlock()
		// send back error response
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error":"Request limit exceeded"}`))
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request received")
	ctx := r.Context()
	select {
	case <-ctx.Done():
		fmt.Println("Omg something happened")
	default:
	}
	fmt.Println("Do your job")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`))
}
