package models

type MenuRoleRel struct {
	Id        int
	MenuId    int    `form:"-"`
	RoleId    int    `form:"-"`
	CreatedAt string `form:"-"`
	UpdatedAt string `form:"-"`
}

// TableName 设置User表名
func (m *MenuRoleRel) TableName() string {
	return RoleMenuTBName()
}
