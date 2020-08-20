package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
)

// Item is an menu item.
type MenuItem struct {
	Title        string     `json:"title"`
	ID           int        `json:"id"`
	Url          string     `json:"url"`
	IsLinkUrl    bool       `json:"isLinkUrl"`
	Icon         string     `json:"icon"`
	Header       string     `json:"header"`
	Active       string     `json:"active"`
	ChildrenList []MenuItem `json:"childrenList"`
}

// Menu contains list of menu items and other info.
type MenuInfo struct {
	List     []MenuItem          `json:"list"`
	Options  []map[string]string `json:"options"`
	MaxOrder int                 `json:"maxOrder"`
}

// MenuModel is menu model structure.
type Menu struct {
	BaseModel
	Id        int    `form:"menu_id"`
	Title     string `form:"title"`
	ParentId  int    `form:"parent_id"`
	Type      int    `form:"-"`
	Order     int    `form:"-"`
	Icon      string `form:"icon"`
	Uri       string `form:"uri"`
	Header    string `form:"-"`
	CreatedAt string `form:"-"`
	UpdatedAt string `form:"-"`
}

// TableName 设置User表名
func (m *Menu) TableName() string {
	return MenuTBName()
}

// GetGlobalMenu return Menu of given user model.
func GetGlobalMenu(user interface{}) *MenuInfo {
	var (
		menuIdsModel []Menu
		menuOption   = make([]map[string]string, 0)
		qb           orm.QueryBuilder
		roleIds      []interface{}
		o            orm.Ormer
	)
	curUser := user.(User)

	qb, _ = orm.NewQueryBuilder("mysql")
	// 执行 SQL 语句
	o = orm.NewOrm()
	if curUser.IsSuper {
		// 构建查询对象 导出 SQL 语句
		qb.Select("*").From(MenuTBName() + " as m").
			OrderBy("m.order").Asc()
		sql := qb.String()
		o.Raw(sql).QueryRows(&menuIdsModel)
	} else {
		roleIds = curUser.GetAllRoleId()
		// 构建查询对象 导出 SQL 语句
		sql := qb.Select("*").From(RoleMenuTBName() + " as rm").
			LeftJoin(MenuTBName() + " as m").On("rm.menu_id = m.id").
			Where("rm.role_id in (?)").
			OrderBy("m.order").Asc().String()
		o.Raw(sql, roleIds).QueryRows(&menuIdsModel)
	}
	for i := 0; i < len(menuIdsModel); i++ {
		menuOption = append(menuOption, map[string]string{
			"id":    strconv.Itoa(menuIdsModel[i].Id),
			"title": menuIdsModel[i].Title,
		})
	}

	menuList := constructMenuTree(menuIdsModel, 0)
	maxOrder := int(0)
	if len(menuIdsModel) > 0 {
		maxOrder = menuIdsModel[len(menuIdsModel)-1].ParentId
	}

	return &MenuInfo{
		List:     menuList,
		Options:  menuOption,
		MaxOrder: maxOrder,
	}
}

func constructMenuTree(menus []Menu, parentID int) []MenuItem {

	branch := make([]MenuItem, 0)

	var title string
	for j := 0; j < len(menus); j++ {
		if parentID == menus[j].ParentId {
			if menus[j].Type == 1 {
				//title = language.Get(menus[j]["title"].(string))
				title = menus[j].Title
			} else {
				title = menus[j].Title
			}

			child := MenuItem{
				Title:        title,
				ID:           int(menus[j].Id),
				Url:          menus[j].Uri,
				Icon:         menus[j].Icon,
				Header:       menus[j].Header,
				Active:       "",
				ChildrenList: constructMenuTree(menus, menus[j].Id),
			}

			branch = append(branch, child)
		}
	}

	return branch
}

