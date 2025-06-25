package main

import "fmt"

// 全局变量
var n1 = 100
var n2 = 200
var name = "jack"

// 可以优化称一次性申明
var (
	n3    = 300
	n4    = 900
	name2 = "mary"
)

func main() {

	//一次性申明多个变量
	//var n1, n2, n3 int
	//fmt.Println("n1=", n1, "n2=", n2, "n3=", n3)

	//一次性什么多个变量方式2
	//var n1, name, n3 = 1000, "tom", 888
	//fmt.Println("n1=", n1, "name=", name, "n3=", n3)

	//一次性申明多个变量的方式3，同样可以使用类型推导
	//n1, name, n3 := 1000, "tom~", 888
	//fmt.Println("n1=", n1, "name=", name, "n3=", n3)

	fmt.Println("n3=", n3, "name2=", name2, "n4=", n4)

}
