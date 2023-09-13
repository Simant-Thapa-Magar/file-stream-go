package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type PizzaServer struct{}

func (ps *PizzaServer) start() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go ps.readPizza((conn))
	}
}

func (ps *PizzaServer) readPizza(conn net.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		pizza := buff[:n]
		fmt.Println(pizza)
		fmt.Printf("Received %d bytes of pizza over the network\n", n)
	}
}

func sendPizza(size int) error {
	pizza := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, pizza)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		return err
	}
	n, err := conn.Write(pizza)
	if err != nil {
		return err
	}
	fmt.Printf("Sent %d bytes of pizza over the network\n", n)
	return nil
}

func main() {
	go func() {
		time.Sleep(2 * time.Second)
		sendPizza(2000)
	}()
	server := &PizzaServer{}
	server.start()
}
