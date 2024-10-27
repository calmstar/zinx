package ziface

type IMessage interface {
	GetData() []byte
	GetDataLen() uint32
	GetMsgId() uint32

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
