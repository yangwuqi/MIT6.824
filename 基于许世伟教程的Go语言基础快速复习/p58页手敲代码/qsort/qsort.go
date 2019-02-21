package qsort

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
