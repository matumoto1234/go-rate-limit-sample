package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimit処理をTransportに持たせる場合
// 1. Transportのフィールドとして*rate.Limiterを持たせる
// 2. RoundTrip()内でrate.Limiter.Wait()を呼ぶ

type RateLimitTransport struct {
	Transport   http.RoundTripper
	rateLimiter *rate.Limiter
}

func (rlt *RateLimitTransport) transport() http.RoundTripper {
	if rlt.Transport == nil {
		return http.DefaultTransport
	}
	return rlt.Transport
}

func (rlt *RateLimitTransport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := rlt.transport().(canceler); ok {
		cr.CancelRequest(req)
	}
}

func (rlt *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := rlt.transport().RoundTrip(req)

	if err := rlt.rateLimiter.Wait(req.Context()); err != nil {
		return nil, err
	}

	return res, err
}

func NewRateLimitTransport(transport http.RoundTripper, rateLimiter *rate.Limiter) *RateLimitTransport {
	return &RateLimitTransport{
		Transport:   transport,
		rateLimiter: rateLimiter,
	}
}

func main() {
	client := &http.Client{
		Transport: NewRateLimitTransport(
			nil,
			rate.NewLimiter(rate.Every(2), 1), // 1 request every 2 seconds
		),
	}

	for i := 0; i < 3; i++ {

		resp, err := client.Get("https://httpbin.org/get")
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("i:", i, "body:", string(body))
	}
}
