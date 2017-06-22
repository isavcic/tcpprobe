package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	hostPtr := flag.String("host", "", "hostname")
	portPtr := flag.Uint64("port", 0, "port")
	timeoutPtr := flag.Float64("timeout", 1, "timeout")
	retriesPtr := flag.Int("retries", 3, "retries")
	sleepPtr := flag.Float64("sleep", 1, "sleep")
	flag.Parse()

	host := *hostPtr
	port := *portPtr
	timeout := *timeoutPtr
	retries := *retriesPtr
	sleep := *sleepPtr

	// Sanity check
	if timeout < 0.1 {
		timeout = 0.1
	}

	if sleep < 0.5 {
		sleep = 0.5
	}

	if retries < 1 {
		retries = 1
	}

	if host == "" || port == 0 {
		os.Exit(1)
	}

	if probePort(host, port, timeout, retries, sleep) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func scanPort(host string, port uint64, timeout float64) bool {
	d := &net.Dialer{Timeout: time.Duration(uint64(timeout*1000)) * time.Millisecond}
	conn, _ := d.Dial("tcp", fmt.Sprintf("%v:%v", host, port))

	if conn != nil {
		conn.Close()
		return true
	} else {
		return false
	}
}

func probePort(host string, port uint64, timeout float64, retries int, sleep float64) bool {
	i := 0
	for {
		if scanPort(host, port, timeout) {
			return true
		} else {
			i++
			if i == retries {
				return false
			}
			time.Sleep(time.Duration(uint64(sleep*1000)) * time.Millisecond)
		}
	}
	return false
}
