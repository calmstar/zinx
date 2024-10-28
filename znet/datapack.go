package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

// 工具结构体，纯粹用来提供解压缩方法给外部使用
type DataPack struct{}

func NewDataPack() *DataPack {
	return new(DataPack)
}

// 数据包格式排列 Id uint32; dataLen uint32; Data 前面制定长度占用
func (dp *DataPack) GetHeadLen() uint32 {
	// message结构体：Id uint32; dataLen uint32 各自占用4个字节，总共8个字节
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	// 按顺序写入，1msg.GetDataLen(); 2msg.GetMsgId(); 3msg.GetData()
	err := binary.Write(buff, binary.LittleEndian, msg.GetDataLen()) // 小端字节序
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	err = binary.Write(buff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		fmt.Println("err: ", err)
		return nil, err
	}

	buff.Write(msg.GetData())
	return buff.Bytes(), nil
}

func (dp *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	// 拆包，解压head头部数据
	dataBuff := bytes.NewReader(data)
	msg := &Message{}

	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}
	err = binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPacketSize > 0 && utils.GlobalObject.MaxPacketSize < msg.DataLen {
		return nil, errors.New("packet too large")
	}

	return msg, nil
}
