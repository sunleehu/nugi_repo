package main

// Transaction Log Ledger Struct
type transaction struct {
	GlnTxNo          string    `json:"GLN_TX_NO"`
	GlnTxHash        string    `json:"GLN_TX_HASH"` // GLN 거래 번호 HASH
	SndrLocalGlnCd   string    `json:"SNDR_LOCAL_GLN_CD"`
	RcvrLocalGlnCd   string    `json:"RCVR_LOCAL_GLN_CD"`
	UsoLocalGlnCcoCd string    `json:"USO_LOCAL_GLN_CCO_CD"`
	TxTpDvCd         string    `json:"TX_TP_DV_CD"`
	CanDvCd          string    `json:"CAN_DV_CD"`
	OriGlnTxNo       string    `json:"ORI_GLN_TX_NO"`
	TxStCd           string    `json:"TX_ST_CD"`
	McNm             string    `json:"MC_NM"`
	McBtCd           string    `json:"MC_BT_CD"`
	McID             string    `json:"MC_ID"`
	McAdr            string    `json:"MC_ADR"`
	TerlID           string    `json:"TERL_ID"`
	RcptNo           string    `json:"RCPT_NO"`
	SndrXchrPbldNo   string    `json:"SNDR_XCHR_PBLD_NO"`
	RcvrXchrPbldNo   string    `json:"RCVR_XCHR_PBLD_NO"`
	TxApyXchr        float64   `json:"TX_APY_XCHR"`
	AdjCriDt         string    `json:"ADJ_CRI_DT"`
	UtcTxDtm         string    `json:"UTC_TX_DTM"`
	SndrTxDtm        string    `json:"SNDR_TX_DTM"`
	RcvrTxDtm        string    `json:"RCVR_TX_DTM"`
	RcvrCurTxPrcp    float64   `json:"RCVR_CUR_TX_PRCP"`
	UsdTxPrcp        float64   `json:"USD_TX_PRCP"`
	RcvrAdjAmt       float64   `json:"RCVR_ADJ_AMT"`
	SndrAdjAmt       float64   `json:"SNDR_ADJ_AMT"`
	UsdMbrPayAmt     float64   `json:"USD_MBR_PAY_AMT"`
	SndrCurMbrPayAmt float64   `json:"SNDR_CUR_MBR_PAY_AMT"`
	TxID             string    `json:"TX_ID"`
	FeeData          []feeList `json:"FEE_DATA"`
	
	//2019.10.31 이선혁 added
	OriRcptNo			string `json:"ORI_RCPT_NO"`
	RcvrCurAdjAmt		string `json:"RCVR_CUR_ADJ_AMT"`
	SndrCurAdjAmt		string `json:"SNDR_CUR_ADJ_AMT"`
}

type feeList struct {
	FeDvCd   string  `json:"FE_DV_CD"`
	BenfoCd  string  `json:"BENFO_CD"`
	BuoCd    string  `json:"BUO_CD"`
	FeUsdAmt float64 `json:"FE_USD_AMT"`
}

type pubData struct {
	GlnTxHash string `json:"GLN_TX_HASH"` // GLN 거래 번호 HASH
	Date      string `json:"DATE"`        // 거래일시
	From      string `json:"FROM"`        // Sender
	To        string `json:"TO"`          // Receiver
	BcTxID    string `json:"TX_ID"`       // 블록체인 TX ID
}

type respStruct struct {
	GlnTxNo string `json:"GLN_TX_NO"`
	BcTxID  string `json:"TX_ID"`
}

// Query JSON struct
type queryArgs struct {
	GlnTxNo      string `json:"GLN_TX_NO"`      // GLN 거래
	ReqStartTime string `json:"REQ_START_TIME"` // 기간 시작값
	ReqEndTime   string `json:"REQ_END_TIME"`   // 기간 끝 값
	LcGlnUnqCd   string `json:"LOCALGLN_CODE"`  // Local GLN 코드
	DivCd        string `json:"DIV_CODE"`       // 구분 코드
	DeTpDvCd     string `json:"TX_TP_DV_CD"`    // 정상 취소 구분 코드
	PageSize     int32  `json:"PAGE_COUNT"`     //Paging size
	BookMark     string `json:"PAGE_NEXT_ID"`   //Bookmark ID
}

// Event Payload Header Json
type hEvt struct {
	Target []string
	Data   interface{}
}

// type queryResponseStruct struct {
// }

const Reciver string = "RCVR_LOCAL_GLN_CD" //
const Sender string = "SNDR_LOCAL_GLN_CD"  //
