package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"gisa/common/config"

	"github.com/astaxie/beego/orm"
)

// TableName 设置User表名
func (a *User) TableName() string {
	return UserTBName()
}

// UserQueryParam 用于查询的类
type UserQueryParam struct {
	BaseQueryParam
	UserNameLike string //模糊查询
	RealNameLike string //模糊查询
	Mobile       string //精确查询
	SearchStatus string //为空不查询，有值精确查询
}

// User is user model structure.
type User struct {
	Id                 int
	RealName           string `orm:"size(32)"`
	UserName           string `orm:"size(24)"`
	Password           string `json:"-"`
	IsSuper            bool
	Status             int
	Mobile             string       `orm:"size(16)"`
	Email              string       `orm:"size(256)"`
	Avatar             string       `orm:"size(256)"`
	RoleIds            []int        `orm:"-" form:"RoleIds"`
	ResourceUrlForList []string     `orm:"-"`
	Roles              []Role       `orm:"-" json:"role"`
	Permissions        []Permission `orm:"-" json:"permissions"`
	MenuIds            []int        `orm:"-" json:"menu_ids"`
	Level              string       `orm:"-" json:"level"`
	LevelName          string       `orm:"-" json:"level_name"`
}

// UserPageList 获取分页数据
func UserPageList(params *UserQueryParam) ([]*User, int64) {
	query := orm.NewOrm().QueryTable(UserTBName())
	data := make([]*User, 0)
	//默认排序
	sortorder := "Id"
	switch params.Sort {
	case "Id":
		sortorder = "Id"
	}
	if params.Order == "desc" {
		sortorder = "-" + sortorder
	}
	query = query.Filter("username__istartswith", params.UserNameLike)
	query = query.Filter("realname__istartswith", params.RealNameLike)
	if len(params.Mobile) > 0 {
		query = query.Filter("mobile", params.Mobile)
	}
	if len(params.SearchStatus) > 0 {
		query = query.Filter("status", params.SearchStatus)
	}
	total, _ := query.Count()
	query.OrderBy(sortorder).Limit(params.Limit, params.Offset).All(&data)
	return data, total
}

