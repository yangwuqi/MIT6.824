package qsort
/*
func quickSort(values []int,left,right int){
	temp:=values[left]
	p:=left
	i,j:=left,right
	for i<=j{
		for j>=p&&values[j]>=temp{//找values[j]<temp
			j--
		}
		if j>=p{//这时values[j]<temp且下标j后面的数全部大于等于temp，把values[p]也就是temp和values[j]交换
			values[p]=values[j]//先把values[p]的位置换成values[j]，原来的values[p]还在temp里存着，values[j]上面的先不管，把p设成j，后面还需要从左到p找比temp大的
			p=j//更新p的值为j，此时values[p]的值变成values[j]而不是之前的values[left]
		}
		if values[i]<=temp&&i<=p{//找values[i]>temp，请注意，从左到p找的是比temp大的，但是这时候p已经变成了上面的j而不是left，这时候下标j后面的已经全部大于等于temp，但是values[j]的值还没被改变
			i++
		}
		if i<=p{//这时values[i]>temp，且下标i前面的全部小于等于temp，把values[p]的值变成values[i]
			values[p]=values[i]
			p=i//更新p的值为i，此时values[p]=values[i]而不是之前的values[j]，之前values[j]的值已经存在之前的values[p]了，唯一没了的values[left]在temp里
		}
	}
	values[p]=temp//最后，把values[p]的值变成temp
	if p-left>1{
		quickSort(values,left,p-1)
	}
	if right-p>1{
		quickSort(values,p+1,right)
	}
}
*/
func QuickSort(values []int){
	quickSort(values,0,len(values)-1)
}

func quickSort(values []int,left,right int){
	if left>=right{
		return
	}
	mid:=values[left]
	i,j:=left,right
	for i<j{
		for j>i&&values[j]>=mid{
			j--
		}
		values[i],values[j]=values[j],values[i]
		for i<j&&values[i]<=mid{
			i++
		}
		values[i],values[j]=values[j],values[i]
	}
	quickSort(values,left,i-1)
	quickSort(values,i+1,right)
}
