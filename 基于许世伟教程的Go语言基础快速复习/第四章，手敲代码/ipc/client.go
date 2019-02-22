package ipc

import "encoding/json"

type IpcClient struct{
	conn chan string
}
func NewIpcClient(server *IpcServer)*IpcClient{
	c:=server.connect()
	return &IpcClient{c}
}
func(client *IpcClient)Call(method,params string)(resp *Response,err error){
	req:=&Request{method,params}
	b,err:=json.Marshal(req)
	if err!=nil{
		return
	}
	client.conn<-string(b)//这里书上的代码似乎漏掉了connect函数，我看看要不要补上，不然难道直接写入channel然后马上又读出来？

	str:=<-client.conn

	var resp1 Response
	err=json.Unmarshal([]byte(str),&resp1)
	resp=&resp1
	return
}

func(client *IpcClient)Close(){
	client.conn<-"close"
}
