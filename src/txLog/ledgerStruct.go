package main

// Transaction Log Ledger Struct
type transaction struct {
	GlnDeNo        string  `json:"GLN_DE_NO"`          // GLN 거래 번호
	Seq            uint64  `json:"Seq"`                // Seq
	GlnMbrUnqk     string  `json:"GLN_MBR_UNQK"`       // 회원 Local GLN 회원 식별키
	SndrLcGlnUnqCd string  `json:"SNDR_LC_GLN_UNQ_CD"` // 회원 Local GLN 고유코드
	MbrSvcDvCd     string  `json:"MBR_SVC_DV_CD"`      // 회원 서비스 구분 코드
	DeTpDvCd       string  `json:"DE_TP_DV_CD"`        // 정상 취소 구분 코드
	GlnDeDtm       string  `json:"GLN_DE_DTM"`         // GLN 거래 일시
	RcvrLcGlnUnqCd string  `json:"RCVR_LC_GLN_UNQ_CD"` // 사용 Local GLN 고유코드
	UsoNm          string  `json:"USO_NM"`             // 사용처 명
	RcvrCurCd      string  `json:"RCVR_CUR_CD"`        // 사용 통화 코드
	RcvrDepoAmt    float64 `json:"RCVR_DEPO_AMT"`      // 사용 금액(상품 금액)
	SndrCurCd      string  `json:"SNDR_CUR_CD"`        // 회원 통화 코드
	SndrWdrwAmt    float64 `json:"SNDR_WDRW_AMT"`      // 회원 출금 금액
	AdCtt          string  `json:"AD_CTT"`             // 추가내용
	RsvAtc         string  `json:"RSV_ATC,omitempty"`  // reserved 필드
}

// Query JSON struct
type queryArgs struct {
	GlnDeNo      string `json:"GLN_DE_NO"`      // GLN 거래
	ReqStartTime string `json:"REQ_START_TIME"` // 기간 시작값
	ReqEndTime   string `json:"REQ_END_TIME"`   // 기간 끝 값
	LcGlnUnqCd   string `json:"LC_GLN_UNQ_CD"`  // Local GLN 코드
	GlnMbrUnqk   string `json:"GLN_MBR_UNQK"`   // GLN 회원 유니크 키
	DeTpDvCd     string `json:"DE_TP_DV_CD"`    // 정상 취소 구분 코드
}
