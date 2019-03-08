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

type hEvt struct {
	Target []string
	Data   interface{}
}