func GetMenuTreeHtml(items interface{}, httpPath string) string {
	var (
		menuItems   []MenuItem
		menuTreeStr string
		dfs         func([]MenuItem, string) (string, bool)
		subStr      string
		liClass     string
		aClass      string
		piClass     string
	)
	dfs = func(menuItems []MenuItem, httpPath string) (tempStr string, parentFlag bool) {
		tempStr = ""
		parentFlag = false
		for _, item := range menuItems {
			flag := false
			subStr = ""
			subFlag := false

			if item.ChildrenList != nil && len(item.ChildrenList) > 0 {
				subStr, subFlag = dfs(item.ChildrenList, httpPath)
			}

			if item.Url == httpPath || subFlag {
				parentFlag = true
				flag = true
			}
			if flag {
				liClass = " menu-open " + httpPath
				aClass = " active "
			} else {
				liClass = " " + httpPath
				aClass = " "
			}
			piClass = ""
			if subStr != "" {
				piClass = `<i class="right fas fa-angle-left"></i>`
			}
			tempStr += `<li class="nav-item has-treeview ` + liClass + ` ">`
			tempStr += `<a href="` + item.Url + `" class="nav-link ` + aClass + `">
                        <i class="nav-icon ` + item.Icon + `"></i>
                        <p>` + item.Title + piClass + `</p>
                    </a>`
			if subStr != "" {
				tempStr += `<ul class="nav nav-treeview ">`
				tempStr += subStr
				tempStr += `</ul>`
			}
			tempStr += `</li>`
		}

		return tempStr, parentFlag
	}
	menuItems = items.([]MenuItem)
	menuTreeStr, _ = dfs(menuItems, httpPath)
	return menuTreeStr
}

func GetMenuNestableHtml() string {
	var (
		menuIdsModel []Menu
		qb           orm.QueryBuilder
		o            orm.Ormer
		menuTreeStr  string
		dfs          func([]MenuItem) string
		menuItems    []MenuItem
		subStr       string
	)
	qb, _ = orm.NewQueryBuilder("mysql")
	o = orm.NewOrm()
	// 构建查询对象 导出 SQL 语句
	qb.Select("*").From(MenuTBName() + " as m").OrderBy("m.order").Asc()
	sql := qb.String()
	o.Raw(sql).QueryRows(&menuIdsModel)
	menuItems = constructMenuTree(menuIdsModel, 0)

	dfs = func(menuItems []MenuItem) (tempStr string) {
		tempStr = `<ol class="dd-list">`
		for _, item := range menuItems {
			subStr = ""
			if item.ChildrenList != nil && len(item.ChildrenList) > 0 {
				subStr = dfs(item.ChildrenList)
			}
			tempStr += `<li class="dd-item" data-id="` + strconv.Itoa(item.ID) + `">`
			tempStr += `<div class="dd-handle">`
			tempStr += `<i class='fa fa-bar-chart'></i>&nbsp;<strong>` + item.Title + `</strong>&nbsp;&nbsp;&nbsp;<a href="` + item.Url + `" class="dd-nodrag">` + item.Url + `</a>`
			tempStr += `	   <span class="float-right dd-nodrag">`
			tempStr += `		   <a href="/info/menu/update/` + strconv.Itoa(item.ID) + `" class="btn btn-xs btn-success"><i class="fa fa-edit"></i></a>`
			tempStr += `		   <a href="javascript:void(0);" data-id="` + strconv.Itoa(item.ID) + `" class="tree_branch_delete btn btn-xs btn-danger"><i class="fa fa-trash"></i></a>`
			tempStr += `	   </span>`
			tempStr += `   </div>`
			if subStr != "" {
				tempStr += subStr
			}

			tempStr += `</li>`
		}
		tempStr += `</ol>`

		return tempStr
	}
	menuTreeStr = dfs(menuItems)
	return menuTreeStr
}
func GetMenuSelectOption(selected int) string {
	var (
		menuIdsModel   []Menu
		qb             orm.QueryBuilder
		o              orm.Ormer
		menuItems      []MenuItem
		dfs            func([]MenuItem, string) string
		subStr         string
		menuTreeStr    string
		optionSelected string
	)
	qb, _ = orm.NewQueryBuilder("mysql")
	o = orm.NewOrm()
	// 构建查询对象 导出 SQL 语句
	qb.Select("*").From(MenuTBName() + " as m").OrderBy("m.order").Asc()
	sql := qb.String()
	o.Raw(sql).QueryRows(&menuIdsModel)

	menuItems = constructMenuTree(menuIdsModel, 0)
	dfs = func(menuItems []MenuItem, identStr string) (tempStr string) {
		tempStr = ""
		optionSelected = ""
		identStr += "&nbsp;&nbsp;&nbsp;&nbsp;"
		for _, item := range menuItems {
			subStr = ""
			if item.ChildrenList != nil && len(item.ChildrenList) > 0 {
				subStr = dfs(item.ChildrenList, identStr)
			}
			if item.ID == selected {
				optionSelected = " selected "
			}
			tempStr += `<option value="` + strconv.Itoa(item.ID) + `" ` + optionSelected + `>` + identStr + `┝  ` + item.Title + `</option>\n`
			if subStr != "" {
				tempStr += subStr
			}
			optionSelected = ""
		}

		return tempStr
	}
	menuTreeStr = dfs(menuItems, "")
	return menuTreeStr
}

