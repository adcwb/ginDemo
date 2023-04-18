package test

type Student struct {
	Name string
	Age  int
}

type Result struct {
	ActiveTime         string `xml:"activeTime"`
	ProdStatusName     string `xml:"prodStatusName"`
	ProdMainStatusName string `xml:"prodMainStatusName"`
	CertNumber         string `xml:"certNumber"`
	Number             string `xml:"number"`
}

type SvcCont struct {
	Result        Result `xml:"RESULT"`
	ResultCode    int    `xml:"resultCode"`
	ResultMsg     string `xml:"resultMsg"`
	TransactionID string `xml:"GROUP_TRANSACTIONID"`
}

type T struct {
	Result struct {
		ActiveTime         string `json:"ActiveTime"`
		ProdStatusName     string `json:"ProdStatusName"`
		ProdMainStatusName string `json:"ProdMainStatusName"`
		CertNumber         string `json:"CertNumber"`
		Number             string `json:"Number"`
	} `json:"Result"`
	ResultCode    int    `json:"ResultCode"`
	ResultMsg     string `json:"ResultMsg"`
	TransactionID string `json:"TransactionID"`
}

type T3 struct {
	ResultCode         string `json:"resultCode"`
	ResultMsg          string `json:"resultMsg"`
	GroupTransactionId string `json:"groupTransactionId"`
	Description        struct {
		PageIndex string `json:"pageIndex"`
		SimList   []struct {
			AccNumber      string   `json:"accNumber"`
			Iccid          string   `json:"iccid"`
			ActivationTime string   `json:"activationTime"`
			CreateTime     string   `json:"createTime"`
			SimStatus      []string `json:"simStatus"`
			CustId         string   `json:"custId"`
			Imsi           string   `json:"imsi"`
			LastChangeDate string   `json:"lastChangeDate"`
			StatusCd       string   `json:"statusCd"`
			StopTypeList   []string `json:"stopTypeList"`
		} `json:"simList"`
	} `json:"description"`
}
