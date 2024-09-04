package utils

//
//import (
//	"bytes"
//	"encoding/binary"
//	"go_code/src/zinx/ziface"
//	"go_code/src/zinx/znet"
//	"net"
//)
//
//type Pack struct{}
//
//var PackUtil *Pack
//
//func init() {
//	PackUtil = &Pack{}
//}
//
//// GetHeadLen 获取包头长度方法
//func (dp *Pack) GetHeadLen() uint32 {
//	// DataLen uint32(4字节) + ID uint32(4字节)
//	return 8
//}
//
//// Pack 封包方法
//// | DataLen | MsgID | Data |
//func (dp *Pack) Pack(msg ziface.IMessage) ([]byte, error) {
//	dataBuf := bytes.NewBuffer([]byte{}) // 创建一个字节缓冲
//
//	// 写DataLen
//	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
//		return nil, err
//	}
//
//	// 写MsgID
//	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
//		return nil, err
//	}
//
//	// 写Data
//	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
//		return nil, err
//	}
//
//	return dataBuf.Bytes(), nil
//}
//
//func (dp *Pack) UnpackHead(headData []byte) (ziface.IMessage, error) {
//	reader := bytes.NewReader(headData)
//
//	msg := &znet.Message{}
//
//	// 读DataLen
//	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
//		return nil, err
//	}
//
//	// 读MsgID
//	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
//		return nil, err
//	}
//
//	return msg, nil
//}
//
//func (dp *Pack) Unpack(conn net.Conn, headData []byte) (ziface.IMessage, error) {
//	// Unpack the header
//	msg, err := dp.UnpackHead(headData)
//	if err != nil {
//		return nil, err
//	}
//
//	// If there is data, read it from the connection
//	if msg.GetDataLen() > 0 {
//		msgData := make([]byte, msg.GetDataLen())
//		if _, err := conn.Read(msgData); err != nil {
//			return nil, err
//		}
//		msg.SetData(msgData)
//	}
//
//	return msg, nil
//}
