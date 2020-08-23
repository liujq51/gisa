package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

// Bot is menu model structure.
type Bot struct {
	BaseModel
	Id        int    `form:"bot_id"`
	Title     string `form:"title"`
	Webhook   string `form:"webhook"`
	CreatedAt string `form:"-"`
	UpdatedAt string `form:"-"`
}

// TableName 设置表名
func (p *Bot) TableName() string {
	return BotTBName()
}

//list permissions
func (this Bot) List(pageIndex, pageCount int) ([]*Bot, int, error) {
	var bots []*Bot
	var total int64
	o := orm.NewOrm()
	_, err := o.QueryTable(this.TableName()).Limit(pageCount).Offset(pageCount * (pageIndex - 1)).All(&bots)

	if err != nil {
		return bots, int(total), err
	}

	total, err = o.QueryTable(this.TableName()).Count()
	return bots, int(total), err
}

func (this *Bot) Insert() (isInsert bool, err error) {
	if this.Title == "" {
		return false, errors.New("名称不能为空")
	}

	this.CreatedAt = time.Now().Format("2001-01-01 00:00:00")
	this.UpdatedAt = this.CreatedAt
	o := orm.NewOrm()
	id, err := o.Insert(this)

	return id > 0, err
}

// 获取单条Bot
func FindOneBot(id int) (*Bot, error) {
	o := orm.NewOrm()
	this := Bot{Id: id}
	err := o.Read(&this)
	if err != nil {
		return nil, err
	}
	return &this, nil
}

//remove current menu from database
func (this *Bot) Delete() (isDelete bool, err error) {
	if this.Id <= 0 {
		return false, errors.New("删除对象不能为空")
	}

	o := orm.NewOrm()
	num, err := o.Delete(this)

	return num > 0, err
}

func (this *Bot) Update() (isUpdate bool, err error) {
	if this.Title == "" {
		return false, errors.New("名称不能为空")
	}

	this.UpdatedAt = time.Now().Format("2001-01-01 00:00:00")
	o := orm.NewOrm()
	id, err := o.Update(this)

	return id > 0, err
}

func (this Bot) GetBotSelectOption(selected int) string {
	var (
		Models []Bot
		//qb      orm.QueryBuilder
		o            orm.Ormer
		generateHtml func([]Bot) string
		TreeStr      string
		selectedStr  string
		err          error
	)
	o = orm.NewOrm()
	if _, err = o.QueryTable(this.TableName()).All(&Models); err != nil {
		fmt.Println(err.Error())
	}

	//qb, _ = orm.NewQueryBuilder("mysql")
	//o = orm.NewOrm()
	//// 构建查询对象 导出 SQL 语句
	//qb.Select("*").From(BotTBName() + " as m").OrderBy("m.id").Asc()
	//sql := qb.String()
	//o.Raw(sql).QueryRows(&Models)

	generateHtml = func(Models []Bot) (tempStr string) {
		tempStr = ""
		for _, item := range Models {
			if item.Id == selected {
				selectedStr = " selected "
			} else {
				selectedStr = ""
			}

			tempStr += `<option value="` + strconv.Itoa(item.Id) + `" ` + selectedStr + `>` + item.Title + `</option>\n`
		}

		return tempStr
	}
	TreeStr = generateHtml(Models)
	return TreeStr
}

type botMsgType struct {
	Type  string
	Title string
}

func (this Bot) GetBotMsgTypeList(selected string) (TreeStr string) {
	var (
		botMsgTypeList []botMsgType
		selectedStr    string
	)
	botMsgTypeList = []botMsgType{
		{Type: "text", Title: "文本"},
		{Type: "markdown", Title: "Markdown"},
		{Type: "news", Title: "图文"},
	}
	selectedStr = ""
	for _, item := range botMsgTypeList {
		if selected == item.Type {
			selectedStr = "selected"
		}
		TreeStr += `<option value="` + item.Type + `" ` + selectedStr + `>` + item.Title + `</option>\n`
		selectedStr = ""
	}

	return
}