// UserOne 根据id获取单条
func UserOne(id int) (*User, error) {
	o := orm.NewOrm()
	m := User{Id: id}
	err := o.Read(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UserOneByUserName 根据用户名密码获取单条
func UserOneByUserName(username, password string) (*User, error) {
	m := User{}
	err := orm.NewOrm().QueryTable(UserTBName()).Filter("user_name", username).Filter("password", password).One(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func inMethodArr(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if strings.EqualFold(arr[i], str) {
			return true
		}
	}
	return false
}

func (u User) CheckPermission(path string, method string) bool {
	var (
		reg       *regexp.Regexp
		matchPath string
		err       error
	)
	path, _ = url.PathUnescape(path)
	//fmt.Println("checkPermission", path, method)
	//return true
	//if t.IsSuperAdmin() {
	//	return true
	//}

	logoutCheck, _ := regexp.Compile(config.Url("/logout") + "(.*?)")

	if logoutCheck.MatchString(path) {
		return true
	}

	if path == "" {
		return false
	}

	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	//path = strings.Replace(path, common.EditPKKey, "id", -1)
	//path = strings.Replace(path, constant.DetailPKKey, "id", -1)
	//
	//path, params := getParam(path)
	//for key, value := range formParams {
	//	if len(value) > 0 {
	//		params.Add(key, value[0])
	//	}
	//}
	//
	for _, v := range u.Permissions {
		fmt.Println("http:", v.HttpPath, v.HttpMethod)
		httpPath := strings.Split(v.HttpPath, "\n")
		httpMethod := strings.Split(v.HttpMethod, ",")
		if httpMethod[0] == "" || inMethodArr(httpMethod, method) {

			if httpPath[0] == "*" {
				return true
			}

			for i := 0; i < len(v.HttpPath); i++ {
				matchPath = config.Url(strings.TrimSpace(httpPath[i]))
				reg, err = regexp.Compile(matchPath)
				if err != nil {
					//logger.Error("CheckPermissions error: ", err)
					continue
				}
				if reg.FindString(path) == path {
					return true
				}
			}
		}
	}

	return true
}

func (u User) GetAllRoleId() []interface{} {

	var ids = make([]interface{}, len(u.Roles))

	for key, role := range u.Roles {
		ids[key] = role.Id
	}
	ids = append(ids, 1)
	ids = append(ids, 2)
	return ids
}

// WithRoles query the role info of the user.
func (u *User) WithRoles() *User {
	//roleModel, _ := u.Table("goadmin_role_users").
	//	LeftJoin("goadmin_roles", "goadmin_roles.id", "=", "goadmin_role_users.role_id").
	//	Where("user_id", "=", t.Id).
	//	Select("goadmin_roles.id", "goadmin_roles.name", "goadmin_roles.slug",
	//		"goadmin_roles.created_at", "goadmin_roles.updated_at").
	//	All()
	//
	//for _, role := range roleModel {
	//	u.Roles = append(u.Roles, new(Role).MapToModel(role))
	//}
	var roles []*Role
	qb, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象
	//qb.Select("r.id", "r.name", "r.slug", "r.creatd_at", "r.updated_at").
	//	From(RoleUserTBName() + " as ru").
	//	LeftJoin(RoleTBName() + " as r").On("ru.role_id = r.id").
	//	Where("ru.user_id = ?").
	//	OrderBy("r.id").Desc()
	//.Limit(10).Offset(0)
	// 导出 SQL 语句
	sql := qb.Select("r.id", "r.name", "r.slug", "r.created_at", "r.updated_at").
		From(RoleUserTBName() + " as ru").
		LeftJoin(RoleTBName() + " as r").On("ru.role_id = r.id").
		Where("ru.user_id = ?").
		OrderBy("r.id").Desc().String()
	// 执行 SQL 语句
	o := orm.NewOrm()
	o.Raw(sql, u.Id).QueryRows(&roles)
	if len(u.Roles) > 0 {
		u.Level = u.Roles[0].Slug
		u.LevelName = u.Roles[0].Name
	}

	return u
}

// WithPermissions query the permission info of the user.
func (u *User) WithPermissions() *User {
	var (
		permissions     = make([]map[string]interface{}, 0)
		userPermissions = make([]map[string]interface{}, 0)
	)
	//roleIds := u.GetAllRoleId()
	roleIds := "1,2"
	qb, _ := orm.NewQueryBuilder("mysql")
	qb1, _ := orm.NewQueryBuilder("mysql")
	// 构建查询对象 导出 SQL 语句
	sql := qb.Select("p.id", "p.name", "p.slug", "p.http_path", "p.http_method", "p.created_at", "p.updated_at").
		From(RolePermissionTBName() + " as rp").
		LeftJoin(PermissionTBName() + " as p").On("rp.permission_id = p.id").
		Where("rp.role_id in (" + roleIds + ")").
		OrderBy("p.id").Desc().String()

	permissionSql := qb1.Select("p.id", "p.name", "p.slug", "p.http_path", "p.http_method", "p.created_at", "p.updated_at").
		From(RolePermissionTBName() + " as rp").
		LeftJoin(PermissionTBName() + " as p").On("rp.permission_id = p.id").
		Where("rp.role_id in (" + roleIds + ")").
		OrderBy("p.id").Desc().String()
	// 执行 SQL 语句
	o := orm.NewOrm()
	o.Raw(sql).QueryRows(&permissions)
	// 执行 SQL 语句
	permissionOrm := orm.NewOrm()
	permissionOrm.Raw(permissionSql).QueryRows(&userPermissions)

	permissions = append(permissions, userPermissions...)
	for i := 0; i < len(permissions); i++ {
		exist := false
		for j := 0; j < len(u.Permissions); j++ {
			if u.Permissions[j].Id == permissions[i]["id"] {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		u.Permissions = append(u.Permissions, new(Permission).MapToModel(permissions[i]))
	}

	return u
}

// WithMenus query the menu info of the user.
func (u *User) WithMenus() *User {
	type MenuModel struct {
		Title    string
		MenuId   int
		ParentId int
	}
	var (
		menuIdsModel []MenuModel
		roleIds      []interface{}
		qb           orm.QueryBuilder
		menuIds      []int
		o            orm.Ormer
	)
	qb, _ = orm.NewQueryBuilder("mysql")
	// 执行 SQL 语句
	o = orm.NewOrm()
	if u.IsSuper {
		// 构建查询对象 导出 SQL 语句
		qb.Select("m.title", "m.id as menu_id", "m.parent_id").
			From(MenuTBName() + " as m").
			OrderBy("m.id").Desc()
		sql := qb.String()
		o.Raw(sql).QueryRows(&menuIdsModel)
	} else {
		roleIds = u.GetAllRoleId()
		// 构建查询对象 导出 SQL 语句
		sql := qb.Select("m.title", "rm.menu_id", "m.parent_id").
			From(RoleMenuTBName() + " as rm").
			LeftJoin(MenuTBName() + " as m").On("rm.menu_id = m.id").
			Where("rm.role_id in (?)").
			OrderBy("m.id").Desc().String()
		o.Raw(sql, roleIds).QueryRows(&menuIdsModel)
	}

	for _, mid := range menuIdsModel {
		if mid.ParentId != 0 {
			for _, mid2 := range menuIdsModel {
				if mid2.MenuId == mid.ParentId {
					menuIds = append(menuIds, mid.MenuId)
					break
				}
			}
		} else {
			menuIds = append(menuIds, mid.MenuId)
		}
	}
	u.MenuIds = menuIds
	return u
}
