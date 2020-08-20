package routers

import (
	"gisa/backend/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Include(&controllers.HomeController{})

	//beego.Include(&controllers.TestController{})
	//beego.Include(&controllers.LogController{})
	//beego.Include(&controllers.AuthController{})
	//beego.Include(&controllers.UserController{})
	//beego.Include(&info.MenuController{})
	//beego.Include(&info.PermissionController{})
	//beego.Include(&controllers.JobController{})

	//beego.ErrorController(&controllers.ErrorController{})

}
