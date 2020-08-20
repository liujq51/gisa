package controllers

type LogController struct {
	BaseController
}

//@desc 获取日志列表
//@router /log/list [get]
func (c *LogController) List() {
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.Ctx.WriteString(c.Layout)
	c.TplName = "log/list.tpl"
	//c.Ctx.WriteString("hi")
}
