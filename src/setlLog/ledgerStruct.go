package main

//Settlement Log Ledger Struct

type settlmentData struct {
	//Common Data
	GlnDeNo     string `json:"GLN_DE_NO"`      // GLN 거래 번호(Key)
	GlnDeDtm    string `json:"GLN_DE_DTM"`     // GLN 거래 일시
	MbrSvcDvCd  string `json:"MBR_SVC_DV_CD"`  // 회원 서비스 구분 코드
	DeTpDvCd    string `json:"DE_TP_DV_CD"`    // 거래 유형 구분 코드
	RecvDeNo    string `json:"RECV_DE_NO"`     // 사용처 거래 추적 번호
	RecvUnqCd   string `json:"RECV_UNQ_CD"`    // 사용처 고유코드
	OriGlnDeNo  string `json:"ORI_GLN_DE_NO"`  // 원 GLN 거래 번호
	CnclGlnDeNo string `json:"CNCL_GLN_DE_NO"` // 취소 GLN 거래 번호
	PbldDtm     string `json:"PBLD_DTM"`       // 기준 환율 고시 일자
	PbldTn      uint64 `json:"PBLD_TN"`        // 기준 환율 고시 회차

	//Receiver Data
	RcvrLcGlnUnqCd  string  `json:"RCVR_LC_GLN_UNQ_CD"`  // Receiver Local GLN 고유코드 (SubKey)
	RcvrNatCd       string  `json:"RCVR_NAT_CD"`         // Receiver 국가 코드
	RcVrCurCd       string  `json:"RCVR_CUR_CD"`         // Receiver 통화 코드
	RcvrDepoAmt     float64 `json:"RCVR_DEPO_AMT"`       // Receiver 입금 금액
	RcvrSndrGlnXchr float64 `json:"RCVR_SNDR_GLN_XCHR"`  // Receiver/Sender GLN 환율
	RcvrXrFe        float64 `json:"RCVR_XR_FE"`          // Receiver 환율 수수료
	SndrXrFeValR    float64 `json:"SNDR_XR_FE_VAL_R"`    // Sender 환율 수수료 금액
	RcvrXrFeValR    float64 `json:"RCVR_XR_FE_VAL_R"`    // Receiver 환율 수수료 금액
	RcvrGlnSprd     float64 `json:"RCVR_GLN_SPRD"`       // GLN Receiver 환율 spread
	SndrGlnSprdValR float64 `json:"SNDR_GLN_SPRD_VAL_R"` // GLN Sender 환율 spread 금액
	RcvrGlnSprdValR float64 `json:"RCVR_GLN_SPRD_VAL_R"` // GLN Receiver 환율 spread 금액
	RcvrFeeIz       []fee   `json:"RCVR_FEE_IZ"`         // 수수료 반복부

	//Sender Data
	SndrLcGlnUnqCd  string  `json:"SNDR_LC_GLN_UNQ_CD"`  //Sender Local GLN 고유코드 (SubKey)
	SndrNatCd       string  `json:"SNDR_NAT_CD"`         //Sender 국가 코드
	SndrCurCd       string  `json:"SNDR_CUR_CD"`         //Sender 통화 코드
	SndrWdrwAmt     float64 `json:"SNDR_WDRW_AMT"`       //Sender 출금 금액
	RcvrSndrCriXchr float64 `json:"RCVR_SNDR_CRI_XCHR"`  //Receiver/Sender 기준 환율
	SndrXrFe        float64 `json:"SNDR_XR_FE"`          // Receiver 환율 수수료
	SndrXrFeValS    float64 `json:"SNDR_XR_FE_VAL_S"`    // Sender 환율 수수료 금액
	RcvrXrFeValS    float64 `json:"RCVR_XR_FE_VAL_S"`    // Receiver 환율 수수료 금액
	GlnSprd         float64 `json:"SNDR_GLN_SPRD"`       // GLN Receiver 환율 spread
	SndrGlnSprdValS float64 `json:"SNDR_GLN_SPRD_VAL_S"` // GLN Sender 환율 spread 금액
	RcvrGlnSprdValS float64 `json:"RCVR_GLN_SPRD_VAL_S"` // GLN Receiver 환율 spread 금액
	SndrFeeIz       []fee   `json:"SNDR_FEE_IZ"`         // 수수료 반복부
}

// FeeIz Ledger Struct
type fee struct {
	FeDvCd   string  `json:"FE_DV_CD"`   // 수수료 구분코드
	FeAmt    float64 `json:"FE_AMT"`     // 수수료 금액
	ExFeAmt  float64 `json:"EX_FE_AMT"`  // 환전 수수료 금액
	SndRcvDv string  `json:"SND_RCV_DV"` // Sender / Receiver 구분(S/R)
	BenfoCd  string  `json:"BENFO_CD"`   // 부담처코드(01:회원,02:직제휴사,03:gln제휴사,04:merchant)
	BuoCd    string  `json:"BUO_CD"`     // 수혜처코드(01:GLN int''l/02:Local GLN)
}

// Query JSON struct
type queryArgs struct {
	GlnDeNo      string `json:"GLN_DE_NO"`      // GLN 거래 번호
	ReqStartTime string `json:"REQ_START_TIME"` // 기간 시작값
	ReqEndTime   string `json:"REQ_END_TIME"`   // 기간 끝 값
	LcGlnUnqCd   string `json:"LC_GLN_UNQ_CD"`  // Local GLN 코드
}

//Query Response Field
var commonField = "\"GLN_DE_NO\", \"GLN_DE_DTM\", \"MBR_SVC_DV_CD\", \"DE_TP_DV_CD\", \"RECV_DE_NO\", \"RECV_UNQ_CD\", \"ORI_GLN_DE_NO\", \"CNCL_GLN_DE_NO\", \"PBLD_DTM\",\"PBLD_TN\""
var rcvrField = "\"RCVR_LC_GLN_UNQ_CD\", \"RCVR_NAT_CD\", \"RCVR_CUR_CD\" ,\"RCVR_DEPO_AMT\" ,\"RCVR_SNDR_GLN_XCHR\", \"RCVR_XR_FE\", \"SNDR_XR_FE_VAL_R\", \"RCVR_XR_FE_VAL_R\", \"RCVR_GLN_SPRD\", \"SNDR_GLN_SPRD_VAL_R\", \"RCVR_GLN_SPRD_VAL_R\", \"RCVR_FEE_IZ\""
var sndrField = "\"SNDR_LC_GLN_UNQ_CD\", \"SNDR_NAT_CD\", \"SNDR_CUR_CD\", \"SNDR_WDRW_AMT\", \"RCVR_SNDR_CRI_XCHR\", \"SNDR_XR_FE\", \"SNDR_XR_FE_VAL_S\", \"RCVR_XR_FE_VAL_S\", \"SNDR_GLN_SPRD\", \"SNDR_GLN_SPRD_VAL_S\", \"RCVR_GLN_SPRD_VAL_S\", \"SNDR_FEE_IZ\""
