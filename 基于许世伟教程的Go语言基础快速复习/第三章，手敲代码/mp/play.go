package mp

import "fmt"

type Player interface {
	Play(source string)
}
func Play(source,mtype string){
	var p Player
	switch mtype{
	case "Mp3":
		p=&Mp3Player{}

	case "WAV":
		//p=&WavPlayer{}
	default:
		fmt.Printf("未知的音乐类型！无法播放！类型为%v\n",mtype)
		return
	}
	p.Play(source)
}