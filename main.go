package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// Scans a ports and if it is open writes it to a chennel
func scanPort(host, port, protocol string, found chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := net.DialTimeout(protocol, host+":"+port, 2500*time.Millisecond)
	if err != nil {
	} else {
		foundport, _ := strconv.Atoi(port)
		found <- foundport
	}
}

// Scans all the ports on a host, prints found ports
func scanHost(host, protocol string) {
	found := make(chan int, 65535)
	var wg sync.WaitGroup

	for i := 1; i <= 65535; i++ {
		wg.Add(1)
		go scanPort(host, strconv.Itoa(i), protocol, found, &wg)
	}

	wg.Wait()
	close(found)

	for port := range found {
		fmt.Println(port)
	}
}

func main() {
	ip := "192.168.100.2"
	protocol := "tcp"
	scanHost(ip, protocol)
}
