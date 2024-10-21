package main

import "zinx/znet"

func main() {
	s := znet.NewServer("[zinx]v1.0")
	s.Serve()

	select {}
}
