package ziface

type IDataPack interface {
	Pack(msg IMessage) ([]byte, error) //封包方法
	UnPack([]byte) (IMessage, error)   // 拆包方法
	GetHeadLen() uint32                // 获取包头长度方法
}
