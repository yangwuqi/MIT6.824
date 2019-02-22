package main
import(
	"./cg"
	"./ipc"
	"fmt"
	"strconv"
	"strings"
	"bufio"
	"os"
)
/*
var counter int=0

func Count(lock *sync.Mutex,index int){
	lock.Lock()
	counter++
	fmt.Printf("已将counter锁住！当前协程编号是%d，当前counter值是%d\n",index,counter)
	lock.Unlock()
}
func main(){
	lock:=sync.Mutex{}//lock是一个锁的实例
	for i:=0;i<10;i++{
		go Count(&lock,i)//把lock这个锁传给Count用于counter计数加+1时确保只有一个函数在操作counter，用指针确保是同一个锁，这里原书代码是错的没加&
	}
	for{
		lock.Lock()
		c:=counter
		lock.Unlock()
		runtime.Gosched()//这个玩意用于出让CPU时间片
		if c>=10{
			break
		}
	}
}
*/
/*
func Count(ch chan int,index int){
	ch<-1
	fmt.Printf("正在处理！编号%d\n",index)
}
func main(){
	chs:=make([]chan int,10)
	for i:=0;i<10;i++{
		chs[i]=make(chan int)
		go Count(chs[i],i)
	}
	//for cc,ch:=range(chs){
	//		fmt.Printf("Count[%d]处理完毕，对应的chs[%d]是%d\n",cc,cc,<-ch)
	//}
	for i:=0;i<10;i++{
		fmt.Printf("Count[%d]处理完毕，对应的chs[%d]是%d\n",i,i,<-chs[i])
	}
}
*/

//原书代码有个问题，client和server之间没有确实的通信，logout无效，我现在赶着去学校今天下午就不改了，先这样了
var centerClient *cg.CenterClient
func startCenterService()error{
	server:=ipc.NewIpcServer(&cg.CenterServer{})
	client:=ipc.NewIpcClient(server)
	centerClient=&cg.CenterClient{client}
	return nil
}
func Help(args []string)int{
	fmt.Printf("请输入以下指令：login<用户名><等级><exp>，logout<用户名>，send<信息>，listplayer，quit(q)，help(h)\n")
	return 0
}
func Quit(args []string)int{
	return 1
}
func Logout(args []string)int{
	if len(args)!=2{
		fmt.Printf("用法是：logout<用户名>\n")
		return 0
	}
	centerClient.RemovePlayer(args[1])
	return 0
}
func Login(args []string)int {
	if len(args) != 4 {
		fmt.Printf("用法是：login<用户名><等级><exp>\n")
		return 0
	}
	level, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("您输入的等级参数有误！检查一下是否是整数！\n")
		return 0
	}
	exp, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Printf("您输入的exp参数有误！检查一下是否是整数！\n")
		return 0
	}
	player := cg.NewPlayer()
	player.Name = args[1]
	player.Level = level
	player.Exp=exp
	err=centerClient.AddPlayer(player)
	if err!=nil{
		fmt.Printf("发生错误！无法添加player！%v\n",err)
	}
	return 0
}
func ListPlayer(args []string)int{
	ps,err:=centerClient.ListPlayer("")
	if err!=nil{
		fmt.Printf("中心客户端构建ListPlayer失败！%v\n",err)
	}else{
		for i,v:=range ps{
			fmt.Printf("player编号下标%v，用户名%v，等级%v，exp %v\n",i,v.Name,v.Level,v.Exp)
		}
	}
	return 0
}
func Send(args []string)int{
	message:=strings.Join(args[1:],"")
	err:=centerClient.Broadcast(message)
	if err!=nil{
		fmt.Printf("广播失败！%v\n",err)
	}
	return 0
}
func GetCommandHandlers()map[string]func(args []string)int{
	return map[string]func([]string)int{//这个map是 字符串---函数 的map，这里要注意一下
		"help":Help,//注意，这里字符串对应的是函数名！
		"h":Help,
		"quit":Quit,
		"q":Quit,
		"login":Login,
		"logout":Logout,
		"listplayer":ListPlayer,
		"send":Send,
	}
}

func main(){
	fmt.Printf("随意游戏服务器解决方案\n")
	startCenterService()
	Help(nil)
	r:=bufio.NewReader(os.Stdin)
	handlers:=GetCommandHandlers()
	for{
		fmt.Print("请输入指令：")
		b,_,_:=r.ReadLine()
		line:=string(b)
		tokens:=strings.Split(line," ")
		if handler,ok:=handlers[tokens[0]];ok{//根据那个map的定义，这里的handler是个函数，返回值是int
			ret:=handler(tokens)
			if ret!=0{
				break
			}
		}else{
			fmt.Printf("无法识别的指令！%v\n",tokens[0])
		}
	}
}