func UpdateMenuParentIdAndOrder(menuId int, parentId int, menuOrder int) bool {
	o := orm.NewOrm()
	menu := Menu{Id: menuId}
	if o.Read(&menu) == nil {
		menu.ParentId = parentId
		menu.Order = menuOrder
		if _, err := o.Update(&menu); err == nil {
			return true
		}
	}
	return false
}

func (this *Menu) Insert() (isInsert bool, err error) {
	if this.Title == "" {
		return false, errors.New("菜单名称不能为空")
	}

	this.CreatedAt = time.Now().Format("2001-01-01 00:00:00")
	this.UpdatedAt = this.CreatedAt
	o := orm.NewOrm()
	id, err := o.Insert(this)

	return id > 0, err
}

//remove current menu from database
func (m *Menu) Delete() (isDelete bool, err error) {
	if m.IsNewRecord() {
		return false, errors.New("删除对象不能为空")
	}

	o := orm.NewOrm()
	MenuModel := Menu{}
	o.QueryTable(m.TableName()).Filter("parent_id", m.Id).One(&MenuModel)

	if !MenuModel.IsNewRecord() {
		return false, errors.New("存在子菜单，不能删除")
	}

	num, err := o.Delete(m)

	return num > 0, err
}

//Find one Menu by id from database
func (m *Menu) FindById(id int) error {
	if id <= 0 {
		return errors.New("菜单ID不能为空")
	}

	o := orm.NewOrm()
	return o.QueryTable(m.TableName()).Filter("id", id).One(m)
}

// 获取单条MenuOne
func MenuOne(id int) (*Menu, error) {
	o := orm.NewOrm()
	m := Menu{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

//check current model is new
func (m *Menu) IsNewRecord() bool {
	return m.Id <= 0
}

// CheckRole check the role if has permission to get the menu.
func (m *Menu) CheckRole(roleId int) bool {
	var (
		o   orm.Ormer
		num int64
	)
	o = orm.NewOrm()
	num, _ = o.QueryTable(RoleMenuTBName()).Filter("menu_id", m.Id).Filter("role_id", roleId).Count()
	return num > 0
}

// AddRole add a role to the menu.
func (m *Menu) AddRole(roleId int) (int64, error) {
	var (
		o   orm.Ormer
		err error
	)
	o = orm.NewOrm()
	if roleId != 0 {
		if !m.CheckRole(roleId) {
			menuRole := MenuRoleRel{
				MenuId: m.Id,
				RoleId: roleId,
			}
			if _, err = o.Insert(&menuRole); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return 0, nil
}

// DeleteRoles delete roles with menu.
func (m *Menu) DeleteRoles() error {
	o := orm.NewOrm()
	_, err := o.QueryTable(RoleMenuTBName()).Filter("menu_id", m.Id).Delete()
	return err
}

//retrieve all Roles
func AllRoleSelectedList(menuId int) string {
	var (
		menuRoleRel  []MenuRoleRel
		roleIdList   []int
		roleListByte []byte
		err          error
	)

	o := orm.NewOrm()
	o.QueryTable(RoleMenuTBName()).Filter("menu_id", menuId).All(&menuRoleRel, "role_id")

	for _, v := range menuRoleRel {
		roleIdList = append(roleIdList, v.RoleId)
	}
	if roleListByte, err = json.Marshal(roleIdList); err != nil {
		fmt.Println(err.Error())
	}

	return string(roleListByte)
}
