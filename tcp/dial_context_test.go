package tcp

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	dl := time.Now().Add(time.Second * 5)
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel() // context GC

	var d net.Dialer
	d.Control = func(_, addr string, _ syscall.RawConn) error {
		time.Sleep(time.Second * 6) // deadline ctx보다는 길어야 함
		return nil
	}

	// dialer에 ctx를 연결 하는 사례
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}

	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Error("error is not timeout")
		}
	}

	// DeadlineExceeded is the error returned by Context.Err when the context's
	// deadline passes.
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected dealine exceeded, got %v", ctx.Err())
	}
}