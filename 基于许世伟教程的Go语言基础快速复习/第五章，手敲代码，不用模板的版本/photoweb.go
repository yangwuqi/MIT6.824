package main

import (
	"net/http"
	"io"
	"log"
	"os"
	"io/ioutil"
	"fmt"
)

/*
func helloHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"HELLO WORLD!鸡哥鸡哥你在干嘛？")
}
func main(){
	http.HandleFunc("/hey",helloHandler)//启动http服务器，回调函数是helloHandler，注意本地服务器ip地址是127.0.0.1
	err:=http.ListenAndServe(":9090",nil)//为这个http服务器监听:9090端口，在浏览器输入127.0.0.1:9090/hey打开
	if err!=nil{
		log.Fatal("ListenAndServe:",err.Error())
	}
}
*/

const(
	UPLOAD_DIR="./uploads"
)
func uploadHandler(w http.ResponseWriter,r *http.Request){
	if r.Method=="GET"{
		io.WriteString(w,"<html><form method=\"POST\" action=\"/upload\" "+" enctype=\"multipart/form-data\">"+"Choose an image to upload: <input name=\"image\" type=\"file\" />"+"<input type=\"submit\" value=\"Upload\" />"+"</form></html>")
		return
	}
	if r.Method=="POST"{
		f,h,err:=r.FormFile("image")
		if err!=nil{
			http.Error(w,err.Error(), http.StatusInternalServerError)
			return
		}
		filename:=h.Filename
		defer f.Close()
		fmt.Printf("即将创建上传的文件，创建路径为%v/%v\n",UPLOAD_DIR,filename)
		t,err:=os.Create(UPLOAD_DIR+"/"+filename)//注意路径这里在最前面的/前面必须加一个.
		if err!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _,err:=io.Copy(t,f);err!=nil{
			http.Error(w,err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w,r,"/view?id="+filename,http.StatusFound)
	}
}
func viewHandler(w http.ResponseWriter,r *http.Request){
	imageId:=r.FormValue("id")//这里原书代码有问题吧，原书是=不是:=
	imagePath:=UPLOAD_DIR+"/"+imageId//这里原书代码有问题吧，原书是=不是:=
	if exists:=isExits(imagePath);!exists{
		http.NotFound(w,r)
		return
	}
	w.Header().Set("Content-Type","image")
	http.ServeFile(w,r,imagePath)
}
func isExits(path string)bool{
	_,err:=os.Stat(path)
	if err==nil{
		return true
	}
	return os.IsExist(err)
}
func listHandler(w http.ResponseWriter,r *http.Request){
	fileInfoArr,err:=ioutil.ReadDir("./uploads")
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	var listHtml string
	for _,fileInfo:=range fileInfoArr{
		listHtml+="<li><a href=\"/view?id="+fileInfo.Name()+"\">图片</a></li>"//原书代码这里这句写错了fileInfo没()，加上()就好了
	}
	io.WriteString(w,"<html><body><ol>"+listHtml+"</ol></body></html>")//原书这里没有<HTML></HTML>包裹，Chrome无法识别成网页
}
func main(){
	http.HandleFunc("/",listHandler)
	http.HandleFunc("/view",viewHandler)
	http.HandleFunc("/upload",uploadHandler)
	err:=http.ListenAndServe(":8080",nil)
	if err!=nil{
		log.Fatal("ListenAndServe: ",err.Error())
	}
}