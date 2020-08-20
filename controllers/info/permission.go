package info

import (
	"fmt"
	"gisa/backend/common"
	"gisa/backend/controllers"

	"gisa/backend/models"
)

type PermissionController struct {
	controllers.BaseController
}

//@desc 权限管理列表页
//@router /info/permission [get]
func (this *PermissionController) Index() {
	//c.Ctx.ResponseWriter.Header().Add("X-PJAX-Url", "/info/menu")
	page_index, err := this.GetInt("page_index")

	pagination := common.Pagination{}
	pagination.PageCount = 20
	pagination.Url = this.URLFor("PermissionController.Index")

	if err == nil {
		pagination.PageIndex = page_index
	} else {
		pagination.PageIndex = 1
	}

	permissions, pageTotal, err := models.Permission{}.List(pagination.PageIndex, pagination.PageCount)

	if err != nil {
		//this.AddErrorMessage(err.Error())
		fmt.Println(err.Error())
	}

	pagination.PageTotal = pageTotal
	this.Data["permissions"] = permissions
	this.Data["pages"] = pagination
	this.AddBreadcrumbs("权限管理", this.URLFor("PermissionController.Index"))
	this.Data["PageTitle"] = "权限列表"
	this.ShowHtml("info/permission/index.html")
}
