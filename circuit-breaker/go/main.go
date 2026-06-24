package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/sony/gobreaker/v2"
)

var cb *gobreaker.CircuitBreaker[[]byte]

var randError = errors.New("random error")

func init() {
	st := gobreaker.Settings{
		Name:        "HTTP GET",
		MaxRequests: 1,
		Interval:    time.Second * 3,
		Timeout:     time.Second * 2,
		IsSuccessful: func(err error) bool {
			if errors.Is(err, randError) {
				return false
			}

			return true
		},
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	cb = gobreaker.NewCircuitBreaker[[]byte](st)
}

var threshold = 0.5

func FetchPage(url string, shouldHaveFailure bool) func() ([]byte, error) {
	return func() ([]byte, error) {

		n := rand.Float64()
		if n < threshold && shouldHaveFailure {
			return nil, randError
		}

		return []byte("Hello World"), nil
	}
}

func Get(url string, shouldHaveFailure bool) ([]byte, error) {
	body, err := cb.Execute(FetchPage(url, shouldHaveFailure))
	if err != nil {
		return nil, err
	}

	return body, nil
}

func PrintState() {
	fmt.Println(cb.State())
}

func main() {

	go func() {
		for {
			time.Sleep(time.Millisecond * 250)
			PrintState()
		}
	}()

	time.Sleep(time.Second * 3)

	for range 10000 {
		Get("http://www.google.com/robots.txt", true)
	}

	fmt.Println("Sleeping for 3 seconds to allow the circuit breaker to reset...")
	time.Sleep(time.Second * 3)

	fmt.Println("Making requests again after the sleep...")
	for range 10000 {
		Get("http://www.google.com/robots.txt", false)
	}

	fmt.Println("Sleeping for 3 seconds to allow the circuit breaker to reset...")
	time.Sleep(time.Second * 3)

	fmt.Println("Making requests again after the sleep...")
	for range 10000 {
		Get("http://www.google.com/robots.txt", false)
	}
}
