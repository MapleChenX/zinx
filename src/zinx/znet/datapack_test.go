package znet

import (
	"io"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Errorf("listener err: %v", err)
		return
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Errorf("accept err: %v", err)
				return
			}

			go func(conn net.Conn) {
				for {
					pack := &DataPack{}
					// 只读取head数据
					headData := make([]byte, pack.GetHeadLen())
					_, err := conn.Read(headData)
					if err != nil {
						if err == io.EOF {
							t.Logf("connection closed by client")
							return
						}
						t.Errorf("read head err: %v", err)
						return
					}

					msg, err := pack.Unpack(conn, headData)
					if err != nil {
						t.Errorf("unpack err: %v", err)
						return
					}

					t.Logf("Received MsgID: %d, DataLen: %d, Data: %s",
						msg.GetMsgId(), msg.GetDataLen(), string(msg.GetData()))
				}
			}(conn)
		}
	}()

	select {}
}

func TestClient(t *testing.T) {
	// Simulate a client sending a message
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Errorf("dial err: %v", err)
		return
	}
	defer conn.Close()

	pack := NewDataPack()
	msg1 := &Message{
		Id:   1,
		Data: []byte("hello, world"),
	}
	// 设置消息长度
	msg1.SetDataLen(uint32(len(msg1.GetData())))

	packedData1, err := pack.Pack(msg1)
	if err != nil {
		t.Errorf("pack err: %v", err)
		return
	}

	msg2 := &Message{
		Id:   1,
		Data: []byte("wocaonima!"),
	}
	// 设置消息长度
	msg2.SetDataLen(uint32(len(msg2.GetData())))

	packedData2, err := pack.Pack(msg2)
	if err != nil {
		t.Errorf("pack err: %v", err)
		return
	}

	data4send := append(packedData1, packedData2...)

	_, err = conn.Write(data4send)
	if err != nil {
		t.Errorf("write err: %v", err)
		return
	}
}
