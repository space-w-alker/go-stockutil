// Utilities for working with networks and network accessories.
package netutil

import (
	"fmt"
	"net"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/jackpal/gateway"
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
		if p := typeutil.Int(port); p == 0 {
			if p, err := EphemeralPort(); err == nil {
				return fmt.Sprintf("%v:%d", host, p)
			} else {
				panic(err.Error())
			}
		}
	}

	return address
}

// Retrieves the default gateway interface.
func DefaultGateway() (net.IP, error) {
	return gateway.DiscoverGateway()
}

type Address struct {
	Address   net.Addr
	Interface *net.Interface
	Gateway   net.IP
}

func (self *Address) IP() net.IP {
	if ipaddr, ok := self.Address.(*net.IPAddr); ok {
		return ipaddr.IP
	}

	return nil
}

// Return a list of routable IP addresses, along with their associated gateways and interfaces.
func RoutableAddresses() ([]*Address, error) {
	addresses := make([]*Address, 0)

	// get the default gateway
	if gw, err := DefaultGateway(); err == nil {
		if ifaces, err := net.Interfaces(); err == nil {
			// for each interface...
			for _, iface := range ifaces {
				if addrs, err := iface.Addrs(); err == nil {
					// for each address on this interface...
					for _, addr := range addrs {
						// only consider IP addresses at the moment
						if ipaddr, ok := addr.(*net.IPAddr); ok {
							network := net.IPNet{
								IP:   ipaddr.IP,
								Mask: ipaddr.IP.DefaultMask(),
							}

							// if this addresses network contains the gateway, we found a usable address
							if network.Contains(gw) {
								addresses = append(addresses, &Address{
									Address:   ipaddr,
									Interface: &iface,
									Gateway:   gw,
								})
							}
						}
					}
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return addresses, nil
}

// Retrieves the first routable IP address on any interface that falls inside of the
// system's default gateway network.  Will return nil if no IP could be found.
func DefaultAddress() net.IP {
	if addrs, err := RoutableAddresses(); err == nil && len(addrs) > 0 {
		return addrs[0].IP()
	}

	return nil
}
