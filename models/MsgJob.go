package models

import (
	"errors"
	"fmt"
	"gisa/common/crontab"
	"gisa/logic"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

// Job is menu model structure.
type MsgJob struct {
	BaseModel
	Id    int    `form:"job_id"`
	Title string `form:"title"`
	//WebhookId     int    `form:"webhook"`
	MsgType       string `form:"msg_type"`
	MsgContent    string `form:"msg_content"`
	MsgTitle      string `form:"msg_title"`
	MsgDesc       string `form:"msg_desc"`
	MentionedList string `form:"mentioned_list"`
	MsgUrl        string `form:"msg_url"`
	MsgPicurl     string `form:"msg_picurl"`
	CronExpr      string `form:"cron_expr"`
	CreatedAt     string `form:"-"`
	UpdatedAt     string `form:"-"`
	Webhook       *Bot   `form:"webhook" orm:"rel(one)"`
}

// TableName 设置表名
func (p *MsgJob) TableName() string {
	return MsgJobTBName()
}

//list permissions
func (this MsgJob) ListJobs(pageIndex, pageCount int) ([]*MsgJob, int, error) {
	var Jobs []*MsgJob
	var total int64
	o := orm.NewOrm()
	_, err := o.QueryTable("gisa_msg_job").Limit(pageCount).Offset(pageCount * (pageIndex - 1)).RelatedSel().All(&Jobs)

	if err != nil {
		return Jobs, int(total), err
	}

	total, err = o.QueryTable(JobTBName()).Count()
	return Jobs, int(total), err
}

func (this *MsgJob) Insert() (isInsert bool, err error) {
	if this.Title == "" {
		return false, errors.New("名称不能为空")
	}

	this.CreatedAt = time.Now().Format("2020-01-01 00:00:00")
	this.UpdatedAt = this.CreatedAt
	o := orm.NewOrm()
	id, err := o.Insert(this)
	if id > 0 {
		SaveEtcdJob(int(id))
	}
	return id > 0, err
}

//保存etcd任务记录
func (this MsgJob) SaveEtcdJob(jobId int) bool {
	var (
		etcdJob crontab.Job
		o       orm.Ormer
		Job     MsgJob
		err     error
	)
	// 执行 SQL 语句
	o = orm.NewOrm()
	if err := o.QueryTable(MsgJobTBName()).Filter("id", jobId).RelatedSel().One(&Job); err != nil {
		fmt.Println(err.Error())
	}
	jobCommand := this.makeJobCommand(Job)
	etcdJob = crontab.Job{
		Name:     "msg_job_" + strconv.Itoa(jobId),
		Command:  jobCommand,
		CronExpr: Job.CronExpr,
	}
	if _, err = logic.G_jobMgr.SaveJob(&etcdJob); err != nil {
		fmt.Println(err.Error())
	}
	return true
}
func (this MsgJob) makeJobCommand(job MsgJob) (jobCommand string) {
	jobCommand = ""
	if job.MsgType == "text" {
		jobCommand = fmt.Sprintf("curl '%v' -H 'Content-Type: application/json' -d '{\"msgtype\":\"text\",\"text\":{\"content\":\"%v\",\"mentioned_list\":[%v],}}'",
			job.Webhook.Webhook, job.MsgContent, job.MentionedList)
	} else if job.MsgType == "markdown" {
		jobCommand = fmt.Sprintf("curl '%v' -H 'Content-Type: application/json' -d '{\"msgtype\":\"markdown\",\"markdown\":{\"content\":\"%v\"}}'",
			job.Webhook.Webhook, job.MsgContent)
	} else if job.MsgType == "news" {
		jobCommand = fmt.Sprintf("curl '%v' -H 'Content-Type: application/json' -d '{\"msgtype\":\"news\",\"news\":{\"title\":\"%v\",\"description\":\"%v\",\"url\":\"%v\",\"picurl\":\"%v\"}}'",
			job.Webhook.Webhook, job.MsgType, job.MsgDesc, job.MsgUrl, job.MsgPicurl)
	}
	return jobCommand
}

// 获取单条Job
func FindOneMsgJob(id int) (*MsgJob, error) {
	o := orm.NewOrm()
	this := MsgJob{Id: id}
	err := o.Read(&this)
	if err != nil {
		return nil, err
	}
	return &this, nil
}

//remove current menu from database
func (this *MsgJob) Delete() (isDelete bool, err error) {
	if this.Id <= 0 {
		return false, errors.New("删除对象不能为空")
	}
	jobName := "job_" + strconv.Itoa(this.Id)
	o := orm.NewOrm()
	num, err := o.Delete(this)
	if num > 0 {
		if _, err = logic.G_jobMgr.DeleteJob(jobName); err != nil {
			fmt.Println(err.Error())
		}
	}
	return num > 0, err
}

func (this *MsgJob) Update() (isUpdate bool, err error) {
	var num int64

	if this.Title == "" {
		return false, errors.New("名称不能为空")
	}

	this.UpdatedAt = time.Now().Format("2001-01-01 00:00:00")

	o := orm.NewOrm()
	if num, err = o.Update(this); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("num:", num, this.Id)
	this.SaveEtcdJob(this.Id)
	return num > 0, err
}
