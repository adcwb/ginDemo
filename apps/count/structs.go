package count

type getUserPersonalPermissionStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

type GetHandleInfoStruct struct {
	Code int `json:"code"`
	Data struct {
		BindCard         string   `json:"bind_card"`
		BindCardDate     string   `json:"bind_card_date"`
		BindCardId       string   `json:"bind_card_id"`
		BindOperator     string   `json:"bind_operator"`
		BindOperatorDate string   `json:"bind_operator_date"`
		ConnetMac        []string `json:"connet_mac"`
		ConnetMacNum     int      `json:"connet_mac_num"`
		CpuRate          int      `json:"cpu_rate"`
		CpuType          string   `json:"cpu_type"`
		CurrentSlot      int      `json:"currentSlot"`
		Down             string   `json:"down"`
		IgnitionCard     string   `json:"ignition_card"`
		Imei             string   `json:"imei"`
		LocBase          string   `json:"loc_base"`
		LocLat           string   `json:"loc_lat"`
		LocLng           string   `json:"loc_lng"`
		LocalTotal       int      `json:"localTotal"`
		LocalCard        string   `json:"local_card"`
		Memory           struct {
		} `json:"memory"`
		MobileSupplier    string `json:"mobile_supplier"`
		MobileType        string `json:"mobile_type"`
		PhoneSignal       string `json:"phone_signal"`
		Speed             string `json:"speed"`
		StartedDeviceTime string `json:"started_device_time"`
		StatusBattery     int    `json:"status_battery"`
		StreamRate        struct {
			Down string `json:"down"`
			Up   string `json:"up"`
		} `json:"stream_rate"`
		Time          string `json:"time"`
		Total         int    `json:"total"`
		Up            string `json:"up"`
		UserInsertSIM int    `json:"userInsertSIM"`
		VsimTotal     int    `json:"vsimTotal"`
	} `json:"data"`
}

type InitPackageStruct struct {
	Code int `json:"code"`
	Data []struct {
		Id                  int         `json:"id"`
		SetMealBaseName     string      `json:"setMealBaseName"`
		Name                string      `json:"name"`
		Price               string      `json:"price"`
		FictitiousPrice     string      `json:"fictitiousPrice"`
		NowPrice            interface{} `json:"nowPrice"`
		PurchaseType        []string    `json:"purchaseType"`
		Integral            string      `json:"integral"`
		Day                 string      `json:"day"`
		IsWholemonth        bool        `json:"is_wholeMonth"`
		IntegralChangePrice string      `json:"integral_change_price"`
		Tips                bool        `json:"tips"`
		TipsMsg             string      `json:"tipsMsg"`
	} `json:"data"`
}

type DeviceGetMobileStruct struct {
	Msg  []string `json:"msg"`
	Code int      `json:"code"`
}

type ExpressDeliveryStruct struct {
	Message   string `json:"message"`
	Nu        string `json:"nu"`
	Ischeck   string `json:"ischeck"`
	Condition string `json:"condition"`
	Com       string `json:"com"`
	Status    string `json:"status"`
	State     string `json:"state"`
	Data      []struct {
		Time    string `json:"time"`
		Ftime   string `json:"ftime"`
		Context string `json:"context"`
	} `json:"data"`
}

