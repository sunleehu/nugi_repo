package main

// Commission fee Ledger Struct
type feeReq struct {
	AdjMnDsbReqNo  string  `json:"ADJ_MN_DSB_REQ_NO"` // 정산대금 지급요청 번호
	LcGlnUnqCd     string  `json:"LC_GLN_UNQ_CD"`     // Local GLN 고유코드
	AdjTpCd        string  `json:"ADJ_TP_CD"`         // 정산유형코드(10:매출, 20:비매출, 30:예수금)
	AdjSD          string  `json:"ADJ_S_D"`           // 정산시작일
	AdjED          string  `json:"ADJ_E_D"`           // 정산종료일
	AdjDtm         string  `json:"ADJ_DTM"`           // 정산기준일자
	WdrwAmt        float64 `json:"WDRW_AMT"`          // Sender 출금 금액
	WdrwAmtSc      float64 `json:"WDRW_AMT_SC"`       // Sender 사용 금액
	WdrwAmtCl      float64 `json:"WDRW_AMT_CL"`       // Sender 취소 금액
	WdrwAmtRf      float64 `json:"WDRW_AMT_RF"`       // Sender 환불 금액
	ColFeSum       float64 `json:"COL_FE_SUM"`        // 수수료 합계 금액
	ColFeSupp      float64 `json:"COL_FE_SUPP"`       // 수수료 공급가액
	ColFeTxamt     float64 `json:"COL_FE_TXAMT"`      // 수수료 세액
	TrxFeVal       float64 `json:"TRX_FE_VAL"`        // 거래 수수료 합계
	XrFeVal        float64 `json:"XR_FE_VAL"`         // 환율 수수료 금액
	GlnSprdVal     float64 `json:"GLN_SPRD_VAL"`      // GLN 환율  spread 금액
	XchrWdrwAmt    float64 `json:"XCHR_WDRW_AMT"`     // 환전 Sender 출금 금액
	XchrWdrwAmtSc  float64 `json:"XCHR_WDRW_AMT_SC"`  // 환전 Sender 사용 금액
	XchrWdrwAmtCl  float64 `json:"XCHR_WDRW_AMT_CL"`  // 환전 Sender 취소 금액
	XchrWdrwAmtRf  float64 `json:"XCHR_WDRW_AMT_RF"`  // 환전 Sender 환불 금액
	XchrColFeSum   float64 `json:"XCHR_COL_FE_SUM"`   // 환전 수수료 합계 금액
	XchrColFeSupp  float64 `json:"XCHR_COL_FE_SUPP"`  // 환전 수수료 공급가액
	XchrColFeTxamt float64 `json:"XCHR_COL_FE_TXAMT"` // 환전 수수료 세액
	XchrTrxFeVal   float64 `json:"XCHR_TRX_FE_VAL"`   // 환전 거래 수수료 합계
	XchrXrFeVal    float64 `json:"XCHR_XR_FE_VAL"`    // 환전 환율 수수료 금액
	XchrGlnSprdVal float64 `json:"XCHR_GLN_SPRD_VAL"` // 환전 GLN 환율  spread 금액
	GlnAdjAcNo     string  `json:"GLN_ADJ_AC_NO"`     // GLN 정산계좌번호
	DsbCompYn      string  `json:"DSB_COMP_YN"`       // 지급완료 여부(10:완료, 20:대기)
	PbldDtm        string  `json:"PBLD_DTM"`          // 기준환율 고시일자
	PbldTn         uint    `json:"PBLD_TN"`           // 기준환율 고시회차
}

// Query JSON struct
type queryArgs struct {
	AdjMnDsbReqNo string `json:"ADJ_MN_DSB_REQ_NO"` // 정산대금 지급요청 번호
	LcGlnUnqCd    string `json:"REQ_START_TIME"`    // Local GLN 고유코드
	ReqStartTime  string `json:"REQ_END_TIME"`      // 기간 시작값
	ReqEndTime    string `json:"LC_GLN_UNQ_CD"`     // 기간 끝 값
}

// Update JSON struct
type resultArgs struct {
	AdjMnDsbReqNo string `json:"ADJ_MN_DSB_REQ_NO"` // 정산 대금 지급요청 번호
	LcGlnUnqCd    string `json:"LC_GLN_UNQ_CD"`     // Local GLN 고유코드
	DsbCompYn     string `json:"DSB_COMP_YN"`       // 지급완료 여부
}

// Event Payload Header Json
type hEvt struct {
	Target []string
	Data   interface{}
}

// Insert Event Payload Data Json
type iEvt struct {
	AdjDtm string `json:"ADJ_DTM"`
}

// Update Event Payload Data Json
type uEvt struct {
	LcGlnUnqCd    string   `json:"LC_GLN_UNQ_CD"`
	AdjMnDsbReqNo []string `json:"ADJ_MN_DSB_REQ_NO"`
}
