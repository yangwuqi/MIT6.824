package main

import (
	"./library"
	"fmt"
	"strconv"
	"./mp"
	"bufio"
	"os"
	"strings"
)
var lib *library.MusicManager
var ctrl,signal chan int
func handleLibCommands(tokens []string){
	switch tokens[1]{
	case "list":
		for i:=0;i<lib.Len();i++{
			e,_:=lib.Get(i)
			fmt.Printf("获取列表的第%d支曲目，名称为%v，作者为%v，类型为%v\n",i+1,e.Name,e.Artist,e.Type)
		}
	case "add":{
		if len(tokens)==7{
			lib.Add(&library.MusicEntry{tokens[2],tokens[3],tokens[4],tokens[5],tokens[6]})
		}else{
			fmt.Println("add指令的参数数量输入错误，应该是6个")
		}
	}
	case "remove":
		if len(tokens)==3{
			index,_:=strconv.Atoi(tokens[2])
			lib.Remove(index)
		}else{
			fmt.Printf("remove指令的参数输入数量错误，应该是两个")
		}
	default:
		fmt.Println("无法被识别的命令！%v\n",tokens[1])
	}
}

func handlePlayCommand(tokens []string){
	if len(tokens)!=2{
		fmt.Printf("出错了！播放的参数应该是2个，play+歌曲名\n")
		return
	}
	e:=lib.Find(tokens[1])
	if e==nil{
		fmt.Println("这个曲目不存在！")
		return
	}
	mp.Play(e.Source,e.Type)
}
func main(){
	fmt.Println("请输入指令来控制播放器，指令第一个是lib，则后面可以是 list, add+歌曲信息（5项）, remove+歌曲下标编号, 指令第一个是play，则后面是 歌曲名，输入e或者q退出")
	lib=library.NewMusicManager()
	r:=bufio.NewReader(os.Stdin)
	for{
		fmt.Print("请输入指令")
		rawline,_,_:=r.ReadLine()
		line:=string(rawline)
		if line=="q"||line=="e"{
			break
		}
		tokens:=strings.Split(line," ")
		if tokens[0]=="lib"{
			handleLibCommands(tokens)
		}else if tokens[0]=="play"{
			handlePlayCommand(tokens)
		}else{
			fmt.Printf("指令无法识别！%v\n",tokens[0])
		}
	}
}