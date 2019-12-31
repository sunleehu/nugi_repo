package main

//Settlement Log Ledger Struct

type glnbill struct {
	AdjPblNo           string `json:"ADJ_PBL_NO"`			//정산요청번호
	SndrLocalGlnCd     string `json:"LOCAL_GLN_CD"`			//LocalGLN코드
	AdjSDt             string `json:"ADJ_S_DT"`				//정산시작일자
	AdjEDt             string `json:"ADJ_E_DT"`				//정산종료일자
	AdjPblTxCnt        uint64 `json:"ADJ_PBL_TX_CNT"`		//정산거래건수
	AdjPblDt           string `json:"ADJ_PBL_DT"`			//정산일자
	TxPrcpSum          string `json:"TX_PRCP_SUM"`			//거래원금합계
	TxPrcpSumSign      string `json:"TX_PRCP_SUM_SIGN"`		//거래원금합계부호
	SndrAdjAmt         string `json:"SNDR_ADJ_AMT"`			//Sender정산금액
	SndrAdjAmtSign     string `json:"SNDR_ADJ_AMT_SIGN"`	//Sender정산금액부호
	RcvrAdjAmt         string `json:"RCVR_ADJ_AMT"`			//Receiver정산금액
	RcvrAdjAmtSign     string `json:"RCVR_ADJ_AMT_SIGN"`	//Receiver정산금액부호
	GlnAdjAmt          string `json:"GLN_ADJ_AMT"`			//GLN정산금액
	GlnAdjAmtSign      string `json:"GLN_ADJ_AMT_SIGN"`		//GLN정산금액부호
	SndrRcvgFeSum      string `json:"SNDR_RCVG_FE_SUM"`		//Sender수취수수료합계
	SndrRcvgFeSumSign  string `json:"SNDR_RCVG_FE_SUM_SIGN"`//Sender수취수수료합계부호
	RcvrRcvgFeSum      string `json:"RCVR_RCVG_FE_SUM"`		//Receiver수취수수료합계
	RcvrRcvgFeSumSign  string `json:"RCVR_RCVG_FE_SUM_SIGN"`//Receiver수취수수료합계부호
	GlnRcvgFeSum       string `json:"GLN_RCVG_FE_SUM"`		//GLN수취수수료합계
	GlnRcvgFeSumSign   string `json:"GLN_RCVG_FE_SUM_SIGN"`	//GLN수취수수료합계부호
	GlnAdjBnkSwiftCd   string `json:"GLN_ADJ_BNK_SWIFT_CD"`	//GLN정산은행 SWIFT코드
	GlnAdjBnkNm        string `json:"GLN_ADJ_BNK_NM"`		//GLN정산은행명
	GlnAdjAcNo         string `json:"GLN_ADJ_AC_NO"`		//GLN정산은행계좌번호
	SndrAdjDfnYn       string `json:"LOCAL_GLN_ADJ_DFN_YN"`	//정산확인여부
	TotalFeeSign       string `json:"TOTAL_FEE_SIGN"`		//전체수수료합계
	TotalFeeAmount     string `json:"TOTAL_FEE_AMOUNT"`		//전체수수료합계부호
	SettlementFileName string `json:"SETTLEMENT_FILE_NAME"`	//정산파일명
	FeeFileName        string `json:"FEE_FILE_NAME"`		//수수료파일명
	Txid               string `json:"TX_ID"`				//TX ID
	//2019.11.05 이선혁 변경
	RcvrAprAmt			string `json:"RCVR_APR_AMT"`			//Receiver승인금액
	RcvrAprAmtSign		string `json:"RCVR_APR_AMT_SIGN"`		//Receiver승인금액부호	
	RcvrCanAmt	 		string `json:"RCVR_CAN_AMT"`			//Receiver취소금액
	RcvrCanAmtSign	 	string `json:"RCVR_CAN_AMT_SIGN"`		//Receiver취소금액부호	
	RcvrChrbkCanAmt	 	string `json:"RCVR_CHRBK_CAN_AMT"`		//Receiver이의제기취소금액	
	RcvrChrbkCanAmtSign	string `json:"RCVR_CHRBK_CAN_AMT_SIGN"`	//Receiver이의제기취소금액부호			
	RcvrCurAdjAmt	 	string `json:"RCVR_CUR_ADJ_AMT"`		//Receiver통화정산금액	
	RcvrCurAdjAmtSign	string `json:"RCVR_CUR_ADJ_AMT_SIGN"`	//Receiver통화정산금액부호		
	RcvrCurFeAdjAmt	 	string `json:"RCVR_CUR_FE_ADJ_AMT"`		//Receiver통화수수료정산금액		
	RcvrCurFeAdjAmtSign	string `json:"RCVR_CUR_FE_ADJ_AMT_SIGN"`//Receiver통화수수료정산금액부호			
	RcvrRpmAmt	 		string `json:"RCVR_RPM_AMT"`			//Receiver환불금액
	RcvrRpmAmtSign	 	string `json:"RCVR_RPM_AMT_SIGN"`		//Receiver환불금액부호	
	SndrAprAmt	 		string `json:"SNDR_APR_AMT"`			//Sender승인금액
	SndrAprAmtSign	 	string `json:"SNDR_APR_AMT_SIGN"`		//Sender승인금액부호	
	SndrCanAmt	 		string `json:"SNDR_CAN_AMT"`			//Sender취소금액
	SndrCanAmtSign	 	string `json:"SNDR_CAN_AMT_SIGN"`		//Sender취소금액부호	
	SndrChrbkCanAmt	 	string `json:"SNDR_CHRBK_CAN_AMT"`		//Sender이의제기취소금액	
	SndrChrbkCanAmtSign	string `json:"SNDR_CHRBK_CAN_AMT_SIGN"`	//Sender이의제기취소금액부호			
	SndrCurAdjAmt	 	string `json:"SNDR_CUR_ADJ_AMT"`		//Sender통화정산금액	
	SndrCurAdjAmtSign	string `json:"SNDR_CUR_ADJ_AMT_SIGN"`	//Sender통화정산금액부호		
	SndrCurFeAdjAmt	 	string `json:"SNDR_CUR_FE_ADJ_AMT"`		//Sender통화수수료정산금액		
	SndrCurFeAdjAmtSign string `json:"SNDR_CUR_FE_ADJ_AMT_SIGN"`//Sender통화수수료정산금액부호			
	SndrRpmAmt	 		string `json:"SNDR_RPM_AMT"`			//Sender환불금액
	SndrRpmAmtSign	 	string `json:"SNDR_RPM_AMT_SIGN"`		//Sender환불금액부호	
	BpLocalGlnCd	 	string `json:"BP_LOCAL_GLN_CD"`			//BPLocalGLN코드	
	SpAdjCurCd	 		string `json:"SP_ADJ_CUR_CD"`			//SP정산통화코드

	
}

