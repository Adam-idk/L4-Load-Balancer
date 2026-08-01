package main

import (
	"net"
	"time"
)

func CheckHealth(pool *ServerPool) {
	for _, backend := range pool.backends {
		conn, err := net.DialTimeout("tcp", backend.URL, 2*time.Second)
		if err != nil {
			backend.SetAlive(false)
		} else {
			backend.SetAlive(true)
			conn.Close()
		}
	}
}

func StartHealthCheck(pool *ServerPool) {
	CheckHealth(pool)
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		CheckHealth(pool)
	}
}
