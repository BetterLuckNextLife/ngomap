package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

func pingPort(host, port, protocol string, found chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := net.DialTimeout(protocol, host+":"+port, 1000*time.Millisecond)
	if err != nil {
	} else {
		foundport, _ := strconv.Atoi(port)
		found <- foundport
	}

}

func main() {
	ip := "127.0.0.1"
	//port := "22"
	protocol := "tcp"

	found := make(chan int, 65535)
	var wg sync.WaitGroup

	for i := 0; i <= 65535; i++ {
		wg.Add(1)
		go pingPort(ip, strconv.Itoa(i), protocol, found, &wg)
	}

	wg.Wait()
	close(found)

	for port := range found {
		fmt.Println(port)
	}
}
