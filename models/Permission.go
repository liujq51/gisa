package models

import (
	"github.com/astaxie/beego/orm"
)

// Permission 权限控制资源表
type Permission struct {
	Id         int
	Name       string
	Slug       string
	HttpMethod string
	HttpPath   string
	CreatedAt  string
	UpdatedAt  string
}

// TableName 设置表名
func (p *Permission) TableName() string {
	return PermissionTBName()
}

// MapToModel get the permission model from given map.
func (p Permission) MapToModel(m map[string]interface{}) Permission {
	p.Id = m["id"].(int)
	p.Name, _ = m["name"].(string)
	p.Slug, _ = m["slug"].(string)

	p.HttpMethod, _ = m["http_method"].(string)
	//if methods != "" {
	//	p.HttpMethod = strings.Split(methods, ",")
	//} else {
	//	p.HttpMethod = []string{""}
	//}

	p.HttpPath, _ = m["http_path"].(string)
	//p.HttpPath = strings.Split(path, "\n")
	p.CreatedAt, _ = m["created_at"].(string)
	p.UpdatedAt, _ = m["updated_at"].(string)
	return p
}

//list permissions
func (this Permission) List(pageIndex, pageCount int) ([]*Permission, int, error) {
	var permissions []*Permission
	var total int64
	o := orm.NewOrm()
	_, err := o.QueryTable(this.TableName()).Limit(pageCount).Offset(pageCount * (pageIndex - 1)).All(&permissions)

	if err != nil {
		return permissions, int(total), err
	}

	total, err = o.QueryTable(this.TableName()).Count()
	return permissions, int(total), err
}
