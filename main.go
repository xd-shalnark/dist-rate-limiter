package main

import (
	"time"

	"ratelimiter/registry"
)

func main() {
	rl := registry.NewRateLimiter(3, time.Minute)

	for i := 0; i < 5; i++ {
		if rl.Allow("user-A") {
			println("user-A: allowed")      // user-A is allowed for the first 3 requests
		} else {
			println("user-A: denied")
		}
	}

	if rl.Allow("user-B") {
		println("user-B: allowed")
	} else {
		println("user-B: denied")
	}
}