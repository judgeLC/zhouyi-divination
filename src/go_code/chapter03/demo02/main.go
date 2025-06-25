package main

import "fmt"

func main() {

	//指定变量类型，声明后不赋值，使用默认值
	var i int
	fmt.Println("i=", i)

	//根据值自行判断变量类型（类型推导）
	var num = 10.11
	fmt.Println("num=", num)

	//省略var，注意：+左侧的变量不应该是已经声明过的，否则会报错
	name := "tom"
	fmt.Println("name=", name)

}
