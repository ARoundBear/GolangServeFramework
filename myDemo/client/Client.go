package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("net dial err : ", err)
		return
	}

	buf := make([]byte, 1024)

	fmt.Println("please input client data...")
	for {
		conn.Write([]byte("111"))

		cnt, err := conn.Read(buf)

		if err != nil {
			fmt.Printf("client read data fail %s\n", err)
			continue
		}

		fmt.Println("client response" + string(buf[:cnt]))
	}
}
