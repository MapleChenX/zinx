package main

import "go_code/src/zinx/znet"

func main() {
	server := znet.NewServer("ZinxV0.1")
	server.Serve()
}
