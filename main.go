package main

import (
	"fmt"
	"log"
	"time"

	"ratelimiter/registry"
)

func main() {

	rl, err := registry.NewRateLimiter(3, time.Minute, 5*time.Minute)
	if err != nil {
		log.Fatalf("failed to create rate limiter: %v", err)
	}

	for i := 0; i < 5; i++ {
		if rl.Allow("user-A") {
			fmt.Println("user-A: allowed")
		} else {
			fmt.Println("user-A: denied")
		}
	}

	if rl.Allow("user-B") {
		fmt.Println("user-B: allowed")
	} else {
		fmt.Println("user-B: denied")
	}

	_, err = registry.NewRateLimiter(5, 0, time.Minute)
	if err != nil {
		fmt.Println("ожидаемая ошибка на невалидном конфиге:", err)
	}
}
