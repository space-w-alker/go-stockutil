// Utilities for working with networks and network accessories.
package netutil

import (
	"fmt"
	"net"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/phayes/freeport"
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

// Retrieve an open ephemeral port.
func EphemeralPort() (int, error) {
	return freeport.GetFreePort()
}

// Takes an address in the form of "host:port", looks for port zero (e.g: ":0"),
// and gets an ephemeral local port and returns that address (e.g.: ":41327").
func ExpandPort(address string) string {
	if host, port, err := net.SplitHostPort(address); err == nil {
		if p := typeutil.V(port).Int(); p == 0 {
			if p, err := EphemeralPort(); err == nil {
				return fmt.Sprintf("%v:%d", host, p)
			} else {
				panic(err.Error())
			}
		}
	}

	return address
}
