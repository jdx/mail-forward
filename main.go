package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
)

var TLSConfig *tls.Config

func main() {
	TLSConfig = setupTLS()
	addr := "0.0.0.0:" + port()
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("listening on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "25"
	}
	return port
}

func setupTLS() *tls.Config {
	cert, err := tls.LoadX509KeyPair("./certs/ssl.pem", "./certs/ssl.key")
	if err != nil {
		log.Println("Error loading certificate:", err)
		return nil
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ServerName:   "mail.dickey.xxx",
	}
	config.Rand = rand.Reader
	return config
}
