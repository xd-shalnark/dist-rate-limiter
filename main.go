package main

import (
	"fmt"
	"time"

	"ratelimiter/registry"
)

func main() {
	// limit=3 запроса, window=1 минута, ttl=5 минут неактивности до эвикшна ключа
	rl := registry.NewRateLimiter(3, time.Minute, 5*time.Minute)

	for i := 0; i < 5; i++ {
		if rl.Allow("user-A") {
			fmt.Println("user-A: allowed") // первые 3 запроса пройдут
		} else {
			fmt.Println("user-A: denied")
		}
	}

	if rl.Allow("user-B") {
		fmt.Println("user-B: allowed")
	} else {
		fmt.Println("user-B: denied")
	}
}
