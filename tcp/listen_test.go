package tcp

import (
	"net"
	"testing"
)

// net.Listen을 통해 소켓 주소 (ip와 포트)에 바인딩이 되었는지 확인합니다.
func TestListener(t *testing.T) {
	network := "tcp" // network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
	listener, err := net.Listen(network, "127.0.0.1:0") // 0번 포트이므로 무작위 포트
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = listener.Close()
	}()
	t.Logf("bound to %q", listener.Addr())
}
