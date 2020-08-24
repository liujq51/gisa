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
type Job struct {
	BaseModel
	Id        int    `form:"job_id"`
	Title     string `form:"title"`
	Command   string `form:"command"`
	CronExpr  string `form:"cron_expr"`
	CreatedAt string `form:"-"`
	UpdatedAt string `form:"-"`
}

// TableName 设置表名
func (this *Job) TableName() string {
	return JobTBName()
}

//list permissions
func ListJobs(pageIndex, pageCount int) ([]*Job, int, error) {
	var Jobs []*Job
	var total int64
	o := orm.NewOrm()
	_, err := o.QueryTable("gisa_job").Limit(pageCount).Offset(pageCount * (pageIndex - 1)).RelatedSel().All(&Jobs)

	if err != nil {
		return Jobs, int(total), err
	}

	total, err = o.QueryTable(JobTBName()).Count()
	return Jobs, int(total), err
}

func (this *Job) Insert() (isInsert bool, err error) {
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
func SaveEtcdJob(jobId int) bool {
	var (
		etcdJob crontab.Job
		o       orm.Ormer
		Job     Job
		err     error
	)
	// 执行 SQL 语句
	o = orm.NewOrm()
	if err := o.QueryTable(JobTBName()).Filter("id", jobId).RelatedSel().One(&Job); err != nil {
		fmt.Println(err.Error())
	}

	etcdJob = crontab.Job{
		Name:     "job_" + strconv.Itoa(jobId),
		Command:  Job.Command,
		CronExpr: Job.CronExpr,
	}
	if _, err = logic.G_jobMgr.SaveJob(&etcdJob); err != nil {
		fmt.Println(err.Error())
	}
	return true
}

// 获取单条Job
func FindOneJob(id int) (*Job, error) {
	o := orm.NewOrm()
	this := Job{Id: id}
	err := o.Read(&this)
	if err != nil {
		return nil, err
	}
	return &this, nil
}

//remove current menu from database
func (this *Job) Delete() (isDelete bool, err error) {
	if this.Id <= 0 {
		return false, errors.New("删除对象不能为空")
	}
	jobName := "job_" + strconv.Itoa(this.Id)
	o := orm.NewOrm()
	num, err := o.Delete(this)
	fmt.Println("del job num:", num, "job_"+strconv.Itoa(this.Id), this.Id, jobName)
	if num > 0 {
		if _, err = logic.G_jobMgr.DeleteJob(jobName); err != nil {
			fmt.Println(err.Error())
		}
	}
	return num > 0, err
}

func (this *Job) Update() (isUpdate bool, err error) {
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
	SaveEtcdJob(this.Id)
	return num > 0, err
}
