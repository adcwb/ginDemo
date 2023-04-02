package users

// testStruct 通用返回信息
type testStruct struct {
	Msg  string `json:"msg"`  // 描述信息
	Code int    `json:"code"` // 状态码
	Data string `json:"data"` // 返回数据
}

type getTokenStruct struct {
	Code   int  `json:"code"`
	IsRoot bool `json:"is_root"`
	Data   struct {
		Token string `json:"token"`
		Msg   string `json:"msg"`
	} `json:"data"`
}

type ShowUrlAllStruct struct {
	Data struct {
		List []struct {
			Path       string `json:"path"`
			Component  string `json:"component"`
			Redirect   string `json:"redirect"`
			Name       string `json:"name"`
			AlwaysShow bool   `json:"alwaysShow"`
			Children   string `json:"children"`
			Title      string `json:"title"`
			Icon       string `json:"icon"`
			Breadcrumb bool   `json:"breadcrumb"`
			ActiveMenu string `json:"activeMenu"`
			Roles      string `json:"roles"`
			IsBaseUrl  bool   `json:"is_base_url"`
			CreateTime string `json:"create_time"`
			UpdateTime string `json:"update_time"`
			IsDelete   bool   `json:"is_delete"`
			CommEnt    string `json:"comm ent"`
			Id         int    `json:"id"`
		} `json:"list"`
		Total   int  `json:"total"`
		Page    int  `json:"page"`
		HasPrev bool `json:"has_prev"`
		HasNext bool `json:"has_next"`
	} `json:"data"`
	Code int `json:"code"`
}

type GetUrlNameStruct struct {
	Msg []struct {
		Value string `json:"value"`
		Id    string `json:"id"`
	} `json:"msg"`
	Code int `json:"code"`
}

type ShowUserAllStruct struct {
	Data struct {
		List []struct {
			UserName       string `json:"user_name"`
			Id             int    `json:"id"`
			UserPhone      string `json:"user_phone"`
			UserMoney      int    `json:"user_money"`
			UserAgentCard  string `json:"user_agent_card"`
			UserAgentDev   string `json:"user_agent_dev"`
			UserMoneyCount int    `json:"user_money_count "`
			CreateTime     string `json:"create_time"`
			UpdateTime     string `json:"update_time"`
			IsDelete       bool   `json:"is_delete"`
			Comment        string `json:"comment"`
			Introduction   string `json:"Introduction"`
			Avatar         string `json:"avatar"`
			Name           string `json:"name"`
		} `json:"list"`
		Total   int  `json:"total"`
		Page    int  `json:"page"`
		HasPrev bool `json:"has_prev"`
		HasNext bool `json:"has_next"`
	} `json:"data"`
	Code int `json:"code"`
}

type GetBaseUrlStruct struct {
	Msg []struct {
		Value string `json:"value"`
		Id    int    `json:"id"`
	} `json:"msg"`
	Code int `json:"code"`
}

type GetChildrenUrlStruct struct {
	Code int `json:"code"`
	Msg  []struct {
		Label string `json:"label"`
		Value int    `json:"value"`
	} `json:"msg"`
}

// ReturnRoleTree 返回树形结构路由
type ReturnRoleTree struct {
	Msg []struct {
		Label    string        `json:"label"`
		Children []interface{} `json:"children"`
	} `json:"msg"`
	Code int `json:"code"`
}

// GetUrlAllStruct 返回所有的路由
type GetUrlAllStruct struct {
	Code int      `json:"code"`
	Msg  []string `json:"msg"`
}

// GetUrlFunctionStruct 返回所有路由下的功能
type GetUrlFunctionStruct struct {
	Code int `json:"code"`
	Msg  map[string]interface {
	} `json:"msg"`
}

