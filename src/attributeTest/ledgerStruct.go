package main

type exchangeRate struct {
	GlnFxNo   string
	PbldDtm   string
	PbldTn    uint16
	NatCd     string
	CurCd     string
	UpDtm     string
	TdCriXchr float64
	SndrXchr  float64
	RcvrXchr  float64
	RsvAtc    string
}

type te struct {
	Pay  string `json:"pay"`
	Data string `json:"data,omitempty"`
}

// type queryArgs struct {
// 	GlnFxNo      string
// 	NatCd        string
// 	PbldTn       uint16
// 	UpDtm        string
// 	ReqStartTime string
// 	ReqEndTime   string
// }
