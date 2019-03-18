package main

//Settlement Log Ledger Struct

type glnbill struct {
	AdjReqNo          string `json:"ADJ_REQ_NO"`
	SndrLocalGlnCd    string `json:"LOCAL_GLN_CD"`
	AdjSDt            string `json:"ADJ_S_DT"`
	AdjEDt            string `json:"ADJ_E_DT"`
	AdjTxCnt          uint64 `json:"ADJ_TX_CNT"`
	AdjDt             string `json:"ADJ_DT"`
	TxPrcpSum         string `json:"TX_PRCP_SUM"`
	TxPrcpSumSign     string `json:"TX_PRCP_SUM_SIGN"`
	SndrAdjAmt        string `json:"SNDR_ADJ_AMT"`
	SndrAdjAmtSign    string `json:"SNDR_ADJ_AMT_SIGN"`
	RcvrAdjAmt        string `json:"RCVR_ADJ_AMT"`
	RcvrAdjAmtSign    string `json:"RCVR_ADJ_AMT_SIGN"`
	GlnAdjAmt         string `json:"GLN_ADJ_AMT"`
	GlnAdjAmtSign     string `json:"GLN_ADJ_AMT_SIGN"`
	SndrRcvgFeSum     string `json:"SNDR_RCVG_FE_SUM"`
	SndrRcvgFeSumSign string `json:"SNDR_RCVG_FE_SUM_SIGN"`
	RcvrRcvgFeSum     string `json:"RCVR_RCVG_FE_SUM"`
	RcvrRcvgFeSumSign string `json:"RCVR_RCVG_FE_SUM_SIGN"`
	GlnRcvgFeSum      string `json:"GLN_RCVG_FE_SUM"`
	GlnRcvgFeSumSign  string `json:"GLN_RCVG_FE_SUM_SIGN"`
	GlnAdjBnkSwiftCd  string `json:"GLN_ADJ_BNK_SWIFT_CD"`
	GlnAdjBnkNm       string `json:"GLN_ADJ_BNK_NM"`
	GlnAdjAcNo        string `json:"GLN_ADJ_AC_NO"`
	SndrAdjDfnYn      string `json:"LOCAL_GLN_ADJ_DFN_YN"`
	Txid              string `json:"TX_ID"`
}

// Query JSON struct
type queryArgs struct {
	AdjReqNo     string `json:"ADJ_REQ_NO"`     //정산요청번호
	ReqStartTime string `json:"REQ_START_TIME"` // 기간 시작값
	ReqEndTime   string `json:"REQ_END_TIME"`   // 기간 끝 값
	LcGlnUnqCd   string `json:"LOCALGLN_CODE"`  // Local GLN 코드
	DivCd        string `json:"DIV_CODE"`       // 구분 코드
	DeTpDvCd     string `json:"DE_TP_DV_CD"`    // 정상 취소 구분 코드
	PageSize     int32  `json:"PAGE_COUNT"`
	BookMark     string `json:"PAGE_NEXT_ID"`
}

// Event Payload Header Json
type hEvt struct {
	Target []string
	Data   interface{}
}
