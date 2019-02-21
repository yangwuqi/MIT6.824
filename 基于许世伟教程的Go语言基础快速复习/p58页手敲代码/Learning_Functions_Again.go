package main

import (
	"fmt"
	"errors"
	"os"
	"io"
	"flag"
	"bufio"
	"strconv"
	"./qsort"
	"time"
	"./bubblesort"
)

func Add(a int,b int)(ret int,err error){
	if a<0||b<0{
		fmt.Printf("出错啦!!\n")
		err=errors.New("出错啦!a或者b小于0!")
		return
	}
	return a+b,nil
}

func myfunc(args ...int){
	for _,arg:=range args{
		fmt.Println(arg)
	}
}


func main(){
fmt.Printf("你好呀\n")


f:=func(x,y int)int{return x+y}(2,5)
fmt.Printf("匿名相加函数的参数是2和5，结果是f:%v\n",f)
//下面是闭包和匿名函数示例，作用是保护私有区域i的值
j:=5
a:=func()(func()){
	var i int=10
	return func(){
		fmt.Printf("i*j:%d*%d=%d\n",i,j,i*j)
	}
}()
a()
j*=2
a()
fmt.Printf("\nGO的匿名函数和闭包演示结束，下面是读写文件和简单排序的演示\n\n")

var infile *string=flag.String("i","infile","用于排序的文件")
var outfile *string=flag.String("o","outfile","用于接受排序值们的文件")
var algorithm *string=flag.String("a","qsort","排序算法")
flag.Parse()//用了这个flag才运行了，将输入的命令行参数处理了，上面flag string捕捉-i或者-o之类的命令后面的字符串，返回其地址给infile或者outfile
if infile!=nil{
	fmt.Println("infile=",*infile,"outfile=",*outfile,"algorithm=",*algorithm)
	}
values,err:=readValues(*infile)
if err==nil{
	fmt.Printf("读取数值们:%v",values)
	t1:=time.Now()
	switch *algorithm{
	case "qsort":
		fmt.Printf("启动快速排序\n")
		qsort.QuickSort(values)
	case "bubblesort":
		fmt.Printf("启动冒泡排序\n")
		bubblesort.BubbleSort(values)
	default:
		fmt.Printf("你输入的第三个参数也就是算法名称是未知的，这里只支持快排和冒泡\n")
	}
	t2:=time.Now()
	fmt.Printf("排序过程耗时%v, 排序结果是%v\n",t2.Sub(t1),values)
	writeValues(values,*outfile)
}else{
	fmt.Printf("%v\n",err)
}
}

func readValues(infile string)(values []int,err error){
	file,err:=os.Open(infile)//打开文件
	if err!=nil{
		fmt.Printf("出错啦！\n")
		return
	}
	defer file.Close()//记得最后要关闭文件
	br:=bufio.NewReader(file)//将打开的文件读出来准备处理，这个处理对象是br
	values=make([]int,0)
	for {
		line,isPrefix,err1:=br.ReadLine()//这里注意，bufio.NewReader().ReadLine()是按\n截取一行，返回三个值，
		if err1!=nil {					 //第一个[]byte类型的文件一行内容，第二个是否\n结尾的bool，第三个error
			if err1 != io.EOF {			 //有时，比如缓冲区满了这一行末尾没有\n，isPrefix就是true，否则是false
				err = err1
			}
			break
		}
		if isPrefix{
			fmt.Printf("太长了，看起来处理不来\n")
			return
		}
		str:=string(line)
		value,err1:=strconv.Atoi(str)
		if err1!=nil{
			err=err1
			return
		}
		values=append(values,value)
	}
	return
}

func writeValues(values []int, outfile string)error{
	file,err:=os.Create(outfile)
	if err!=nil{
		fmt.Println("创建输出文件失败啦！",outfile)
		return err
	}
	defer file.Close()
	for _,value:=range values{
		str:=strconv.Itoa(value)
		fmt.Printf("写入的这行内容%v\n",str)//在windows下，str+"\n"在记事本里是不显示换行的，得写成str+"\r\n"
		file.WriteString(str+"\r\n")//os.Create("文件名").WriteString(str+"\n")，表示创建一个名字为 文件名 的文件，然后往里面写一行str字符串
	}
	return nil
}

