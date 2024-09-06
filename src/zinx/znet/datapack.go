package znet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go_code/src/zinx/ziface"
	"net"
)

/*
	封包拆包模块，协议为TLV
*/

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// Pack 封包方法
// | DataLen | MsgID | Data |
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{}) // 创建一个字节缓冲

	// 写DataLen
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写MsgID
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写Data
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (dp *DataPack) UnpackHead(headData []byte) (ziface.IMessage, error) {
	reader := bytes.NewReader(headData)

	msg := &Message{}

	// 读DataLen
	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	return msg, nil
}

func (dp *DataPack) Unpack(conn net.Conn, headData []byte) (ziface.IMessage, error) {
	// Unpack the header
	msg, err := dp.UnpackHead(headData)
	if err != nil {
		return nil, err
	}

	// If there is data, read it from the connection
	if msg.GetDataLen() > 0 {
		msgData := make([]byte, msg.GetDataLen())
		if _, err := conn.Read(msgData); err != nil {
			return nil, err
		}
		msg.SetData(msgData)
	}

	return msg, nil
}

func (dp *DataPack) GetMsgFromConn(conn net.Conn) (ziface.IMessage, error) {
	headData := make([]byte, 8)
	_, err := conn.Read(headData)
	if err != nil {
		return nil, err
	}

	msg, err := dp.Unpack(conn, headData)
	if err != nil {
		fmt.Println("unpack err ", err)
		return nil, err
	}

	return msg, nil
}
