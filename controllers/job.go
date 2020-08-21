package controllers

import (
	"fmt"
	"gisa/common/crontab"
	"gisa/logic"
)

type JobController struct {
	BaseController
	jobMgr *logic.JobMgr
}

//@desc 菜单首页
//@router /job [get]
func (this *JobController) Get() {
	var (
		jobList []*crontab.Job
		err     error
	)
	//  任务管理器
	if err = logic.InitJobMgr(); err != nil {
		fmt.Println(err.Error())
	}
	// 获取任务列表
	if jobList, err = this.jobMgr.ListJobs(); err != nil {
		fmt.Println(jobList)
	}
	fmt.Println(this.jobMgr)
	this.ShowHtml("job/index.html")
}

//@desc 任务添加
//@router /job/add [get]
func (this *JobController) Save() {

	var (
		job    crontab.Job
		oldJob *crontab.Job
		err    error
	)

	job = crontab.Job{
		Name:     "test",
		Command:  "ls -al",
		CronExpr: "* * * * * * *",
	}
	if err = logic.InitJobMgr(); err != nil {
		fmt.Println(err.Error())
	}

	if oldJob, err = logic.G_jobMgr.SaveJob(&job); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("old Job:", oldJob)

	data := JsonData{}
	data.Code = 200
	data.Message = "message"
	data.Data = map[string]interface{}{
		"jobs": oldJob,
	}
	this.ShowJSON(&data)
}
