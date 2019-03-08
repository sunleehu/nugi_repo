package main

// Payment Ledger Struct
type pmtReq struct {
	AdjMnDsbReqNo        string  `json:"ADJ_MN_DSB_REQ_NO"`         // 정산대금 지급요청 번호
	SndrLcGlnUnqCd       string  `json:"SNDR_LC_GLN_UNQ_CD"`        // Sender Local GLN 고유코드
	RcvrLcGlnUnqCd       string  `json:"RCVR_LC_GLN_UNQ_CD"`        // Receiver Local GLN 고유코드
	AdjSD                string  `json:"ADJ_S_D"`                   // 정산시작일
	AdjED                string  `json:"ADJ_E_D"`                   // 정산종료일
	AdjDtm               string  `json:"ADJ_DTM"`                   // 정산기준일자
	AmtDue               float64 `json:"AMT_DUE"`                   // 정산금액(USD)
	SndrNatCd            string  `json:"SNDR_NAT_CD"`               // Sender 국가 코드
	SndrCurCd            string  `json:"SNDR_CUR_CD"`               // Sender 통화 코드
	SndrWdrwAmt          float64 `json:"SNDR_WDRW_AMT"`             // Sender 출금 금액
	SndrWdrwAmtSc        float64 `json:"SNDR_WDRW_AMT_SC"`          // Sender 사용 금액
	SndrWdrwAmtCl        float64 `json:"SNDR_WDRW_AMT_CL"`          // Sender 취소 금액
	SndrWdrwAmtRf        float64 `json:"SNDR_WDRW_AMT_RF"`          // Sender 환불 금액
	SndrXrFeValS         float64 `json:"SNDR_XR_FE_VAL_S"`          // Sender Spread
	RcvrXrFeValS         float64 `json:"RCVR_XR_FE_VAL_S"`          // Receiver Spread
	SndrGlnSprdValS      float64 `json:"SNDR_GLN_SPRD_VAL_S"`       // Sender Gln Spread
	RcvrGlnSprdValS      float64 `json:"RCVR_GLN_SPRD_VAL_S"`       // Receiver Gln Spread
	TrxFeeIntAmtS        float64 `json:"TRX_FEE_INT_AMT_S"`         // GLN international 거래수수료
	TrxFeeSenderAmtS     float64 `json:"TRX_FEE_SENDER_AMT_S"`      // Sender 거래수수료
	TrxFeeReceiverAmtS   float64 `json:"TRX_FEE_RECEIVER_AMT_S"`    // Receiver 거래수수료
	ExSndrWdrwAmt        float64 `json:"EX_SNDR_WDRW_AMT"`          // 환전 Sender 출금 금액
	ExSndrWdrwAmtSc      float64 `json:"EX_SNDR_WDRW_AMT_SC"`       // 환전 Sender 사용 금액
	ExSndrWdrwAmtCl      float64 `json:"EX_SNDR_WDRW_AMT_CL"`       // 환전 Sender 취소 금액
	ExSndrWdrwAmtRf      float64 `json:"EX_SNDR_WDRW_AMT_RF"`       // 환전 Sender 환불 금액
	ExSndrXrFeValS       float64 `json:"EX_SNDR_XR_FE_VAL_S"`       // 환전 Sender Spread
	ExRcvrXrFeValS       float64 `json:"EX_RCVR_XR_FE_VAL_S"`       // 환전 Receiver Spread
	ExSndrGlnSprdValS    float64 `json:"EX_SNDR_GLN_SPRD_VAL_S"`    // 환전 Sender Gln Spread
	ExRcvrGlnSprdValS    float64 `json:"EX_RCVR_GLN_SPRD_VAL_S"`    // 환전 Receiver Gln Spread
	ExTrxFeeIntAmtS      float64 `json:"EX_TRX_FEE_INT_AMT_S"`      // 환전 GLN international 거래수수료
	ExTrxFeeSenderAmtS   float64 `json:"EX_TRX_FEE_SENDER_AMT_S"`   // 환전 Sender 거래수수료
	ExTrxFeeReceiverAmtS float64 `json:"EX_TRX_FEE_RECEIVER_AMT_S"` // 환전 Receiver 거래수수료
	RcvrNatCd            string  `json:"RCVR_NAT_CD"`               // Receiver 국가 코드
	RcvrCurCd            string  `json:"RCVR_CUR_CD"`               // Receiver 통화 코드
	RcvrDepoAmt          float64 `json:"RCVR_DEPO_AMT"`             // Receiver 통화 환산 Sender 출금 금액
	RcvrDepoAmtSc        float64 `json:"RCVR_DEPO_AMT_SC"`          // Receiver 통화 환산 Sender 사용 금액
	RcvrDepoAmtCl        float64 `json:"RCVR_DEPO_AMT_CL"`          // Receiver 통화 환산 Sender 취소 금액
	RcvrDepoAmtRf        float64 `json:"RCVR_DEPO_AMT_RF"`          // Receiver 통화 환산 Sender 환불 금액
	SndrXrFeValR         float64 `json:"SNDR_XR_FE_VAL_R"`          // Receiver 통화 환산 Sender Spread
	RcvrXrFeValR         float64 `json:"RCVR_XR_FE_VAL_R"`          // Receiver 통화 환산 Receiver Spread
	SndrGlnSprdValR      float64 `json:"SNDR_GLN_SPRD_VAL_R"`       // Receiver 통화 환산 Sender Gln Spread
	RcvrGlnSprdValR      float64 `json:"RCVR_GLN_SPRD_VAL_R"`       // Receiver 통화 환산 Receiver Gln Spread
	TrxFeeIntAmtR        float64 `json:"TRX_FEE_INT_AMT_R"`         // Receiver 통화 환산 GLN international 거래수수료
	TrxFeeSenderAmtR     float64 `json:"TRX_FEE_SENDER_AMT_R"`      // Receiver 통화 환산 Sender 거래수수료
	TrxFeeReceiverAmtR   float64 `json:"TRX_FEE_RECEIVER_AMT_R"`    // Receiver 통화 환산 Receiver 거래수수료
	ExRcvrDepoAmt        float64 `json:"EX_RCVR_DEPO_AMT"`          // 환전 Receiver 통화 환산 Sender 출금 금액
	ExRcvrDepoAmtSc      float64 `json:"EX_RCVR_DEPO_AMT_SC"`       // 환전 Receiver 통화 환산 Sender 사용 금액
	ExRcvrDepoAmtCl      float64 `json:"EX_RCVR_DEPO_AMT_CL"`       // 환전 Receiver 통화 환산 Sender 취소 금액
	ExRcvrDepoAmtRf      float64 `json:"EX_RCVR_DEPO_AMT_RF"`       // 환전 Receiver 통화 환산 Sender 환불 금액
	ExSndrXrFeValR       float64 `json:"EX_SNDR_XR_FE_VAL_R"`       // 환전 Receiver 통화 환산 Sender Spread
	ExRcvrXrFeValR       float64 `json:"EX_RCVR_XR_FE_VAL_R"`       // 환전 Receiver 통화 환산 Receiver Spread
	ExSndrGlnSprdValR    float64 `json:"EX_SNDR_GLN_SPRD_VAL_R"`    // 환전 Receiver 통화 환산 Sender Gln Spread
	ExRcvrGlnSprdValR    float64 `json:"EX_RCVR_GLN_SPRD_VAL_R"`    // 환전 Receiver 통화 환산 Receiver Gln Spread
	ExTrxFeeIntAmtR      float64 `json:"EX_TRX_FEE_INT_AMT_R"`      // 환전 Receiver 통화 환산 GLN international 거래수수료
	ExTrxFeeSenderAmtR   float64 `json:"EX_TRX_FEE_SENDER_AMT_R"`   // 환전 Receiver 통화 환산 Sender 거래수수료
	ExTrxFeeReceiverAmtR float64 `json:"EX_TRX_FEE_RECEIVER_AMT_R"` // 환전 Receiver 통화 환산 Receiver 거래수수료
	BankNm               string  `json:"BANK_NM"`                   // 은행명
	AdjAcNo              string  `json:"ADJ_AC_NO"`                 // 정산 계좌번호
	SwiftCd              string  `json:"SWIFT_CD"`                  // SWIFT CODE
	AdjCompYn            string  `json:"ADJ_COMP_YN"`               // 지급완료 여부(10:완료, 20:대기)
	PbldDtm              string  `json:"PBLD_DTM"`                  // 기준환율 고시일자
	PbldTn               uint    `json:"PBLD_TN"`                   // 기준환율 고시회차
}

// Query Json key
type queryArgs struct {
	AdjMnDsbReqNo string `json:"ADJ_MN_DSB_REQ_NO"`
	LcGlnUnqCd    string `json:"REQ_START_TIME"`
	ReqStartTime  string `json:"REQ_END_TIME"`
	ReqEndTime    string `json:"LC_GLN_UNQ_CD"`
}

// Update Json Key
type resultArgs struct {
	AdjMnDsbReqNo  string `json:"ADJ_MN_DSB_REQ_NO"`
	SndrLcGlnUnqCd string `json:"SNDR_LC_GLN_UNQ_CD"`
	RcvrLcGlnUnqCd string `json:"RCVR_LC_GLN_UNQ_CD"`
	AdjCompYn      string `json:"ADJ_COMP_YN"`
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
