package ziface

type IRouter interface {
	Handle(r IRequest)
	PreHandle(r IRequest)
	PostHandle(r IRequest)
}
