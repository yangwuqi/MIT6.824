package ipc

import (
	"encoding/json"
	"fmt"
)

type Request struct{
	Method string "method"//说实话这个在结构体里面的元素后面加字符串的写法我没见过，是设定初始值的用法？
	Params string "params"
}
type Response struct {
	Code string "code"
	Body string "params"
}
type Server interface {
	Name() string
	Handle(method,params string)*Response
}
type IpcServer struct{
	Server//哇，直接写一个接口名放这就行了？
}
func NewIpcServer(server Server)*IpcServer{
	return &IpcServer{server}
}
func(server *IpcServer)connect()chan string{
	session:=make(chan string,0)
	go func(c chan string){
		for{
			request:=<-c
			if request=="close"{
				break
			}
			var req Request
			err:=json.Unmarshal([]byte(request),&req)
			if err!=nil{
				fmt.Println("非法的请求格式！",request)
			}
			resp:=server.Handle(req.Method,req.Params)
			b,err:=json.Marshal(resp)
			c<-string(b)
		}
		fmt.Printf("session关闭！\n")
	}(session)
	fmt.Println("一个新的session被成功创建！")
	return session
}