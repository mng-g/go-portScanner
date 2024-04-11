package main

import (
	"flag"
	"fmt"

	"github.com/mng-g/go-portScanner/portformat"

	"net"
	"sort"
)

func worker(host string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

var (
	targetHost  string
	targetPorts string
)

func init() {
	flag.StringVar(&targetHost, "h", "localhost", "The host to scan")
	flag.StringVar(&targetPorts, "p", "1-1024", "The range of ports to scan")
	flag.Parse()
}

func main() {

	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	portsList, err := portformat.Parse(targetPorts)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

	for i := 0; i < cap(ports); i++ {
		go worker(targetHost, ports, results)
	}

	go func() {
		for _, port := range portsList {
			ports <- port
		}
	}()

	for i := 0; i < len(portsList); i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
