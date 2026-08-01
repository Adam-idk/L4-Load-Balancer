package main

func main() {
	cfg := NewConfig()

	pool := NewServerPool()
	for _, url := range cfg.Backends {
		pool.AddBackend(NewBackend(url))
	}

	go StartHealthCheck(pool)

	StartServer(cfg, pool)
}
