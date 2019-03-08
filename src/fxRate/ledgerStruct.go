package main

// Local GLN Exchange Rate INFO
type exchangeRate struct {
	LocalGlnXchrInfUnqno string  `json:"LOCAL_GLN_XCHR_INF_UNQNO"`
	LocalGlnCd           string  `json:"LOCAL_GLN_CD"`
	UsdBidr              float64 `json:"USD_BIDR"`
	UsdOfferr            float64 `json:"USD_OFFERR"`
	BasicXchr            float64 `json:"BASIC_XCHR"`
	Dtm                  string  `json:"DTM"`
	HoprRgYn             string  `json:"HOPR_RG_YN"`
	TxID                 string  `json:"TX_ID"`
	XchrPbldDt           string  `json:"XCHR_PBLD_DT"`
	XchrPbldTn           string  `json:"XCHR_PBLD_TN"`
	XchrPbldHr           string  `json:"XCHR_PBLD_HR"`
}

// Query Key struct
type queryArgs struct {
	LocalGlnXchrInfUnqno string `json:"LOCAL_GLN_XCHR_INF_UNQNO"` // Local GLN 환율 정보
	ReqStartTime         string `json:"REQ_START_TIME"`           //기간 시작 값
	ReqEndTime           string `json:"REQ_END_TIME"`             //기간 끝 값
	LocalGlnCd           string `json:"LOCALGLN_CODE"`            //LocalGLN 고유코드
	BookMark             string `json:"PAGE_COUNT"`               // BookMark
	PageSize             int32  `json:"PAGE_NEXT_ID"`             //Paging Size
}

type queryResults struct {
	LocalGlnXchrInfUnqno string `json:"LOCAL_GLN_XCHR_INF_UNQNO"`
	Txid                 string `json:"TX_ID"`
}