// ShowPermissionAllStruct 权限组
type ShowPermissionAllStruct struct {
	Data struct {
		List []struct {
			NamePermission string `json:"NamePermission"`
			Id             int    `json:"id"`
			FunctionName   string `json:"function_name"`
			UpdateTime     string `json:"update_time"`
			CreateTime     string `json:"create_time"`
			IsDelete       bool   `json:"is_delete"`
			Status         string `json:"status"`
			Style          string `json:"style"`
			Comment        string `json:"comment"`
		} `json:"list"`
		Total   int  `json:"total"`
		Page    int  `json:"page"`
		HasPrev bool `json:"has_prev"`
		HasNext bool `json:"has_next"`
	} `json:"data"`
	Code int `json:"code"`
}

type GetPermissionGroupInfoStruct struct {
	Code      int         `json:"code"`
	Msg       interface{} `json:"msg"`
	Invisible interface{} `json:"invisible"`
}

//type GetPermissionGroupInfoStruct struct {
//	Code      string            `json:"code"`
//	Msg       map[string]*item  `json:"msg"`
//	Invisible map[string]*item2 `json:"invisible"`
//}

//type item struct {
//	FuncAll []string `json:"func_all"`
//	Title   string   `json:"title"`
//	Active  []string `json:"active"`
//}
//
//type item2 struct {
//	Func  string   `json:"func"`
//	Title []string `json:"title"`
//}

type GetPermissionGroupAllAndUserPermissionStruct struct {
	Code int `json:"code"`
	Msg  struct {
		UserTemplate  []interface{} `json:"user_template"`
		PermissionAll []struct {
			Key   int    `json:"key"`
			Label string `json:"label"`
		} `json:"permission_all"`
	} `json:"msg"`
}

type getUserPersonalPermissionStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

// SendDataStruct 发送短信请求接口
type SendDataStruct struct {
	PhoneNumbers string `json:"phone_numbers"` // 接收验证码的手机号
	OpenID       string `json:"open_id"`       // 微信客户端专用，OpenID
}

type WeChatOpenIDRequestStruct struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

type ChangeUserRequestStruct struct {
	UserName       string `json:"user_name"`
	Id             int    `json:"id"`
	UserPhone      string `json:"user_phone"`
	UserMoney      string `json:"user_money"`
	UserAgentCard  string `json:"user_agent_card"`
	UserAgentDev   string `json:"user_agent_dev"`
	UserMoneyCount string `json:"user_money_count "`
	CreateTime     string `json:"create_time"`
	UpdateTime     string `json:"update_time"`
	IsDelete       bool   `json:"is_delete"`
	Comment        string `json:"comment"`
	Introduction   string `json:"Introduction"`
	Avatar         string `json:"avatar"`
	Name           string `json:"name"`
}

type UserRechargeStruct struct {
	Mobile    string `json:"mobile"`
	UserMoney int    `json:"user_money"`
	UserId    int    `json:"user_id"`
	Mode      string `json:"mode"`
	Money     string `json:"money"`
}

type GetUserRealNameStatusStruct struct {
	Code int `json:"code"`
	Data struct {
		Local []struct {
			Iccid    string      `json:"iccid"`
			IsCert   bool        `json:"is_cert"`
			Url      string      `json:"url"`
			Operator int         `json:"operator"`
			UTel     string      `json:"u_tel"`
			DtCreate string      `json:"dt_create"`
			DtUpdate string      `json:"dt_update"`
			IdCard   string      `json:"id_card"`
			ApNumber string      `json:"ap_number"`
			DeviceId string      `json:"device_id"`
			Brand    string      `json:"brand"`
			Category int         `json:"category"`
			Channal  interface{} `json:"channal"`
		} `json:"local"`
		Vsim []struct {
			Iccid    string      `json:"iccid"`
			IsCert   bool        `json:"is_cert"`
			Url      string      `json:"url"`
			Operator int         `json:"operator"`
			UTel     string      `json:"u_tel"`
			DtCreate string      `json:"dt_create"`
			DtUpdate string      `json:"dt_update"`
			IdCard   string      `json:"id_card"`
			ApNumber string      `json:"ap_number"`
			DeviceId string      `json:"device_id"`
			Brand    string      `json:"brand"`
			Category int         `json:"category"`
			Channal  interface{} `json:"channal"`
		} `json:"vsim"`
	} `json:"data"`
}
