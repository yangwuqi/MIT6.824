package main

import (
	"fmt"
	"os"
	"net/rpc"
	"log"
	"strconv"
)

type ArgsTwo struct{
	A int
	B int
}
type QuotientTwo struct{
	Quo int
	Rem int
}

func main(){
	fmt.Printf("os*******",os.Args,"*********")
	if len(os.Args)!=4{
		fmt.Printf("你输入的参数数量不对！检查一下！%v\n",os.Args[0])
		os.Exit(1)
	}else{
		fmt.Printf("你输入了%d个参数，参数数量是对的，参数分别是%v,%v,%v\n",len(os.Args)-1,os.Args[1],os.Args[2],os.Args[3])
	}
	serverAddress:=os.Args[1]
	fmt.Printf("服务器地址是%v\n",serverAddress)
	client,err:=rpc.DialHTTP("tcp",serverAddress)
	if err!=nil{
		log.Fatal("客户端的HTTP请求出错了！\n",err)
	}else{
		fmt.Printf("客户端到服务器的RPC HTTP连接成功，%v\n",os.Args[1])
	}
	x,_:=strconv.Atoi(os.Args[2])
	y,_:=strconv.Atoi(os.Args[3])
	args:=ArgsTwo{x,y}
	reply:=0
	err=client.Call("Arith.Multiply",args,&reply)
	if err!=nil{
		log.Fatal("客户端对服务器端的Arith数据类型的Multiply函数的调用出现问题了！\n",err)
	}
	fmt.Printf("执行了服务器端的对Arith数据类型的Multiply函数，传递的参数是%d, %d, 结果是 %d*%d = %d\n",args.A,args.B,args.A,args.B,reply)
	var quot QuotientTwo
	err=client.Call("Arith.Divide",args,&quot)
	if err!=nil{
		log.Fatal("客户端对服务器端的Arith数据类型的Divide函数的调用出现问题了！\n",err)
	}
	fmt.Printf("执行了服务器端的对Arith数据类型的Divide函数，传递的参数是%d, %d, 结果是 %d/%d = %d，余数是%d\n",args.A,args.B,args.A,args.B,quot.Quo,quot.Rem)
}