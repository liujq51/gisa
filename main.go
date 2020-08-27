package main

import (
	"fmt"
	_ "gisa/routers"
	_ "gisa/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	beego.Run()
}
