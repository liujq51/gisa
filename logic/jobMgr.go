package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gisa/common/crontab"
	"time"

	"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// 任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	// 单例
	G_jobMgr *JobMgr
)

// 初始化管理器
func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	// 初始化配置
	duration, _ := beego.AppConfig.Int64("crontab::etcdDialTimeout")
	config = clientv3.Config{
		Endpoints:   beego.AppConfig.Strings("crontab::etcdEndpoints"), // 集群地址
		DialTimeout: time.Duration(duration) * time.Millisecond,        // 连接超时
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}
	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// 保存任务
func (jobMgr *JobMgr) SaveJob(job *crontab.Job) (oldJob *crontab.Job, err error) {
	// 把任务保存到/cron/jobs/任务名 -> json
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj crontab.Job
	)

	// etcd的保存key
	jobKey = crontab.JOB_SAVE_DIR + job.Name
	// 任务信息json
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}
	// 保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		fmt.Println(err.Error())
	}
	// 如果是更新, 那么返回旧值
	if putResp.PrevKv != nil {
		// 对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			fmt.Println(err.Error())
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (oldJob *crontab.Job, err error) {
	var (
		jobKey    string
		delResp   *clientv3.DeleteResponse
		oldJobObj crontab.Job
	)

	// etcd中保存任务的key
	jobKey = crontab.JOB_SAVE_DIR + name

	// 从etcd中删除它
	if delResp, err = jobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	// 返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 解析一下旧值, 返回它
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}
	return
}

// 列举任务
func (jobMgr *JobMgr) ListJobs() (jobList []*crontab.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *crontab.Job
	)
	// 任务保存的目录
	dirKey = crontab.JOB_SAVE_DIR
	// 获取目录下所有任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		fmt.Println(err.Error())
	}
	// 初始化数组空间
	jobList = make([]*crontab.Job, 0)
	// len(jobList) == 0
	// 遍历所有任务, 进行反序列化
	for _, kvPair = range getResp.Kvs {
		job = &crontab.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

// 查看任务
func (jobMgr *JobMgr) JobDetail(name string) (job *crontab.Job, err error) {
	var (
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		jobKey  string
	)
	// 任务保存的目录
	jobKey = crontab.JOB_SAVE_DIR + name
	// 获取任务信息
	if getResp, err = jobMgr.kv.Get(context.TODO(), jobKey); err != nil {
		fmt.Println(err.Error())
	}
	job = &crontab.Job{}
	// len(jobList) == 0
	// 遍历所有任务, 进行反序列化
	for _, kvPair = range getResp.Kvs {
		tempJob := &crontab.Job{}
		if err = json.Unmarshal(kvPair.Value, tempJob); err != nil {
			err = nil
			continue
		}
		fmt.Println("job name:", tempJob, tempJob.Name, name)
		if tempJob.Name == name {
			job = tempJob
		}
	}
	fmt.Println("job:", job)
	return
}

// 杀死任务
func (jobMgr *JobMgr) KillJob(name string) (err error) {
	// 更新一下key=/cron/killer/任务名
	var (
		killerKey      string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)

	// 通知worker杀死对应任务
	killerKey = crontab.JOB_KILLER_DIR + name

	// 让worker监听到一次put操作, 创建一个租约让其稍后自动过期即可
	if leaseGrantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	// 租约ID
	leaseId = leaseGrantResp.ID

	// 设置killer标记
	if _, err = jobMgr.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}
	return
}
