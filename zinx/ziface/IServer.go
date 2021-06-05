package ziface

//server interface
type IServer interface{
	Start()
	Stop()
	Serve()

	///Router method
	AddRouter(router IRouter)
}