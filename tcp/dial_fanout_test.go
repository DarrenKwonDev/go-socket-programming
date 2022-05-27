package tcp

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))

	// random port listener
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		// only one dialer connection
		if err == nil {
			conn.Close()
		}
	}()

	dial := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup) {
		defer wg.Done()

		var d net.Dialer
		// context 지닌 dialer 생성
		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			return
		}
		c.Close()

		select {
			case <- ctx.Done():
			case response <- id: // 생성된 id를 response 채널에 넣기
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res
	cancel() // deadline에 도달하기 전에 cancel하므로 context.Canceled 에러 반환
	wg.Wait()
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled, got %v", ctx.Err())
	}

	t.Logf("dialer %d retrieved the resource", response)
}