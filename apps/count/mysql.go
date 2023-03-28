package count

/*
权限控制
`gorm:"<-:create"`		允许读取和创建
`gorm:"<-:update"`		允许读取和更新
`gorm:"<-"`				允许读写（创建和更新）
`gorm:"<-:false"`		允许读，禁用写权限
`gorm:"->"`				只读（禁用写权限，除非它已配置）
`gorm:"->;<-:create"`	允许读取和创建
`gorm:"->:false;<-:create"`
`gorm:"-"`				使用 struct 读写时忽略此字段
`gorm:"-:all"`			读写时忽略此字段migrate with struct
 `gorm:"-:migration"`	使用 struct 迁移时忽略该字段
*/

type UserTable struct {
	Id                int    `gorm:"column:id"`
	NameClass         string `gorm:"column:name_class"`
	Status            bool   `gorm:"column:status"`
	IsDelete          bool   `gorm:"column:is_delete"`
	CreatedTime       string `gorm:"column:created_time"`
	UpdateTime        string `gorm:"column:update_time"`
	Comment           string `gorm:"column:comment"`
	UserName          string `gorm:"column:user_name"`
	UserPassword      string `gorm:"column:user_password"`
	UserPhone         string `gorm:"column:user_phone"`
	UserMoney         int    `gorm:"column:user_money"`
	UserAgentCard     string `gorm:"column:user_agent_card"`
	UserAgentDev      string `gorm:"column:user_agent_dev"`
	UserMoneyCount    int    `gorm:"column:user_money_count"`
	UserPermission    int    `gorm:"column:user_permission"`
	UserLastLoginTime int    `gorm:"column:user_last_login_time"`
	UserCount         int    `gorm:"column:user_count"`
	IsShowAdmin       int    `gorm:"column:is_show_admin"`
	Integral          int    `gorm:"column:integral"`
}

type DeviceTable struct {
	Id                int    `gorm:"column:id"`
	NameClass         string `gorm:"column:name_class"`
	Status            bool   `gorm:"column:status"`
	IsDelete          bool   `gorm:"column:is_delete"`
	CreatedTime       string `gorm:"column:created_time"`
	UpdateTime        string `gorm:"column:update_time"`
	Comment           string `gorm:"column:comment"`
	UserName          string `gorm:"column:user_name"`
	UserPassword      string `gorm:"column:user_password"`
	UserPhone         string `gorm:"column:user_phone"`
	UserMoney         int    `gorm:"column:user_money"`
	UserAgentCard     string `gorm:"column:user_agent_card"`
	UserAgentDev      string `gorm:"column:user_agent_dev"`
	UserMoneyCount    int    `gorm:"column:user_money_count"`
	UserPermission    int    `gorm:"column:user_permission"`
	UserLastLoginTime int    `gorm:"column:user_last_login_time"`
	UserCount         int    `gorm:"column:user_count"`
	IsShowAdmin       int    `gorm:"column:is_show_admin"`
	Integral          int    `gorm:"column:integral"`
}
