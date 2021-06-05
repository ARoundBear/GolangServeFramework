package znet

import (
	"ZinxLearning/zinx/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    ziface.IRouter
}

// func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
// 	// fmt.Println("[Conn Handle] CallbackToClient...")
// 	if _, err := conn.Write(data[:cnt]); err != nil {
// 		fmt.Println("write back buf err", err)
// 		return errors.New("CallBackToClient error")
// 	}
// 	return nil
// }

func (s *Server) Start() {
	fmt.Printf("[Start] Server LIstener at IP: %s , Port : %d , is starting\n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error : ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen err:", s.IPVersion, "err :", err)
			return
		}

		fmt.Println("start zinx server succ , ", s.Name, "succ , LIstenning...")
		var cid uint32
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err : ", err)
				continue
			}

			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			go dealConn.Start()
			/*go func() {
				buf := make([]byte, 4096)
				for {
					cnt , err := conn.Read(buf)
					if cnt == 0 || err != nil{
						conn.Close()
						break
					}

					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err : ", err)
						continue
					}
				}
			}()*/
		}
	}()
}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!!")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
