package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//负责测试的datapack拆包 封包的测试单元
func TestDataPack(t *testing.T) {
	listener , err := net.Listen("tcp" , "127.0.0.1:7777")
	if err != nil{
		fmt.Println("Server listen err:" , err)
		return
	}

	//创建一个go承载 负责从客户端处理业务
	go func(){
		//从客户端读取数据，拆包处理
		for{
			conn , err := listener.Accept()
			if err != nil{
				fmt.Println("server accept error" , err)
				return
			}

			go func (conn net.Conn)  {
				//处理客户端的请求消息
				dp := NewDataPack()
				for{
					headData := make([]byte , dp.GetHeadLen())
					_ , err := io.ReadFull(conn , headData)
					if err != nil{
						fmt.Println("read headerror")
						break
					}

					msgHead , err := dp.UnPack(headData)
					if err != nil{
						fmt.Println("Server uppack err" , err)
						return
					}

					if msgHead.GetMsgLen() > 0{
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度再次从io流中读取
						_,err := io.ReadFull(conn , msg.Data)
						if err != nil{
							fmt.Println("server unpack data err: " , err)
							return
						}

						//完整的消息已经读取完毕
						fmt.Println("--> Recv MsgId:" , msg.Id , "datalen = " , msg.DataLen , "data = " , string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	conn , err := net.Dial("tcp" , "127.0.0.1:7777")
	if err != nil{
		fmt.Println("client dial err: " , err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程 ， 封装两个msg一同发送
	//封装第一个msg1
	msg1 := &Message{
		Id : 100,
		DataLen: 4,
		Data: []byte{'z' , 'i' , 'n' , 'x'},
	}
	sendData1 , err := dp.Pack(msg1)
	if err != nil{
		fmt.Println("client pack msg1 error" , err)
		return
	}

	msg2 := &Message{
		Id : 1,
		DataLen: 7,
		Data: []byte{'n' , 'i' , 'h' , 'a' , 'o' , '!' , '!'},
	}
	sendData2 , err := dp.Pack(msg2)
	if err != nil{
		fmt.Println("client pack msg1 error" , err)
		return
	}

	sendData1 = append(sendData1, sendData2...)

	conn.Write(sendData1)

	select{}
 }
