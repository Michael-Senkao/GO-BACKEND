package main

import "fmt"

func calcSum(arr []int) int{
	sum := 0
	for _, val := range arr{
		sum += val
	}

	return sum
}
func main(){
	fmt.Println("Hello, World!")
}
