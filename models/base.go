package models

import "github.com/astaxie/beego"

// Base is base model structure.
type BaseModel struct {
}

// TableName 下面是统一的表名管理
func TableName(name string) string {
	prefix := beego.AppConfig.String("db_dt_prefix")
	return prefix + name
}

// UserTBName 获取 User 对应的表名称
func UserTBName() string {
	return TableName("user")
}

// PermissionTBName 获取 Permission 对应的表名称
func PermissionTBName() string {
	return TableName("permissions")
}

// RoleTBName 获取 Role 对应的表名称
func RoleTBName() string {
	return TableName("roles")
}

// RolePermissionTBName 角色与资源多对多关系表
func RolePermissionTBName() string {
	return TableName("role_permissions")
}

// RoleUserRelTBName 角色与用户多对多关系表
func RoleUserTBName() string {
	return TableName("role_users")
}

func MenuTBName() string {
	return TableName("menu")
}

func RoleMenuTBName() string {
	return TableName("role_menu")
}
