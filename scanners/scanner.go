package scanners

import (
	"net"
	"strconv"
	"sync"
	"time"
)

// Scans a ports and if it is open writes it to a chennel
func ScanPort(host, port, protocol string, timeout_optional int, found chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := net.DialTimeout(protocol, host+":"+port, time.Duration(timeout_optional)*time.Millisecond)
	if err != nil {
	} else {
		foundport, _ := strconv.Atoi(port)
		found <- foundport
	}
}

// TODO: make threads work
// Scans all the ports on a host, prints found ports, returns a sclie of founf ports
func ScanHost(host, protocol string, timeout int) []int {
	found := make(chan int, 65535)
	var wg sync.WaitGroup

	for i := 1; i <= 65535; i++ {
		wg.Add(1)
		go ScanPort(host, strconv.Itoa(i), protocol, timeout, found, &wg)
	}

	wg.Wait()
	close(found)

	foundPorts := []int{}
	for port := range found {
		foundPorts = append(foundPorts, port)
	}
	return foundPorts
}
