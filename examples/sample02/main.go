package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// RateLimit処理をClientに持たせる場合
// 1. Clientのフィールドとして*rate.Limiterを持たせる
// 2. Do()などの実際にリクエストを送る際に、rate.Limiter.Wait()を呼ぶ

//RateLimitHTTPClient : Rate Limited HTTP Client
type RateLimitHTTPClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

//Do dispatches the HTTP request to the network
func (c *RateLimitHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Comment out the below 5 lines to turn off ratelimiting
	ctx := context.Background()
	err := c.RateLimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//NewClient return http client with a ratelimiter
func NewClient(rl *rate.Limiter) *RateLimitHTTPClient {
	c := &RateLimitHTTPClient{
		client:      http.DefaultClient,
		RateLimiter: rl,
	}
	return c
}

func main() {
	rl := rate.NewLimiter(rate.Every(10*time.Second), 50) // 50 request every 10 seconds
	c := NewClient(rl)
	reqURL := "http://localhost:4000"
	req, _ := http.NewRequest("GET", reqURL, nil)
	for i := 0; i < 300; i++ {
		resp, err := c.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(resp.StatusCode)
			return
		}
		if resp.StatusCode == 429 {
			fmt.Printf("Rate limit reached after %d requests", i)
			return
		}
	}
}