// Query JSON struct
type queryArgs struct {
	AdjPblNo     string `json:"ADJ_PBL_NO"`     //정산요청번호
	ReqStartTime string `json:"REQ_START_TIME"` // 기간 시작값
	ReqEndTime   string `json:"REQ_END_TIME"`   // 기간 끝 값
	LcGlnUnqCd   string `json:"LOCALGLN_CODE"`  // Local GLN 코드 -- 이건은 요청하는 gln_code 가 회원인지 확인할때만 사용하도록 변경 
	DivCd        string `json:"DIV_CODE"`       // 구분 코드
	DeTpDvCd     string `json:"DE_TP_DV_CD"`    // 정상 취소 구분 코드
	PageSize     int32  `json:"PAGE_COUNT"`
	BookMark     string `json:"PAGE_NEXT_ID"`
	SpLocalGlnCd string `json:"SEL_SSP_CD"`		//2019.12.27 추가 LOCAL_GLN_CODE 로 조회 가능하도록 ..
	BpLocalGlnCd string 	//2019.12.27 json에서 같은 값을 꺼내는 이유는 BP_LOCAL_CODE는 조회기관 KOEXKR 의 조건이 되어야 한다. 
	//2019.12.27 요건 SBP에서 SSP코드로 조회하는 기능 블록체인 입장에서 BP_LOCAL_CODE == MSP(KOEXKR)만 허용 
	//실제 조회 대상은 LOCAL_GLN_CD == SEL_SSP_CD(TOSSKR) 이다 
}

// Event Payload Header Json
type hEvt struct {
	Target []string
	Data   interface{}
}

const endorserMsp = "EndorserMSP"
const channelID = "glnchannel"
const libEp = "libep"
