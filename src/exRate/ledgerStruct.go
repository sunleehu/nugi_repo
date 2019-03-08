package main

// Exchange Rate Ledger Struct
type exchangeRate struct {
	GlnFxNo   string  `json:"GLN_FX_NO"`         // 환율 이력 고유 번호
	PbldDtm   string  `json:"PBLD_DTM"`          // 고시일자
	PbldTn    uint    `json:"PBLD_TN"`           // 고시회차
	NatCd     string  `json:"NAT_CD"`            // 국가코드
	CurCd     string  `json:"CUR_CD"`            // 통화코드
	UpDtm     string  `json:"UP_DTM"`            // 갱신시간
	TdCriXchr float64 `json:"TD_CRI_XCHR"`       // 대미환율(매매기준환율)
	SndrXchr  float64 `json:"SNDR_XCHR"`         // Sender 기준 환율
	RcvrXchr  float64 `json:"RCVR_XCHR"`         // Receiver 기준 환율
	RsvAtc    string  `json:"RSV_ATC,omitempty"` // reserved 항목
}

// Query Key struct
type queryArgs struct {
	GlnFxNo      string `json:"GLN_FX_NO"`      // 환율 이력 고유 번호
	NatCd        string `json:"NAT_CD"`         // 국가코드
	PbldDtm      string `json:"PBLD_DTM"`       // 고시일자
	PbldTn       uint16 `json:"PBLD_TN"`        // 고시회차
	UpDtm        string `json:"UP_DTM"`         // 갱신시간
	ReqStartTime string `json:"REQ_START_TIME"` // 기간 시작 값
	ReqEndTime   string `json:"REQ_END_TIME"`   // 기간 끝 값
}
