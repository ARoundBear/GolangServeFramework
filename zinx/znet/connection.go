package znet

import (
	"ZinxLearning/zinx/ziface"
	"fmt"
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
		buf := make([]byte, 4096)
		cnt, err := c.Conn.Read(buf)
		if cnt == 0 || err != nil {
			c.Conn.Close()
			break
		}
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		req := Request{
			conn: c,
			data: buf,
		}

		//find router call from register bind
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		// fmt.Printf("server read buf is %s\n", buf[:cnt])

		// if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		// 	fmt.Println("ConnID", c.ConnID, "handle is error", err)
		// 	break
		// }
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

// send data to remote client
func (c *Connection) Send(data []byte) error {
	return nil
}
