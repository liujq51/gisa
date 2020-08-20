package controllers

type TestController struct {
	BaseController
}

//@router /test [get]
func (c *TestController) Get() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.Ctx.WriteString(c.Layout)
	c.TplName = "test.html"
}
