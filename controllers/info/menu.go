package info

import (
	"encoding/json"
	"fmt"
	"gisa/controllers"
	"gisa/models"
	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
)

type MenuController struct {
	controllers.BaseController
}
type MenuOrder struct {
	Id       int          `json:"id"`
	Children []*MenuOrder `json:"children"`
}

type MenuOrders []MenuOrder

//@desc 菜单管理列表页
//@router /info/menu [get]
func (c *MenuController) Index() {
	//c.Ctx.ResponseWriter.Header().Add("X-PJAX-Url", "/info/menu")
	c.Data["MenuNestableHtml"] = models.GetMenuNestableHtml()
	c.Data["treeId"] = string(utils.RandomCreateBytes(10))
	c.Data["MenuSelectOption"] = models.GetMenuSelectOption(0)
	c.Data["PageTitle"] = "菜单列表"
	c.ShowHtml("info/menu/index.html")
}

type Menu struct {
	Title    string `form:"title"`
	ParentId int    `form:"parent_id"`
	Icon     string `form:"icon"`
	Uri      string `form:"uri"`
}

//@desc 菜单添加
//@router /info/menu/add [post]
func (c *MenuController) Add() {
	menu := models.Menu{}
	if err := c.ParseForm(&menu); err == nil {
		menu.Order = 999
		if _, err := menu.Insert(); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
	c.Redirect(c.URLFor("MenuController.Index"), 302)
}

//@desc 菜单删除
//@router /info/menu/delete [post]
func (c *MenuController) Delete() {
	data := controllers.JsonData{}

	menu_id, err := c.GetInt("menu_id")
	if err != nil {
		data.Code = 400
		data.Message = "数据获取失败"
		c.ShowJSON(&data)
	}

	menuModel := models.Menu{}
	if err := menuModel.FindById(menu_id); err != nil {
		data.Code = 400
		data.Message = "数据获取失败"
		c.ShowJSON(&data)
	}

	if isDelete, err := menuModel.Delete(); isDelete {
		data.Code = 200
		data.Message = "删除成功"
		c.ShowJSON(&data)
	} else {
		data.Code = 400
		data.Message = err.Error()
		c.ShowJSON(&data)
	}

	c.ShowJSON(&data)
}

//@desc 菜单排序管理
//@router /info/menu/order [post]
func (c *MenuController) SaveMenuOrder() {
	var (
		_orderStr string
		mo        []*MenuOrder
		err       error
		dfs       func([]*MenuOrder, int)
	)
	_orderStr = c.Input().Get("_order")
	if err = json.Unmarshal([]byte(_orderStr), &mo); err != nil {
		fmt.Println(err.Error())
	}

	dfs = func(mo []*MenuOrder, parentId int) {
		for k, v := range mo {
			models.UpdateMenuParentIdAndOrder(v.Id, parentId, int(k+1))
			if v.Children != nil && len(v.Children) > 0 {
				dfs(v.Children, v.Id)
			}
		}

	}
	dfs(mo, 0)
	c.ServeJSON()
}

//@desc 修改菜单页面
//@router /info/menu/update/:menu_id [get,post,put]
func (c *MenuController) Update(menu_id int) {
	if c.IsPost() {
		c.DoUpdate()
	}
	var (
		err error
	)
	if menu_id == 0 {
		c.RedirectMessage(c.URLFor("MenuController.Index"), "数据不存在", controllers.MESSAGE_TYPE_ERROR)
	}

	menuModel := models.Menu{}
	if err = menuModel.FindById(menu_id); err != nil {
		c.RedirectMessage(c.URLFor("MenuController.Index"), "数据获取失败", controllers.MESSAGE_TYPE_ERROR)
	}
	c.Data["menu_model"] = menuModel
	c.Data["menu_model_icon"] = `<button class="btn btn-default" id="icon" data-icon="` + menuModel.Icon + `" name="icon" role="iconpicker"></button>`
	fmt.Println(c.Data["menu_model_icon"])
	c.Data["MenuSelectOption"] = models.GetMenuSelectOption(menuModel.ParentId)
	c.Data["role_select_list"] = models.AllRoleSelectList()
	c.Data["role_selected_list"] = models.AllRoleSelectedList(menuModel.Id)
	//c.Data["routes"] = selectRoutes
	//c.Data["parents"] = selectParents
	c.AddBreadcrumbs("菜单管理", c.URLFor("MenuController.Index"))
	c.AddBreadcrumbs("修改", c.URLFor("MenuController.DoUpdateMenu", "menu_id", menu_id))
	c.ShowHtml("info/menu/update.html")
}

//@desc 菜单编辑
//@router /info/menu/doupdate [post,put]
func (c *MenuController) DoUpdate() {
	var (
		o   orm.Ormer
		oM  *models.Menu
		err error
	)
	roleIds := c.GetStrings("roles[]")
	menu := models.Menu{}
	if err := c.ParseForm(&menu); err != nil {
		fmt.Println("parse form:", err.Error())
		return
	}
	o = orm.NewOrm()
	if menu.Id > 0 {
		oM, err = models.MenuOne(menu.Id)
		oM.ParentId = menu.ParentId
		oM.Title = menu.Title
		oM.Icon = menu.Icon
		oM.Uri = menu.Uri
		fmt.Println("om:", oM, menu)
		if _, err = o.Update(oM); err != nil {
			fmt.Println("orm update:", err.Error())
			return
		}
		if err = oM.DeleteRoles(); err != nil {
			fmt.Println("orm delete:", err.Error())
		}
		for _, v := range roleIds {
			roleId, _ := strconv.Atoi(v)
			if _, err = oM.AddRole(roleId); err != nil {
				fmt.Println("role add:", err.Error())
			}
		}
	}
}
