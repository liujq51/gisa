package routers

import (
	"gisa/controllers"
	"gisa/controllers/info"

	"github.com/astaxie/beego"
)

func init() {

	beego.Include(&controllers.HomeController{})
	beego.Include(&controllers.TestController{})
	beego.Include(&controllers.LogController{})
	beego.Include(&controllers.AuthController{})
	beego.Include(&controllers.UserController{})
	beego.Include(&controllers.JobController{})
	beego.Include(&controllers.MsgJobController{})
	beego.Include(&info.MenuController{})
	beego.Include(&info.PermissionController{})
	beego.Include(&controllers.BotController{})
	//beego.ErrorController(&controllers.ErrorController{})
}
