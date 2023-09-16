package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const PORT = ":8000"

type PizzaServer struct{}

func (ps *PizzaServer) start() {
	ln, err := net.Listen("tcp", PORT)
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
	buff := new(bytes.Buffer)
	for {
		var size int64
		binary.Read(conn, binary.LittleEndian, &size)
		n, err := io.CopyN(buff, conn, size)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(buff.Bytes())
		fmt.Printf("Received %d bytes of pizza over the network\n", n)
	}
}

func sendPizza(size int) error {
	pizza := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, pizza)
	if err != nil {
		return err
	}
	conn, err := net.Dial("tcp", PORT)
	if err != nil {
		return err
	}
	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(pizza), int64(size))
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
