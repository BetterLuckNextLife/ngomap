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

// Scans a ports using raw SYN packets and if it is open writes it to a channel
func ScanPortRAW(localIP net.IP, host string, port int, protocol string, timeout int) (int, bool) {
	packet, err := BuildSYN(localIP, host, 12345, port)
	if err != nil {
		fmt.Printf("BuildSYN error on port %d: %v\n", port, err)
		return 0, false
	}
	err = SendRawPacket(host+":"+strconv.Itoa(port), packet)
	if err != nil {
		fmt.Printf("SendRawPacket error on port %d: %v\n", port, err)
		return int(port), false
	}
	return int(port), true
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

func ScanHost(host, protocol string, timeout int, workerCount int) []int {
	fmt.Printf("\033[1m\033[34m[*]\033[0m Scanning host %s\n", host)

	localIP, err := GetOutIP(host)
	if err != nil {
		fmt.Printf("\033[1m\033[33m[!]\033[0m Skipping host %s: failed to determine route (%v)\n", host, err)
		return nil
	}

	var wg sync.WaitGroup
	jobs := make(chan int, 1000)
	result := make(chan int, 1000)

	const totalPorts = 65535
	bar := progressbar.Default(int64(totalPorts), fmt.Sprintf("Host %s", host))

	// Fill jobs
	go func() {
		for port := 1; port <= totalPorts; port++ {
			jobs <- port
		}
		close(jobs)
	}()

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				_, open := ScanPortRAW(localIP.String(), port, protocol, timeout)
				bar.Add(1)
				if open {
					result <- port
				}
			}
		}()
	}

	// Close result channel after all workers done
	go func() {
		wg.Wait()
		close(result)
	}()

	// Collect results
	var resultPorts []int
	for port := range result {
		resultPorts = append(resultPorts, port)
	}

	if len(resultPorts) > 0 {
		fmt.Printf("\033[1m\033[32m[+]\033[0m Host %s: %d open ports\n", host, len(resultPorts))
	} else {
		fmt.Printf("\033[1m\033[31m[-]\033[0m Host %s: no open ports found\n", host)
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
