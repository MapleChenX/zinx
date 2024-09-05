package utils

import (
	"encoding/json"
	"fmt"
	"go_code/src/zinx/ziface"
	"os"
)

/*
   全局参数配置模块
*/

type GlobalObj struct {
	// Server
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         `json:"host"`    // 当前服务器主机监听的IP
	TcpPort   int            `json:"tcpPort"` // 当前服务器主机监听的端口号
	Name      string         `json:"name"`

	// Zinx
	Version          string `json:"version"`        // 当前Zinx版本号
	MaxConn          int    `json:"maxConn"`        // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 `json:"maxPackageSize"` // 当前Zinx框架数据包的最大值
	WorkerPoolSize   uint32 `json:"workerPoolSize"` // 当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 `json:"maxWorkerTask"`  // 每个worker对应的消息队列的任务的数量最大值
}

var GlobalVar *GlobalObj

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	GlobalVar = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		Host:             "0.0.0.0",
		TcpPort:          1234,
		MaxConn:          100,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 从配置文件中加载一些参数
	GlobalVar.Reload()

	// 打印全局变量的数据
	fmt.Println(GlobalVar)
}

func (g *GlobalObj) Reload() {
	// 从配置文件中加载一些参数
	file, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &GlobalVar)
	if err != nil {
		panic(err)
	}
}
