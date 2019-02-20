package main

import (
	"fmt"
	"net/rpc"
	"net/http"
	"errors"
)

func main(){
	rpcDemo()
}
type Arith int
func rpcDemo(){
	arith:=new(Arith)
	rpc.Register(arith)//注意这里这个Register必不可少，必须将这个实例化的类型注册进rpc才能通过客户端调用到这个类型的方法或者说函数
	fmt.Printf("我们在服务器程序定义了一个玩意叫Arith，其实就是个int，并且会给Arith这个玩意写配套的函数来接受数据处理并像客户端返回结果\n")
	fmt.Printf("Arith: %v\n",arith)
	rpc.HandleHTTP()
	err:=http.ListenAndServe(":9999",nil)
	if err!=nil{
		fmt.Printf("出现错误啦！错误是\n")
		fmt.Printf("%v\n",err.Error())
	}
}

type Args struct{
	A int
	B int
}
type Quotient struct{
	Quo int
	Rem int
}

func (t *Arith)Multiply(args *Args,reply *int)error{
	*reply=args.A*args.B
	fmt.Printf("在服务器端执行了%d*%d的乘法，结果是%d\n",args.A,args.B,*reply)
	return nil
}

func (t *Arith)Divide(args *Args,quo *Quotient)error{
	if args.B==0{
		return errors.New("出错啦！除数是0！")
	}
	quo.Quo=args.A/args.B//int除法的整数结果
	quo.Rem=args.A%args.B//余数
	fmt.Printf("在服务器端执行了%d/%d的int除法，结果是%d，相应的余数是%d\n",args.A,args.B,quo.Quo,quo.Rem)
	return nil
}
