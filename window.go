package go_otzi

import (
	"time"
)

const channelBufferSize = 1000

type DataHandler[T any] func(data []T)

// FixedWindow is a fixed window implementation
// capture data for a fixed duration and then call the handler
type FixedWindow[T any] struct {
	duration      time.Duration
	dataCh        chan []T
	addCh         chan T
	requestDataCh chan struct{}
	handler       DataHandler[T]
}

func NewFixedWindow[T any](duration time.Duration, handler DataHandler[T]) *FixedWindow[T] {
	return &FixedWindow[T]{
		duration:      duration,
		dataCh:        make(chan []T, 1), // Buffered channel to send/receive the data slice
		addCh:         make(chan T, channelBufferSize),
		handler:       handler,
		requestDataCh: make(chan struct{}),
	}
}

func (w *FixedWindow[T]) Start() {
	go w.manageData()

	timer := time.NewTicker(w.duration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			w.requestDataCh <- struct{}{} // Request current data
			currentData := <-w.dataCh     // Fetch the current data
			go w.handler(currentData)     // Handle data in a separate goroutine
		}
	}
}

func (w *FixedWindow[T]) AddData(data T) {
	w.addCh <- data
}

func (w *FixedWindow[T]) manageData() {
	var data []T
	for {
		select {
		case item := <-w.addCh:
			data = append(data, item)
		case <-w.requestDataCh:
			w.dataCh <- data
			data = make([]T, 0) // Reset after sending data
		}
	}
}
