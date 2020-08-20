package main

import (
	_ "gisa/backend/routers"
	_ "gisa/backend/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
