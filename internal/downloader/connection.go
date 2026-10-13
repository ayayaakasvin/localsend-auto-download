package downloader

import (
	"fmt"
	"net"
	"time"
)

const localsendDefaultPort = ":53317"

func (d *Downloader) establishIPInsideNetwork(host string) error {
	if err := checkHost(host); err != nil {
		return fmt.Errorf("Host check error: %w", err)
	}

	ip, err := localIP(host)
	if err != nil {
		return fmt.Errorf("Local IP error: %w", err)
	}

    d.host = host
	d.sessionID = ip.String()

	return nil
}

func localIP(host string) (net.IP, error) {
    conn, err := net.Dial("udp", host+localsendDefaultPort)
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    addr, ok := conn.LocalAddr().(*net.UDPAddr)
    if !ok {
        return nil, fmt.Errorf("unexpected local address type")
    }

	if addr.IP.String() == "" {
		return nil, fmt.Errorf("ip is empty")
	}

    return addr.IP, nil
}

func checkHost(host string) error {
    conn, err := net.DialTimeout(
        "tcp",
        host+localsendDefaultPort,
        2*time.Second,
    )
    if err != nil {
        return err
    }
    return conn.Close()
}