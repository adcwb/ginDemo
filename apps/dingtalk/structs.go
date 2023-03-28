package dingtalk

type testStruct struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type GetTokenStruct struct {
	ErrCode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
	ErrMsg      string `json:"errmsg"`
	ExpiresIn   int    `json:"expires_in"`
}

type GetAllFormsStruct struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	RequestId string `json:"request_id"`
	Result    []struct {
		AttendanceType int    `json:"attendance_type"`
		FlowTitle      string `json:"flow_title"`
		GmtModified    string `json:"gmt_modified"`
		IconName       string `json:"icon_name"`
		IconUrl        string `json:"icon_url"`
		IsNewProcess   bool   `json:"is_new_process"`
		ProcessCode    string `json:"process_code"`
	} `json:"result"`
	Success bool `json:"success"`
}

type GetFromProcessCodeStruct struct {
	Result struct {
		CreatorUserId string `json:"creatorUserId"`
		GmtModified   string `json:"gmtModified"`
		FormUuid      string `json:"formUuid"`
		OwnerIdType   string `json:"ownerIdType"`
		FormCode      string `json:"formCode"`
		Memo          string `json:"memo"`
		EngineType    int    `json:"engineType"`
		OwnerId       string `json:"ownerId"`
		GmtCreate     string `json:"gmtCreate"`
		SchemaContent struct {
			Icon  string `json:"icon"`
			Title string `json:"title"`
			Items []struct {
				ComponentName string `json:"componentName"`
				Props         struct {
					StaffStatusEnabled bool          `json:"staffStatusEnabled"`
					HolidayOptiIOns    []interface{} `json:"holidayOpti i ons,omitempty"`
					BizType            string        `json:"bizType,omitempty"`
					StatField          []interface{} `json:"statField,omitempty"`
					BizAlias           string        `json:"bizAlias,omitempty"`
					Id                 string        `json:"id"`
					Label              string        `json:"label"`
					Push               struct {
					} `json:"push"`
					HolidayOptions []interface{} `json:"holidayOptions,omitempty"`
					Placeholder    string        `json:"placeholder,omitempty"`
				} `json:"props"`
			} `json:"items"`
		} `json:"schemaContent"`
		AppUuid      string `json:"appUuid"`
		AppType      int    `json:"appType"`
		VisibleRange string `json:"visibleRange"`
		ListOrder    int    `json:"listOrder"`
		Name         string `json:"name"`
		Status       string `json:"status"`
	} `json:"result"`
}
