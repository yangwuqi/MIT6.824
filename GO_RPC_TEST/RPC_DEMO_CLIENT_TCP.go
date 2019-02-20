package main

import (
	"os"
	"fmt"
	"net/rpc"
	"log"
	"strconv"
)

func main(){
	if len(os.Args)!=4{
		fmt.Printf("你的参数输入有误，必须输入三个参数，ip地址和端口号，运算数字1，运算数字2\n")
	}
	service:=os.Args[1]//这个是ip地址
	client,err:=rpc.Dial("tcp",service)//使用TCP拨号连接这个IP地址,连接成功产生的对象是client
	if err!=nil{
		log.Fatal("客户端程序在请求建立TCP连接的时候出错\n")
	}
	x,_:=strconv.Atoi(os.Args[2])
	y,_:=strconv.Atoi(os.Args[3])
	reply:=0
	err=client.Call("Num.Multiply",Input1{x,y},&reply)
	if err!=nil{
		log.Fatal("客户端向服务器建立TCP连接成功了，但是客户端向服务器调用Num.Multiply失败了,",err)
	}
	fmt.Printf("客户端向服务器成功调用了乘法函数，%d*%d=%d\n",x,y,reply)
	var reply2 Output1
	err=client.Call("Num.Divide",Input1{x,y},&reply2)
	if err!=nil{
		log.Fatal("客户端向服务器建立TCP连接成功了，但是客户端向服务器调用Num.Divide失败了,",err)
	}
	fmt.Printf("客户端向服务器成功调用了int除法函数，%d/%d=%d，余数是%d\n",x,y,reply2.C,reply2.D)
}
type Input1 struct {
	A int
	B int
}
type Output1 struct {
	C int
	D int
}