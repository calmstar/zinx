package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"zinx/ziface"
)

/*
定义一个全局的对象
*/
var GlobalObject *GlobalObj

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数也可以通过 用户根据 zinx.json来配置
*/
type GlobalObj struct {
	TcpServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TcpPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	Version   string         //当前Zinx版本号

	MaxPacketSize uint32 //都需数据包的最大值
	MaxConn       uint32 //当前服务器主机允许的最大链接个数

	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

func init() {
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          8080,
		Name:             "default-zinx",
		Version:          "v0.3",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   5,
		MaxWorkerTaskLen: 1024,
	}
	GlobalObject.reload()
}

func (g *GlobalObj) reload() {
	fileData, err := ioutil.ReadFile("../conf/zinx.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(fileData, &GlobalObject)
	if err != nil {
		fmt.Println("fileData Unmarshal err: ", err)
		return
	}
}
