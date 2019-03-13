//顺便自己手写了一下堆排序

func main() {
    array:=[]int{3,9,5,7,1,0,2,3}
    heap(array)
    fmt.Printf("heap sort! %v\n",array)
}

func heap(nums []int){
    for i:=len(nums)/2-1;i>-1;i--{
        maxheap(nums,i,len(nums)-1)
    }//建初始堆，时间复杂度是O(N)
    for i:=len(nums)-1;i>-1;i--{//开始重复调整，每调整一次时间复杂度是O(logN)
        nums[0],nums[i]=nums[i],nums[0]
        maxheap(nums,0,i-1)
    }
}

func maxheap(nums []int, start,end int){
    father:=start
    child:=2*father+1
    for child<=end{
        if child+1<=end&&nums[child]<nums[child+1]{
            child++
        }
        if nums[father]>=nums[child]{
            return
        }else{
            nums[father],nums[child]=nums[child],nums[father]
            father=child
            child=2*father+1
        }
    }
}
