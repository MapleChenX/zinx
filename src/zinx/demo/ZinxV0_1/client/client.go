package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		_, err := conn.Write([]byte("hello zinx v0.1"))
		if err != nil {
			fmt.Println("write conn err:", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err:", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf[:cnt-2], cnt)

		time.Sleep(1 * time.Second)

	}
}
