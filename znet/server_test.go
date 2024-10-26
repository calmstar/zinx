package znet

import (
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer()
	s.Serve()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Error("conn err: ", err)
		return
	}
	_, err = conn.Write([]byte("hello unit test Server"))
	if err != nil {
		t.Error("Write err: ", err)
		return
	}

	readBuf := make([]byte, 512)
	_, err = conn.Read(readBuf)
	if err != nil {
		t.Error("read err: ", err)
		return
	}

	t.Logf("read info : %s", string(readBuf))
}
