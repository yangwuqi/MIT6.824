package mp

import (
	"fmt"
	"time"
)

type Mp3Player struct {
	stat int
	progress int
}

func(p *Mp3Player)Play(source string){
	fmt.Printf("正在播放Mp3音乐，%v\n",source)
	p.progress=0
	for p.progress<500{
		time.Sleep(100*time.Millisecond)//假装正在播放
		fmt.Print(".")
		p.progress+=10
	}
	fmt.Printf("已完成播放，%v\n",source)
}
