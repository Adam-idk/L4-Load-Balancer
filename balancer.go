package main

import (
	"fmt"
	"io"
	"log"
	"net"
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

	peer := pool.GetNextValidPeer()
	if peer == nil {
		log.Println("No valid peer available")
		return
	}

	backendConn, err := net.DialTimeout("tcp", peer.URL, 2*time.Second)
	if err != nil {
		log.Println("Error connecting to backend:", err)
		return
	}
	defer backendConn.Close()

	go func() {
		_, err = io.Copy(backendConn, conn)
		if err != nil {
			log.Println("Error copying data:", err)
		}
	}()

	_, err = io.Copy(conn, backendConn)
	if err != nil {
		log.Println("Error copying data:", err)
	}
}
