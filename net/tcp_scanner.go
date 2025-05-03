package net

import (
	"fmt"
	"net"
)

type TcpScanner struct{}

func NewTcpScanner() *TcpScanner {
	return &TcpScanner{}
}

func (s *TcpScanner) Scan(host string, maxPort int) {
	ports := make(chan int, 100)
	result := make(chan int)

	go func() {
		for i := 0; i < 100; i++ {
			go worker(ports, host, result)
		}
	}()

	go func() {
		for p := 1; p <= maxPort; p++ {
			ports <- p
		}
	}()

	openPorts := make([]int, 0)
	for i := 0; i < maxPort; i++ {
		p := <-result
		if p != 0 {
			openPorts = append(openPorts, p)
		}
	}

	close(ports)
	close(result)
	for _, p := range openPorts {
		fmt.Println("Open port: ", p)
	}
}

func worker(ports chan int, host string, result chan int) {
	for p := range ports {
		addr := fmt.Sprintf("%s:%d", host, p)
		if conn, err := net.Dial("tcp", addr); err == nil {
			conn.Close()
			result <- p
			//fmt.Println(addr, " open")
		} else {
			result <- 0
			//fmt.Println(addr, " close")
		}
	}
}
