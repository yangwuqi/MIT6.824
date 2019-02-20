package main

import (
	"errors"
	"net/rpc"
	"net"
	"fmt"
	"os"
)

type Input struct {
	A int
	B int
}

type Output struct {
	C int
	D int
}

type Num int

func (n *Num)Multiply(args *Input,reply *int)error{
	*reply=args.A*args.B
	return nil
}

func (n *Num)Divide(args *Input,output *Output)error{
	if args.B==0{
		return errors.New("除数是0，出错啦！")
	}
	output.C=args.A/args.B
	output.D=args.A%args.B
	return nil
}

func main(){
	num:=new(Num)
	rpc.Register(num)
	tcpAddr,err:=net.ResolveTCPAddr("tcp",":9999")
	if err!=nil{
		fmt.Printf("服务器端的TCP初始化出错！\n")
		os.Exit(1)
	}
	listener,err:=net.ListenTCP("tcp",tcpAddr)
	if err!=nil{
		fmt.Printf("服务器端的TCP监听出错！\n")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("服务器端TCP监听Accept()出了一个错\n")
			os.Exit(1)
		} else {
			fmt.Printf("服务器端成功监听到客户端的请求\n")
		}
		rpc.ServeConn(conn) //如果TCP监听Accept()成功，那么将其绑定RPC
	}
}