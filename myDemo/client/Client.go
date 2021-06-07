package main

import (
	"ZinxLearning/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("net dial err : ", err)
		return
	}

	fmt.Println("please input client data...")
	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.5 client Test Message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("conn write err : ", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read header error")
			break
		}

		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack head msg err :", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("client read message body err :", err)
				break
			}

			fmt.Println("--> Recv Server Msg : ID = ", msg.Id, ", len = ",
				msg.GetMsgLen(), ", msg Data : ", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
}
