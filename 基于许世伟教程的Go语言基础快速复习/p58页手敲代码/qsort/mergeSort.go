//自己顺便写的归并排序

func main() {
    array:=[]int{3,9,5,7,1,0,2,3}
    fmt.Printf("merge sort! %v\n",mergesort(array))
}

func mergesort(array []int)[]int{
    if len(array)==1{
        return array
    }
    return merge(mergesort(array[:len(array)/2]),mergesort(array[len(array)/2:]))
}

func merge(x,y []int)[]int{
    i,j:=0,0
    var result []int
    for i<len(x)&&j<len(y){
        if x[i]<y[j]{
            result=append(result,x[i])
            i++
        }else{
            result=append(result,y[j])
            j++
        }
    }
    if i<len(x){
        result=append(result,x[i:]...)
    }
    if j<len(y){
        result=append(result,y[j:]...)
    }
    return result
}
