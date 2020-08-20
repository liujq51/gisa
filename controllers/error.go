package controllers

import "github.com/astaxie/beego"

type ErrorController struct {
	beego.Controller
}

//@router /error
func (c *ErrorController) Error(msg string) {
	c.Data["content"] = msg
	c.TplName = "error/error.tpl"
}

//@router /error/404
func (c *ErrorController) Error404() {
	c.Data["content"] = "page not found"
	c.TplName = "error/404.tpl"
}

func (c *ErrorController) Error401() {
	c.Data["content"] = "no nuthorization"
	c.TplName = "error/401.tpl"
}
