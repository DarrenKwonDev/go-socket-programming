package tcp

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer func() {
			conn.Close()
			close(sync)
		}()

		// 1. read, write timeout 설정
		err = conn.SetDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1)
		_, err = conn.Read(buf) // dialer가 데이터 보낼 때 까지 블로킹
		nErr, ok := err.(net.Error)
		if !ok || !nErr.Timeout() { // 2
			t.Errorf("expected timeout error, got %v", err)
		}
		
		sync <- struct{}{}

		// 3
		err = conn.SetDeadline(time.Now().Add(time.Second * 5))
		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-sync
	_, err = conn.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != io.EOF { // 4
		t.Errorf("expected EOF, got %v", err)
	}
}