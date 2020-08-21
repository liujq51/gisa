package main

import (
	_ "gisa/routers"
	_ "gisa/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
