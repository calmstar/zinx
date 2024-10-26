package znet

import "zinx/ziface"

type BaseRouter struct{}

func (*BaseRouter) Handle(r ziface.IRequest)     {}
func (*BaseRouter) PreHandle(r ziface.IRequest)  {}
func (*BaseRouter) PostHandle(r ziface.IRequest) {}
