package controllers

type JobController struct {
	BaseController
}

//@desc 菜单首页
//@router /job [get]
func (this *JobController) JobList() {
	//c.Data["Menus"] = c.getMenu()
	this.ShowHtml("job/index.html")
}
