package port

import (
	"net"
	"strconv"
	"sync"
	"time"
)

type ScanResult struct {
	Port  int
	State string
}

const (
	open   = "Open"
	closed = "Closed"
)

func ScanPort(protocol, hostname string, port int) ScanResult {
	result := ScanResult{Port: port}
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	if err != nil {
		result.State = closed
		return result
	}

	defer conn.Close()

	result.State = open
	return result
}

func InitialScan(hostname string) map[string][]ScanResult {
	portMap := make(map[string][]ScanResult)
	ports := make(chan ScanResult, 65535)
	var wg sync.WaitGroup

	for i := 1; i <= 65535; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done() // Ensure wg.Done() is called to prevent deadlock
			portResult := ScanPort("tcp", hostname, port)
			ports <- portResult
		}(i)
	}

	go func() {
		wg.Wait()
		close(ports)
	}()

	for port := range ports {
		if _, exists := portMap[port.State]; !exists {
			portMap[port.State] = []ScanResult{}
		}
		portMap[port.State] = append(portMap[port.State], port)
	}

	return portMap
}

func GetOpenPorts(hostname string) []ScanResult {
	return InitialScan(hostname)[open]
}

func GetClosedPorts(hostname string) []ScanResult {
	return InitialScan(hostname)[closed]
}

func GetFullScan(hostname string) []ScanResult {
	portMap := InitialScan(hostname)
	return append(portMap[open], portMap[closed]...)
}
