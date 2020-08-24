package controllers

import (
	"fmt"
	"gisa/common"
	"gisa/logic"
	"gisa/models"
)

type MsgJobController struct {
	BaseController
	jobMgr *logic.JobMgr
}

func (this *MsgJobController) Prepare() {
	this.BaseController.Prepare()
	if err := logic.InitJobMgr(); err != nil {
		fmt.Println(err.Error())
	}
}

//@desc 菜单首页
//@router /msg_job [get]
func (this *MsgJobController) Index() {
	page_index, err := this.GetInt("page_index")
	pagination := common.Pagination{}
	pagination.PageCount = 20
	pagination.Url = this.URLFor("MsgJobController.Index")

	if err == nil {
		pagination.PageIndex = page_index
	} else {
		pagination.PageIndex = 1
	}

	jobs, pageTotal, err := models.MsgJob{}.ListJobs(pagination.PageIndex, pagination.PageCount)
	if err != nil {
		//this.AddErrorMessage(err.Error())
		fmt.Println(err.Error())
	}

	pagination.PageTotal = pageTotal
	this.Data["modles"] = jobs
	this.Data["pages"] = pagination
	this.AddBreadcrumbs("Msg Job 管理", this.URLFor("MsgJobController.Index"))
	this.Data["PageTitle"] = "Message Job列表"
	this.ShowHtml("msg_job/index.html")
}

//@desc 消息任务添加
//@router /msg_job/add [get,post]
func (this *MsgJobController) Add() {
	if this.IsPost() {
		this.Save()
	}
	//c.Data["routes"] = selectRoutes
	this.Data["webhook_select_option"] = models.Bot{}.GetBotSelectOption(0)
	this.Data["msg_type_select_option"] = models.Bot{}.GetBotMsgTypeList("")
	this.AddBreadcrumbs("任务管理", this.URLFor("MsgJobController.Index"))
	this.AddBreadcrumbs("新增", "")
	this.ShowHtml("msg_job/add.html")
}

//@desc 任务添加
func (this *MsgJobController) Save() {
	var (
		job models.MsgJob
		bot models.Bot
		err error
	)

	job = models.MsgJob{}
	if err := this.ParseForm(&job); err != nil {
		fmt.Println(err.Error())
	}
	botId, _ := this.GetInt("webhook")
	bot = models.Bot{
		Id: botId,
	}
	job.Webhook = &bot
	fmt.Printf("%+v", job)
	//if _, err = logic.G_jobMgr.SaveJob(&job); err != nil {
	//	fmt.Println(err.Error())
	//}
	if _, err = job.Insert(); err != nil {
		fmt.Println(err.Error())
	}
	this.Redirect(this.URLFor("JobCrontroller.Get"), 302)
}

//@desc 删除任务接口
//@router /msg_job/delete   [post,delete]
func (this *MsgJobController) Delete() {
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
//@router /msg_job/edit/:jobId [get,post,put]
func (this *MsgJobController) Edit(jobId int) {
	if this.IsPost() {
		this.DoUpdate()
	}
	var (
		err error
		job *models.MsgJob
	)
	if jobId == 0 {
		this.RedirectMessage(this.URLFor("MsgJobController.Index"), "数据不存在", MESSAGE_TYPE_ERROR)
	}

	// 获取任务信息
	if job, err = models.FindOneMsgJob(jobId); err != nil {
		fmt.Println(err.Error())
	}
	this.Data["model"] = job
	this.Data["webhook_select_option"] = models.Bot{}.GetBotSelectOption(job.Webhook.Id)
	this.Data["msg_type_select_option"] = models.Bot{}.GetBotMsgTypeList(job.MsgType)
	this.AddBreadcrumbs("任务管理", this.URLFor("MsgJobController.Index"))
	this.AddBreadcrumbs("修改", "")
	this.ShowHtml("msg_job/edit.html")
}

func (this *MsgJobController) DoUpdate() {
	var (
		job models.MsgJob
		//err error
		bot models.Bot
	)
	job = models.MsgJob{}
	if err := this.ParseForm(&job); err != nil {
		fmt.Println(err.Error())
	}
	botId, _ := this.GetInt("webhook")
	bot = models.Bot{
		Id: botId,
	}
	job.Webhook = &bot
	if _, err := job.Update(); err != nil {
		fmt.Println(err.Error())
	}
	//if _, err = logic.G_jobMgr.SaveJob(&job); err != nil {
	//	fmt.Println(err.Error())
	//}

	//this.Redirect(this.URLFor("MsgJobController.Index"), 302)
}
