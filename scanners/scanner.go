package scanners

import (
	"net"
	"strconv"
	"sync"
	"time"
)

// Scans a ports and if it is open writes it to a channel
func ScanPort(host string, port int, protocol string, timeout int) (int, bool) {
	address := host + ":" + strconv.Itoa(port)
	_, err := net.DialTimeout(protocol, address, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return port, false
	} else {
		return port, true
	}
}

// Scans all the ports on a host, returns a slice of found ports
func ScanHost(host, protocol string, timeout int, workerCount int) []int {

	var wg sync.WaitGroup

	// Create a job channel and fill it with queued for scan
	jobs := make(chan int, 100)
	go func() {
		for port := 1; port < 65535; port++ {
			jobs <- port
		}
		close(jobs)
	}()

	// Create a result channel
	result := make(chan int, 65535)
	// Start workers and write found ports to the result channel
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				port, open := ScanPort(host, port, protocol, timeout)
				if open {
					result <- port
				}
			}
		}()
	}

	// Wait for all workers to finish and close the channel
	go func() {
		wg.Wait()
		close(result)
	}()

	resultPorts := []int{}
	for port := range result {
		resultPorts = append(resultPorts, port)
	}

	return resultPorts
}
