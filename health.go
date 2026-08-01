package main

import (
	"net"
	"sync"
	"time"
)

func CheckHealth(pool *ServerPool) {
	var wg sync.WaitGroup
	for _, backend := range pool.backends {
		wg.Add(1)
		go func(b *Backend) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", b.URL, 2*time.Second)
			if err != nil {
				b.SetAlive(false)
			} else {
				b.SetAlive(true)
				conn.Close()
			}
		}(backend)
	}
	wg.Wait()
}

func StartHealthCheck(cfg *Config, pool *ServerPool) {
	CheckHealth(pool)
	ticker := time.NewTicker(cfg.HealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		CheckHealth(pool)
	}
}
