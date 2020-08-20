package controllers

type AuthController struct {
	BaseController
}

//@desc 菜单首页
//@router /auth/menu/list [get]
func (c *AuthController) MenuList() {
	//c.Data["Menus"] = c.getMenu()
	c.TplName = "auth/menu/index.tpl"
}
