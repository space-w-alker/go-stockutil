// Utilities for working with networks and network accessories.
package netutil

import (
	"fmt"
	"net"
	"os"
	"strings"
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

type IPAddress struct {
	IP        net.IP
	Mask      net.IPMask
	Interface net.Interface
	Gateway   net.IP
}

// Return a list of routable IP addresses, along with their associated gateways and interfaces.
func RoutableAddresses() ([]*IPAddress, error) {
	addresses := make([]*IPAddress, 0)

	// get the default gateway
	if gw, err := DefaultGateway(); err == nil {
		if ifaces, err := net.Interfaces(); err == nil {
			// for each interface...
			for _, iface := range ifaces {
				if addrs, err := iface.Addrs(); err == nil {
					// for each address on this interface...
					for _, addr := range addrs {

						// only consider IP addresses at the moment
						if network, ok := addr.(*net.IPNet); ok {
							// if this addresses network contains the gateway, we found a usable address
							if network.Contains(gw) {
								addresses = append(addresses, &IPAddress{
									IP:        network.IP,
									Mask:      network.Mask,
									Interface: iface,
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
func DefaultAddress() *IPAddress {
	if addrs, err := RoutableAddresses(); err == nil && len(addrs) > 0 {
		return addrs[0]
	}

	return nil
}

// Return the current machine's Fully-qualified domain name,
func FQDN() string {
	if hostname, err := os.Hostname(); err == nil {
		if responses, err := net.LookupIP(hostname); err == nil {
			for _, addr := range responses {
				if ipv4 := addr.To4(); ipv4 != nil {
					if ip, err := ipv4.MarshalText(); err == nil {
						if hosts, err := net.LookupAddr(string(ip)); err == nil && len(hosts) > 0 {
							fqdn := hosts[0]
							return strings.TrimSuffix(fqdn, ".")
						}
					}
				}
			}
		}

		return hostname
	} else {
		return ``
	}
}
