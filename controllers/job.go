package controllers

import (
	"fmt"
	"gisa/common"
	"gisa/logic"
	"gisa/models"
)

type JobController struct {
	BaseController
	jobMgr *logic.JobMgr
}

func (this *JobController) Prepare() {
	this.BaseController.Prepare()
	if err := logic.InitJobMgr(); err != nil {
		fmt.Println(err.Error())
	}
}

//@desc 菜单首页
//@router /job [get]
func (this *JobController) Index() {
	//var (
	//	jobList []*crontab.Job
	//	err     error
	//)
	////  任务管理器
	//if err = logic.InitJobMgr(); err != nil {
	//	fmt.Println(err.Error())
	//}
	//// 获取任务列表
	//if jobList, err = logic.G_jobMgr.ListJobs(); err != nil {
	//	fmt.Println(jobList)
	//}
	//this.Data["job_list"] = jobList

	page_index, err := this.GetInt("page_index")
	pagination := common.Pagination{}
	pagination.PageCount = 20
	pagination.Url = this.URLFor("JobController.Index")

	if err == nil {
		pagination.PageIndex = page_index
	} else {
		pagination.PageIndex = 1
	}

	jobs, pageTotal, err := models.ListJobs(pagination.PageIndex, pagination.PageCount)
	if err != nil {
		//this.AddErrorMessage(err.Error())
		fmt.Println(err.Error())
	}

	pagination.PageTotal = pageTotal
	this.Data["modles"] = jobs
	this.Data["pages"] = pagination
	this.AddBreadcrumbs("Job 管理", this.URLFor("BotController.Index"))
	this.Data["PageTitle"] = "Job列表"
	this.ShowHtml("job/index.html")
}

//@desc 任务添加
//@router /job/add [get,post]
func (this *JobController) Add() {
	if this.IsPost() {
		this.Save()
	}
	//c.Data["routes"] = selectRoutes
	this.Data["webhook_select_option"] = models.Bot{}.GetBotSelectOption(0)
	this.Data["msg_type_select_option"] = models.Bot{}.GetBotMsgTypeList("")
	this.AddBreadcrumbs("任务管理", this.URLFor("JobController.Index"))
	this.AddBreadcrumbs("新增", "")
	this.ShowHtml("job/add.html")
}

//@desc 任务添加
func (this *JobController) Save() {
	var (
		job models.Job
		err error
	)

	job = models.Job{}
	if err := this.ParseForm(&job); err != nil {
		fmt.Println(err.Error())
	}
	//fmt.Printf("%+v", job)
	//if _, err = logic.G_jobMgr.SaveJob(&job); err != nil {
	//	fmt.Println(err.Error())
	//}
	if _, err = job.Insert(); err != nil {
		fmt.Println(err.Error())
	}
	this.Redirect(this.URLFor("JobCrontroller.Get"), 302)
}

//@desc 删除任务接口
//@router /job/delete   [post,delete]
func (this *JobController) Delete() {
	var (
		err   error // interface{}
		data  JsonData
		jobId int
		job   models.Job
	)
	data = JsonData{}
	// 删除的任务名
	jobId, _ = this.GetInt("job_id")
	if err != nil {
		data.Code = 400
		data.Message = "数据获取失败"
		this.ShowJSON(&data)
	}
	job = models.Job{
		Id: jobId,
	}
	// 删除任务处理
	if _, err := job.Delete(); err != nil {
		fmt.Println(err.Error())
	}

	data.Code = 200
	data.Message = "success"
	this.ShowJSON(&data)
}

//@desc 修改页面
//@router /job/edit/:jobId [get,post,put]
func (this *JobController) Edit(jobId int) {
	if this.IsPost() {
		this.DoUpdate()
	}
	var (
		err error
		job *models.Job
	)
	if jobId == 0 {
		this.RedirectMessage(this.URLFor("MenuController.Index"), "数据不存在", MESSAGE_TYPE_ERROR)
	}

	// 获取任务信息
	if job, err = models.FindOneJob(jobId); err != nil {
		fmt.Println(err.Error())
	}
	this.Data["model"] = job
	this.AddBreadcrumbs("任务管理", this.URLFor("JobController.Index"))
	this.AddBreadcrumbs("修改", "")
	this.ShowHtml("job/edit.html")
}

func (this *JobController) DoUpdate() {
	var (
		job models.Job
		//err error
	)
	job = models.Job{}
	if err := this.ParseForm(&job); err != nil {
		fmt.Println(err.Error())
	}
	if _, err := job.Update(); err != nil {
		fmt.Println(err.Error())
	}
	//if _, err = logic.G_jobMgr.SaveJob(&job); err != nil {
	//	fmt.Println(err.Error())
	//}

	//this.Redirect(this.URLFor("JobController.Index"), 302)
}
