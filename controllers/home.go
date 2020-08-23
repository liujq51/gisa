package controllers

import (
	"gisa/logic"
	"gisa/models"
	"gisa/utils"
	"strings"

	"github.com/astaxie/beego"
)

type HomeController struct {
	BaseController
}

//@desc 首页
//@router / [get]
func (this *HomeController) Get() {
	this.ShowHtml("index.html")
}

//@desc Dashboard
//@router /dashboard [get]
func (this *HomeController) Dashboard() {
	this.Data["content"] = "dashboard"
	this.ShowHtml("dashboard.html")
}

//@desc Login
//@router /login [get]
func (c *HomeController) Login() {
	beego.ReadFromRequest(&c.Controller)
	c.Layout = ""
	//c.TplName = "dashboard.tpl"
	c.ShowHtml("home/login.html")
}

//@desc Login
//@router /login [post]
func (c *HomeController) DoLogin() {
	flash := beego.NewFlash()
	username := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))
	if len(username) == 0 || len(password) == 0 {
		//c.jsonResult(logic.JRCodeFailed, "用户名和密码不正确", "")
		flash.Error("用户名和密码不正确!")
	}
	flash.Error("用户名或者密码错误!")
	password = utils.String2md5(password)
	user, err := models.UserOneByUserName(username, password)
	if user != nil && err == nil {
		if user.Status == logic.Disabled {
			//c.jsonResult(logic.JRCodeFailed, "用户被禁用，请联系管理员", "")
			flash.Error("用户被禁用，请联系管理员!")
		} else {
			//保存用户信息到session
			c.setUserSession(user.Id)
			//获取用户信息
			//c.jsonResult(services.JRCodeSucc, "登录成功", "")
			c.Ctx.Redirect(302, "/")
		}
	} else {
		//c.jsonResult(logic.JRCodeFailed, "用户名或者密码错误", "")
		flash.Error("用户名或者密码错误!")
	}
	flash.Error("用户名或者密码错误!")
	flash.Store(&c.Controller)
	c.Redirect("/login", 302)
}

//@desc logout
//@router /logout [*]
func (c *HomeController) Logout() {
	user := models.User{}
	c.SetSession("user", user)
	url := c.URLFor("HomeController.Login")
	c.Redirect(url, 302)
}

func (c *HomeController) error() {
	c.TplName = "error/error.tpl"
}
