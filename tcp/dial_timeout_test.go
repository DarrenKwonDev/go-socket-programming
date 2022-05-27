package tcp

import (
	"net"
	"syscall"
	"testing"
	"time"
)

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	// net.DialTimeout를 사용할 수도 있지만 좀 더 자세한 설정을 위해서 Dialer를 직접 사용
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err: "connection timed out",
				Name: addr,
				Server: "127.0.0.1",
				IsTimeout: true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	c, err := DialTimeout("tcp", "10.0.0.1:http", time.Second*5)
	if err == nil {
		c.Close()
		t.Fatal("timeout not occurred")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(err)
	}
	if !nErr.Timeout() {
		t.Fatal("error is not timeout")
	}
}