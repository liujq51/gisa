package controllers

import (
	"gisa/backend/logic"
	"gisa/backend/models"
	"strings"

	"github.com/lhtzbj12/sdrms/utils"
)

type HomeController struct {
	BaseController
}

//@desc 首页
//@router / [get]
func (c *HomeController) Get() {
	c.ShowHtml("index.html")
}

//@desc Dashboard
//@router /dashboard [get]
func (c *HomeController) Dashboard() {
	c.Data["content"] = "dashboard"
	c.ShowHtml("dashboard.html")
	//c.Ctx.WriteString("Dashboard.")
}

//@desc Login
//@router /login [get]
func (c *HomeController) Login() {
	c.Layout = ""
	//c.Data["content"] = "dashboard"
	//c.TplName = "dashboard.tpl"
	c.ShowHtml("home/login.html")
}

//@desc Login
//@router /login [post]
func (c *HomeController) DoLogin() {
	//c.Data["content"] = "dashboard"
	//c.TplName = "dashboard.tpl"
	username := strings.TrimSpace(c.GetString("username"))
	password := strings.TrimSpace(c.GetString("password"))
	if len(username) == 0 || len(password) == 0 {
		c.jsonResult(logic.JRCodeFailed, "用户名和密码不正确", "")
	}
	password = utils.String2md5(password)
	user, err := models.UserOneByUserName(username, password)
	if user != nil && err == nil {
		if user.Status == logic.Disabled {
			c.jsonResult(logic.JRCodeFailed, "用户被禁用，请联系管理员", "")
		}
		//保存用户信息到session
		c.setUserSession(user.Id)
		//获取用户信息
		//c.jsonResult(services.JRCodeSucc, "登录成功", "")
		c.Ctx.Redirect(302, "/")
	} else {
		c.jsonResult(logic.JRCodeFailed, "用户名或者密码错误", "")
	}
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
