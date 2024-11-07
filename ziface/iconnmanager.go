package ziface

type IConnManager interface {
	Add(i IConnection)
	Remove(id uint32)
	Get(id uint32) (IConnection, error)
	ConnLen() uint32
	ClearConn()
}
