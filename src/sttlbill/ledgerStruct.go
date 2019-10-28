package main

//Settlement Log Ledger Struct

type glnbill struct {
	AdjPblNo           string `json:"ADJ_PBL_NO"`
	SndrLocalGlnCd     string `json:"LOCAL_GLN_CD"`
	AdjSDt             string `json:"ADJ_S_DT"`
	AdjEDt             string `json:"ADJ_E_DT"`
	AdjPblTxCnt        uint64 `json:"ADJ_PBL_TX_CNT"`
	AdjPblDt           string `json:"ADJ_PBL_DT"`
	TxPrcpSum          string `json:"TX_PRCP_SUM"`
	TxPrcpSumSign      string `json:"TX_PRCP_SUM_SIGN"`
	SndrAdjAmt         string `json:"SNDR_ADJ_AMT"`
	SndrAdjAmtSign     string `json:"SNDR_ADJ_AMT_SIGN"`
	RcvrAdjAmt         string `json:"RCVR_ADJ_AMT"`
	RcvrAdjAmtSign     string `json:"RCVR_ADJ_AMT_SIGN"`
	GlnAdjAmt          string `json:"GLN_ADJ_AMT"`
	GlnAdjAmtSign      string `json:"GLN_ADJ_AMT_SIGN"`
	SndrRcvgFeSum      string `json:"SNDR_RCVG_FE_SUM"`
	SndrRcvgFeSumSign  string `json:"SNDR_RCVG_FE_SUM_SIGN"`
	RcvrRcvgFeSum      string `json:"RCVR_RCVG_FE_SUM"`
	RcvrRcvgFeSumSign  string `json:"RCVR_RCVG_FE_SUM_SIGN"`
	GlnRcvgFeSum       string `json:"GLN_RCVG_FE_SUM"`
	GlnRcvgFeSumSign   string `json:"GLN_RCVG_FE_SUM_SIGN"`
	GlnAdjBnkSwiftCd   string `json:"GLN_ADJ_BNK_SWIFT_CD"`
	GlnAdjBnkNm        string `json:"GLN_ADJ_BNK_NM"`
	GlnAdjAcNo         string `json:"GLN_ADJ_AC_NO"`
	SndrAdjDfnYn       string `json:"LOCAL_GLN_ADJ_DFN_YN"`
	TotalFeeSign       string `json:"TOTAL_FEE_SIGN"`
	TotalFeeAmount     string `json:"TOTAL_FEE_AMOUNT"`
	SettlementFileName string `json:"SETTLEMENT_FILE_NAME"`
	FeeFileName        string `json:"FEE_FILE_NAME"`
	Txid               string `json:"TX_ID"`
	//2019.10.16 이선혁 추가 
	SndrTxCnt		   string `json:"SNDR_TX_CNT"`
	RcvrTxCnt 		   string `json:"RCVR_TX_CNT"`
	BpLocalGlnCd	   string `json:"BP_LOCAL_GLN_CD"`
	SpAdjCurCd		   string `json:"SP_ADJ_CUR_CD"`
	SndrCurAdjAmt	   string `json:"SNDR_CUR_ADJ_AMT"`
	RcvrCurAdjAmt	   string `json:"RCVR_CUR_ADJ_AMT"`
	SndrCurFeAdjAmt	   string `json:"SNDR_CUR_FE_ADJ_AMT"`
	RcvrCurFeAdjAmt	   string `json:"RCVR_CUR_FE_ADJ_AMT"`
	AdjFileNm 		   string `json:"ADJ_FILE_NM"`
	FeFileNm		   string `json:"FE_FILE_NM"`
	BlcTxId			   string `json:"BLC_TX_ID"`
	RgDtm			   string `json:"RG_DTM"`
	RgrId			   string `json:"RGR_ID"`
	ChDtm			   string `json:"CH_DTM"`
	ChrId			   string `json:"CHR_ID"`

}

// Query JSON struct
type queryArgs struct {
	AdjPblNo     string `json:"ADJ_PBL_NO"`     //정산요청번호
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

const endorserMsp = "EndorserMSP"
const channelID = "glnchannel"
const libEp = "libep"
