package library

import "errors"

type MusicEntry struct {//原书的这个结构体里的第一个元素是Id也就是歌曲编号，实际上原书的这个设定很有问题，
	Name string	//因为按编号删除非最后一个元素会导致编号混乱，直接用slice的元素下标作为编号就行了
	Artist string
	Style string
	Source string
	Type string
}
type MusicManager struct{
	musics []MusicEntry
}
func(m *MusicManager)Add(newmusic *MusicEntry){
	m.musics=append(m.musics,*newmusic)
}
func NewMusicManager() *MusicManager{
return &MusicManager{make([]MusicEntry,0)}
}
func (m *MusicManager)Len()int{
	return len(m.musics)
}
func(m*MusicManager)Get(index int)(music *MusicEntry,err error){
	if index<0||index>=len(m.musics){
		return nil,errors.New("出错啦，歌曲编号越界！")
	}
	return &m.musics[index],nil
}
func(m *MusicManager)Find(name string)*MusicEntry{
	if len(m.musics)==0{
		return nil
	}
	for _,x:=range m.musics{
		if x.Name==name{
			return &x
		}
	}
	return nil
}
func (m *MusicManager)Remove(index int)*MusicEntry{
	if index<0||index>len(m.musics)-1{
		return nil
	}
	temp:=m.musics[index]
	m.musics=append(m.musics[:index],m.musics[index+1:]...)//这个地方许世伟的书上的代码append里面第一个竟然是m.musics[:index-1]，他书上的代码错了
	return &temp
}
