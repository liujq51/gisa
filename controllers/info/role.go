package info

import (
	"fmt"
	"gisa/controllers"
	"gisa/models"
)

type RoleController struct {
	controllers.BaseController
}

//@desc 首页
//@router /info/menu [get]
func (this *RoleController) Index() {
	var (
		err error
	)

	IndexParams := models.JobListParams{}
	if err := this.ParseForm(&IndexParams); err != nil {
		//handle error
	}
	modles, pagination, err := models.ListRoles(IndexParams)
	if err != nil {
		fmt.Println(err.Error())
	}

	this.Data["modles"] = modles
	this.Data["pages"] = pagination
	this.AddBreadcrumbs("Role 管理", this.URLFor("RoleController.Index"))
	this.Data["PageTitle"] = "角色列表"
	this.ShowHtml("info/role/index.html")
}
