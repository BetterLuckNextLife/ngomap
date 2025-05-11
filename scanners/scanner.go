package scanners

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
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

func ScanPortRAW(host string, port int, protocol string, timeout int) (int, bool) {
	localIP, err := GetOutIP(host, protocol)
	if err != nil {
		return 0, false
	}
	packet, err := BuildSYN(localIP, host, port, port)
	if err != nil {
		return 0, false
	}
	err = SendRawPacket(host+":"+strconv.Itoa(port), packet)
	if err != nil {
		return int(port), false
	} else {
		return int(port), true
	}
}

// Grabs a banner from port as a string
func GrabBanner(host string, port string) (bool, string) {
	address := host + ":" + port

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return false, ""
	}

	return true, string(buffer[:n])
}

// Scans all the ports on a host, returns a slice of found ports
func ScanHost(host, protocol string, timeout int, workerCount int) []int {
	fmt.Printf("\033[1m\033[34m[*]\033[0m Scanning host %s\n", host)

	var wg sync.WaitGroup

	// Create a job channel and fill it with queued for scan
	jobs := make(chan int, 1000)
	totalPorts := 65535
	bar := progressbar.Default(int64(totalPorts), fmt.Sprintf("Host %s", host))
	go func() {
		for port := 1; port <= 65535; port++ {
			jobs <- port
		}
		close(jobs)
	}()

	// Create a result channel
	result := make(chan int, 1000)
	// Start workers and write found ports to the result channel
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				//fmt.Printf("Starting scan on %d\n", port)
				_, open := ScanPortRAW(host, port, protocol, timeout)
				bar.Add(1)
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
	if len(resultPorts) > 0 {
		fmt.Printf("\033[1m\033[32m[+]\033[0m Host %s: %d open ports\n", host, len(resultPorts))
	} else {
		fmt.Printf("\033[1m\033[31m[-]\033[0m Host %s: %d open ports\n", host, len(resultPorts))
	}
	return resultPorts
}

func ScanNetwork(address string, mask string, protocol string, timeout int, workerCount int, threads int) []ScanResult {

	var wg sync.WaitGroup

	jobs := make(chan uint32, workerCount)
	mask_int, _ := strconv.Atoi(mask)
	hostCount := uint32(1 << (32 - mask_int))
	ip := net.ParseIP(address)

	go func() {
		startAddress := IpToUint32(ip)

		endAddress := startAddress + hostCount - 1

		for host := startAddress; host <= endAddress; host++ {
			jobs <- host
		}
		close(jobs)
	}()

	// Create a result channel
	result := make(chan ScanResult, hostCount)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range jobs {
				foundPorts := ScanHost(Uint32ToIP(host).String(), protocol, timeout, threads)
				if len(foundPorts) > 0 {
					result <- ScanResult{host, foundPorts}
				}
			}
		}()
	}

	// Wait for all workers to finish and close the channel
	go func() {
		wg.Wait()
		close(result)
	}()

	hostResults := []ScanResult{}
	for hostResult := range result {
		hostResults = append(hostResults, hostResult)
	}

	return hostResults
}
