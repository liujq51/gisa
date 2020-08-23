package controllers

import (
	"fmt"
	"gisa/common"
	"gisa/models"
)

type BotController struct {
	BaseController
}

//@desc 消息机器人列表
//@router /bot [get]
func (this *BotController) Index() {
	page_index, err := this.GetInt("page_index")

	pagination := common.Pagination{}
	pagination.PageCount = 20
	pagination.Url = this.URLFor("BotController.Index")

	if err == nil {
		pagination.PageIndex = page_index
	} else {
		pagination.PageIndex = 1
	}

	bots, pageTotal, err := models.Bot{}.List(pagination.PageIndex, pagination.PageCount)
	fmt.Println("bot index:", bots, pageTotal, err)
	if err != nil {
		//this.AddErrorMessage(err.Error())
		fmt.Println(err.Error())
	}

	pagination.PageTotal = pageTotal
	this.Data["modles"] = bots
	this.Data["pages"] = pagination
	this.AddBreadcrumbs("Webhook管理", this.URLFor("BotController.Index"))
	this.Data["PageTitle"] = "Webhook列表"
	this.ShowHtml("bot/index.html")
}

//@desc 添加
//@router /bot/add [get,post]
func (this *BotController) Add() {
	if this.IsPost() {
		this.Save()
	}
	this.AddBreadcrumbs("Webhook管理", this.URLFor("BotController.Index"))
	this.AddBreadcrumbs("新增", "")
	this.ShowHtml("bot/add.html")
}

//@desc 添加
func (this *BotController) Save() {
	bot := models.Bot{}
	if err := this.ParseForm(&bot); err == nil {
		if _, err := bot.Insert(); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
	this.Redirect(this.URLFor("BotController.Index"), 302)
}

//@desc 菜单删除
//@router /bot/delete [post]
func (this *BotController) Delete() {
	data := JsonData{}
	bot_id, err := this.GetInt("bot_id")
	if err != nil {
		data.Code = 400
		data.Message = "数据获取失败"
		this.ShowJSON(&data)
	}

	botModel, _ := models.FindOneBot(bot_id)
	if isDelete, err := botModel.Delete(); isDelete {
		data.Code = 200
		data.Message = "删除成功"
		this.ShowJSON(&data)
	} else {
		data.Code = 400
		data.Message = err.Error()
		this.ShowJSON(&data)
	}

	this.ShowJSON(&data)
}

//@desc 修改菜单页面
//@router /bot/edit/:bot_id [get,post,put]
func (this *BotController) Update(bot_id int) {
	var (
		err      error
		botModel *models.Bot
	)
	if bot_id == 0 {
		this.RedirectMessage(this.URLFor("BotController.Index"), "数据不存在", MESSAGE_TYPE_ERROR)
	}

	if botModel, err = models.FindOneBot(bot_id); err != nil {
		this.RedirectMessage(this.URLFor("BotController.Index"), "数据获取失败", MESSAGE_TYPE_ERROR)
	}
	this.Data["model"] = botModel
	//c.Data["routes"] = selectRoutes
	//c.Data["parents"] = selectParents
	this.AddBreadcrumbs("菜单管理", this.URLFor("BotController.Index"))
	this.AddBreadcrumbs("修改", this.URLFor("BotController.DoUpdateMenu", "bot_id", bot_id))
	this.ShowHtml("bot/edit.html")
}

//@desc 修改
//@router /bot/do_update [post,put]
func (this *BotController) DoUpdate() {
	bot := models.Bot{}
	if err := this.ParseForm(&bot); err == nil {
		if _, err := bot.Update(); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
	this.Redirect(this.URLFor("BotController.Index"), 302)
}
