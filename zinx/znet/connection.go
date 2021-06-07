package znet

import (
	"ZinxLearning/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

//link module
type Connection struct {
	//The currently link socket
	Conn *net.TCPConn

	//link id
	ConnID uint32

	//The currently link status
	isClosed bool

	//Notify currently link channel which had exit
	ExitChan chan bool

	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit , remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆解包的对象
		dp := NewDataPack()

		//读取客户端的msg head 二进制流8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgId 和 msgDatalen 放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据DataLen 再次读取Data , 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTCPConnection(), data)
			if err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//find router call from register bind
		go func(request ziface.IRequest) {
			//c.Router.PreHandle(request)
			c.Router.Handle(request)
			//c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start() ... ConnID = ", c.ConnID)

	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	close(c.ExitChan)
}

//get the link bind socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//get Id of the link module
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//get remote client ip port status pf tcp
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.LocalAddr()
}

//提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("Connection closed when send msg")
	}
	dp := NewDataPack()

	binaryMsg , err :=dp.Pack(NewMsgPackage(msgId , data))
	if err != nil{
		fmt.Println("dp pack err :" , err);
		return err
	}

	_ , err = c.Conn.Write(binaryMsg)
	if err != nil{
		fmt.Println("send msg err :" , err)
		return errors.New("conn write error")
	}

	return nil
}
