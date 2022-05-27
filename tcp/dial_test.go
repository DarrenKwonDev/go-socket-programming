package tcp

import (
	"fmt"
	"io"
	"net"
	"testing"
)


func TestDial(t *testing.T) {
	// random port listener를 생성한다.
	listener, err := net.Listen("tcp", "127.0.0.1:") 
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer func() { done <- struct{}{} }()

		for {
			// dialer에게 ACK와 SYC 패킷을 보내 연결을 허가, 수립한다.
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024) // 1024 byte buffer
				for {
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
					fmt.Printf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	// listener한테 SYC 패킷을 담아 날려 connection 생성을 요청한다.
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// dialer와 listener를 닫는다.
	conn.Close() 
	<-done
	listener.Close()
	<-done
}
