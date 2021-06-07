package main

import (
	"ZinxLearning/zinx/ziface"
	"ZinxLearning/zinx/znet"
	"fmt"
)

type PingRouter struct {
	//znet.BaseRouter
}

func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle..")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Hadnle...")
	// _, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	// if err != nil {
	// 	fmt.Println("call back ping...ping...ping error")
	// }
	fmt.Println("recv from client : msgID  = ", request.GetMsgID(),
		", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1 , []byte("ping...ping...ping..."))
	if err != nil{
		fmt.Println("request send message err : " , err)
		return 
	}
}

func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := znet.NewServer("[zinx V0.1]")

	s.AddRouter(&PingRouter{})
	s.Serve()
}
