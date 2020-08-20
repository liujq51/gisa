package controllers

type UserController struct {
	BaseController
}

//@desc 后台用户管理首页
//@router /admin/user [get]
func (c *UserController) UserList() {
	//c.Data["Menus"] = c.getMenu()
	c.TplName = "user/index.tpl"
}
