// Utilities for working with networks and network accessories.
package netutil

import (
	"fmt"
	"net"
	"time"
)

var DefaultWaitForOpenConnectionTimeout = 5 * time.Second
var WaitForOpenCheckInterval = time.Second

func WaitForOpen(network string, address string, totaltime time.Duration, timeouts ...time.Duration) error {
	started := time.Now()
	var timeout time.Duration

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	} else {
		timeout = DefaultWaitForOpenConnectionTimeout
	}

	for time.Since(started) < totaltime {
		if conn, _ := net.DialTimeout("tcp", address, timeout); conn != nil {
			conn.Close()
			return nil
		}

		time.Sleep(WaitForOpenCheckInterval)
	}

	return fmt.Errorf("Timed out waiting for %s/%s to open", network, address)
}
