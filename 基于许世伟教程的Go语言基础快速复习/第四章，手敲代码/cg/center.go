package cg
import(
	"../ipc"
	"sync"
	"encoding/json"
	"errors"
	"fmt"
)
type Message struct{
	From string "from"
	To string "to"
	Content string "content"
}
type CenterServer struct{
	Servers map[string]ipc.Server
	Players []*Player
	//rooms []*Room
	Mutex sync.RWMutex
}
var _ ipc.Server=&CenterServer{}
func NewCenterServer()*CenterServer{
	servers:=make(map[string]ipc.Server)
	players:=make([]*Player,0)
	return &CenterServer{Servers:servers,Players:players}
}
func(server *CenterServer)addPlayer(params string)error{
	player:=NewPlayer()
	err:=json.Unmarshal([]byte(params),&player)
	if err!=nil{
		return err
	}
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	server.Players=append(server.Players,player)
	return nil
}
func(server *CenterServer)removePlayer(params string)error{
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	for i,v:=range server.Players{
		if v.Name==params{//这个地方原书代码写得太蠢了，一句话搞定的事情写了四句还没写对
			server.Players=append(server.Players[:i],server.Players[i+1:]...)
			fmt.Printf("删除%v曲目成功！\n",params)
			return nil
		}
	}
	return errors.New("您输入的曲目名称没有找到，无法执行删除操作")
}
func(server *CenterServer)listPlayer(params string)(players string,err error){
	server.Mutex.RLock()
	defer server.Mutex.RUnlock()
	if len(server.Players)>0{
		b,_:=json.Marshal(server.Players)
		players=string(b)
	}else{
		err=errors.New("无法播放")
	}
	return
}
func(server *CenterServer)broadcast(params string)error{
	var message Message
	err:=json.Unmarshal([]byte(params),&message)
	if err!=nil{
		return err
	}
	server.Mutex.Lock()
	defer server.Mutex.Unlock()
	if len(server.Players)>0{
		for _,player:=range server.Players{
			player.mq<-&message
		}
	}else{
		err=errors.New("无法广播！")
	}
	return err
}
func(server *CenterServer)Name()string{
	return "名字是中心化服务器！"
}
func(server *CenterServer)Handle(method string,params string)*ipc.Response{
	switch method{
	case "addplayer":
		err:=server.addPlayer(params)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
	case "removeplayer":
		err:=server.removePlayer(params)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
	case "listplayer":
		players,err:=server.listPlayer(params)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response{"200",players}
	case "broadcast":
		err:=server.broadcast(params)
		if err!=nil{
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response{Code:"200"}
	default:
		return &ipc.Response{"404",method+":"+"params"}
	}
	return &ipc.Response{Code:"200"}
}