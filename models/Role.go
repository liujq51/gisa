package models

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/orm"
)

// TableName 设置表名
func (a *Role) TableName() string {
	return RoleTBName()
}

// RoleQueryParam 用于搜索的类
type RoleQueryParam struct {
	BaseQueryParam
	NameLike string
}

// Role is role model structure.
type Role struct {
	BaseModel
	Id        int
	Name      string
	Slug      string
	CreatedAt string
	UpdatedAt string
}

// RoleBatchDelete 批量删除
func RoleBatchDelete(ids []int) (int64, error) {
	query := orm.NewOrm().QueryTable(RoleTBName())
	num, err := query.Filter("id__in", ids).Delete()
	return num, err
}

// RoleOne 获取单条
func RoleOne(id int) (*Role, error) {
	o := orm.NewOrm()
	m := Role{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// MapToModel get the role model from given map.
func (r Role) MapToModel(m map[string]interface{}) Role {
	r.Id = m["id"].(int)
	r.Name, _ = m["name"].(string)
	r.Slug, _ = m["slug"].(string)
	r.CreatedAt, _ = m["created_at"].(string)
	r.UpdatedAt, _ = m["updated_at"].(string)
	return r
}

//retrieve all Roles
func (this Role) FindAll() []*Role {
	var roles []*Role
	o := orm.NewOrm()
	o.QueryTable(this.TableName()).All(&roles)

	return roles
}

type roleSelectItem struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

//retrieve all Roles
func AllRoleSelectList() string {
	var (
		roles        []Role
		roleList     []roleSelectItem
		item         roleSelectItem
		roleListByte []byte
		err          error
	)

	o := orm.NewOrm()
	o.QueryTable(RoleTBName()).All(&roles)

	for _, v := range roles {
		item = roleSelectItem{}
		item.Id = v.Id
		item.Text = v.Name
		roleList = append(roleList, item)
	}
	if roleListByte, err = json.Marshal(roleList); err != nil {
		fmt.Println(err.Error())
	}

	return string(roleListByte)
}