## GISA
Group information sending assistant
群组信息发送助手

#### 项目背景
```
   为了解决日常项目开发定时任务需要，原基于PHP、Laravel开发一套单节点任务管理系统来统一配置项目中定时任务，需依赖服务器的crontab服务。
   本次基于Golang开发支持分布式部署，对服务环境支持更广的定时任务管理系统，另平时使用企业微信时，常需要有特定时间群组信息公告之类，故在定时调度的基础上扩展对企业微信、钉钉群组信息支持。
```
#### 系统介绍
```
- 基于Beego(v1.12.2),前端基于AdminLte(v3.0.5)。实现RBAC认证，实现企业微信群消息定时分发，支持跨服务器分发任务。
- etcd 服务发现
- 支持通过docker compose部署
```

#### 系统截图

###### 任务调度管理
![任务调度管理](./static/img/gisa_job.jpg)

![任务配置](./static/img/gisa_job_show.jpg)

###### 权限管理
![菜单](./static/img/gisa_rbac_menu.jpg)