type ExpressDeliveryMapStruct struct {
	Message string `json:"message"`
	Nu      string `json:"nu"`
	Ischeck string `json:"ischeck"`
	Com     string `json:"com"`
	Status  string `json:"status"`
	Data    []struct {
		Time       string `json:"time"`
		Context    string `json:"context"`
		Ftime      string `json:"ftime"`
		AreaCode   string `json:"areaCode"`
		AreaName   string `json:"areaName"`
		Status     string `json:"status"`
		Location   string `json:"location"`
		AreaCenter string `json:"areaCenter"`
		AreaPinYin string `json:"areaPinYin"`
		StatusCode string `json:"statusCode"`
	} `json:"data"`
	State     string `json:"state"`
	Condition string `json:"condition"`
	RouteInfo struct {
		From struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"from"`
		Cur struct {
			Number string `json:"number"`
			Name   string `json:"name"`
		} `json:"cur"`
		To interface{} `json:"to"`
	} `json:"routeInfo"`
	IsLoop      bool   `json:"isLoop"`
	TrailUrl    string `json:"trailUrl"`
	ArrivalTime string `json:"arrivalTime"`
	TotalTime   string `json:"totalTime"`
	RemainTime  string `json:"remainTime"`
}

type GetAutonumberStruct struct {
	LengthPre int    `json:"lengthPre"`
	ComCode   string `json:"comCode"`
	Name      string `json:"name"`
	Message   string `json:"message"`
	Result    bool   `json:"result"`
}

type ExpressDeliveryPoolStruct struct {
	BalkExpressOperator string `json:"balk_express_operator"` // 快递公司
	AfterSaleId         string `json:"after_sale_id"`         // 快递单号
	FromCity            string `json:"from_city"`             // 出发城市
	ToCity              string `json:"to_city"`               // 到达城市
	PhoneNumber         string `json:"phone_number"`          // 顺丰快递手机号
}

type ExpressDeliveryCallbackStruct struct {
	Status     string `json:"status" bson:"status"`
	Billstatus string `json:"billstatus" bson:"billstatus"`
	Message    string `json:"message" bson:"message"`
	AutoCheck  string `json:"autoCheck" bson:"autoCheck"`
	ComOld     string `json:"comOld" bson:"comOld"`
	ComNew     string `json:"comNew" bson:"comNew"`
	LastResult struct {
		Message   string `json:"message" bson:"message"`
		State     string `json:"state" bson:"state"`
		Status    string `json:"status" bson:"status"`
		Condition string `json:"condition" bson:"condition"`
		Ischeck   string `json:"ischeck" bson:"ischeck"`
		Com       string `json:"com" bson:"com"`
		Nu        string `json:"nu" bson:"nu"`
		Data      []struct {
			Context    string `json:"context" bson:"context"`
			Time       string `json:"time" bson:"time"`
			Ftime      string `json:"ftime" bson:"ftime"`
			Status     string `json:"status" bson:"status"`
			AreaCode   string `json:"areaCode" bson:"areaCode"`
			AreaName   string `json:"areaName" bson:"areaName"`
			Location   string `json:"location" bson:"location"`
			AreaCenter string `json:"areaCenter" bson:"areaCenter"`
			AreaPinYin string `json:"areaPinYin" bson:"areaPinYin"`
			StatusCode string `json:"statusCode" bson:"statusCode"`
		} `json:"data" bson:"data"`
		RouteInfo struct {
			From struct {
				Number string `json:"number" bson:"number"`
				Name   string `json:"name" bson:"name"`
			} `json:"from" bson:"from"`
			Cur struct {
				Number string `json:"number" bson:"number"`
				Name   string `json:"name" bson:"name"`
			} `json:"cur"`
			To struct {
				Number string `json:"number" bson:"number"`
				Name   string `json:"name" bson:"name"`
			} `json:"to" bson:"to"`
		} `json:"routeInfo" bson:"routeInfo"`
		IsLoop      bool   `json:"isLoop" bson:"isLoop"`
		TrailUrl    string `json:"trailUrl" bson:"trailUrl"`
		ArrivalTime string `json:"arrivalTime" bson:"arrivalTime"`
		TotalTime   string `json:"totalTime" bson:"totalTime"`
		RemainTime  string `json:"remainTime" bson:"remainTime"`
	} `json:"lastResult" bson:"lastResult"`
}

type K100MongoDBDataStruct struct {
	Status     string
	Billstatus string
	Message    string
	AutoCheck  string
	ComOld     string
	ComNew     string
	LastResult struct {
		Message   string
		State     string
		Status    string
		Condition string
		Ischeck   string
		Com       string
		Nu        string
		Data      []struct {
			Context    string
			Time       string
			Ftime      string
			Status     string
			AreaCode   string
			AreaName   string
			Location   string
			AreaCenter string
			AreaPinYin string
			StatusCode string
		}
		RouteInfo struct {
			From struct {
				Number string
				Name   string
			}
			Cur struct {
				Number string
				Name   string
			}
			To struct {
				Number string
				Name   string
			}
		}
		IsLoop      bool
		TrailUrl    string
		ArrivalTime string
		TotalTime   string
		RemainTime  string
	}
}

type MongoGetK100 struct {
	CourierNumber string                `bson:"courier_number" json:"courier_number"`
	Data          K100MongoDBDataStruct `bson:"data" json:"data"`
}
