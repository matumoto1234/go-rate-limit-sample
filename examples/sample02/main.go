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

// RateLimitHTTPClient : Rate Limited HTTP Client
type RateLimitHTTPClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

func (c *RateLimitHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()

	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

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

	req, _ := http.NewRequest("GET", "http://localhost:3000", nil)

	for i := 0; i < 100; i++ {
		resp, err := c.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(resp.StatusCode)
			return
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit reached after %d requests", i)
			return
		}
	}
}
