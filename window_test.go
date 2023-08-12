package go_otzi

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func DataHandlerInt(data []int) {
	fmt.Println("Int Handler called!")
	fmt.Println("Int Data captured in the window:", data)
}

func DataHandlerString(data []string) {
	fmt.Println("String Handler called!")
	fmt.Println("String Data captured in the window:", data)
}

func TestFixedWindow(t *testing.T) {
	// Test with integers
	windowInt := NewFixedWindow(2*time.Second, DataHandlerInt)
	go windowInt.Start()

	windowInt.AddData(1)
	windowInt.AddData(2)
	time.Sleep(3 * time.Second)
	windowInt.AddData(3)

	// Allow the handler to execute before moving to the next test
	time.Sleep(3 * time.Second)

	// Test with strings
	windowString := NewFixedWindow(2*time.Second, DataHandlerString)
	go windowString.Start()

	windowString.AddData("data-1")
	windowString.AddData("data-2")
	time.Sleep(3 * time.Second)
	windowString.AddData("data-3")

	// Allow the handler to execute before ending the test
	time.Sleep(3 * time.Second)
}

func SumHandler(data []int) {
	sum := 0
	for _, v := range data {
		sum += v
	}
	fmt.Printf("Handler called! Sum of data in the window: %d\n", sum)
}

func TestBurstDataWithGoroutines(t *testing.T) {
	window := NewFixedWindow(1*time.Second, SumHandler)
	go window.Start()

	var wg sync.WaitGroup
	numGoroutines := 10
	numDataPerGoroutine := 2000
	dataRate := time.Millisecond * 5 // 5ms interval between each data addition

	// Start multiple goroutines to rapidly add data
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numDataPerGoroutine; j++ {
				window.AddData(rand.Intn(1000))
				time.Sleep(dataRate)
			}
		}()
	}

	wg.Wait()
}
