package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func StartServer(cfg *Config, pool *ServerPool) {
	address := fmt.Sprintf(":%d", cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()
	log.Println("Listening on", address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go HandleConnection(connection, pool)
	}
}

func HandleConnection(conn net.Conn, pool *ServerPool) {
	defer conn.Close()

	var backendConn net.Conn
	var err error

	attempts := len(pool.backends)
	for i := 0; i < attempts; i++ {
		peer := pool.GetNextValidPeer()
		if peer == nil {
			break
		}

		backendConn, err = net.DialTimeout("tcp", peer.URL, 2*time.Second)
		if err == nil {
			break
		}

		log.Printf("Failed to dial backend %s, marking offline\n", peer.URL)
		peer.SetAlive(false)
	}

	if backendConn == nil {
		log.Println("No healthy backend available")
		return
	}
	defer backendConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = io.Copy(backendConn, conn)
		if tc, ok := backendConn.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
	}()

	go func() {
		defer wg.Done()
		_, _ = io.Copy(conn, backendConn)
		if tc, ok := conn.(*net.TCPConn); ok {
			tc.CloseWrite()
		}
	}()

	wg.Wait()
}
