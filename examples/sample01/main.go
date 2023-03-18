package main

import (
	"context"
	"fmt"

	"golang.org/x/time/rate"
)

// RateLimit処理をmain()内で行う場合

func main() {
	// NewLimiter()の第1引数は、1秒間に何回リクエストを送るかを指定する
	// 詳細 : https://pkg.go.dev/golang.org/x/time/rate#NewLimiter
	// もしくはトークンバケットアルゴリズムについて調べる
	l := rate.NewLimiter(5.0, 1) // 1 request in 5 req/sec

	ctx := context.Background()

	for i := 0; i < 10; i++ {
		if err := l.Wait(ctx); err != nil {
			panic(err)
		}

		fmt.Printf("%d\n", i)
	}
}
