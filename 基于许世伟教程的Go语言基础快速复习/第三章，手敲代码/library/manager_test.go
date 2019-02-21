package library

import (
	"testing"
)

func TestOps(t *testing.T){
	mm:=NewMusicManager()
	if mm==nil{
		t.Error("新建音乐管理程序失败！")
	}
	if mm.Len()!=0{
		t.Error("新建了音乐管理程序，但是竟然初始列表长度不为空！失败！")
	}
	m0:=&MusicEntry{
		"黄黄nb！",
		"鸡哥作曲",
		"Pop",
		"http://www.baidu.com",
		"Mp3",
	}
	mm.Add(m0)
	if mm.Len()!=1{
		t.Error("新增第一首歌失败！")
	}
	m:=mm.Find(m0.Name)
	if m==nil{
		t.Error("音乐管理程序按名称查找第一首歌失败！")
	}
	m,_=mm.Get(0)
	if m==nil{
		t.Error("音乐管理程序的按索引获取曲面功能失败！")
	}
	m=mm.Remove(0)
	if m==nil||mm.Len()!=0{
		t.Error("音乐管理程序的删除曲目功能失败！")
	}
}