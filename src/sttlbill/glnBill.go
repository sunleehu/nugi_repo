/*
* This Chaincode is GLN Settlement Log code
* And it has functions insert and query,
* International GLN can insert Settlement Log.
* Interational and Local GLN can query Data
**/
package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type glnBillCC struct {
}

var logger = shim.NewLogger("STTLBILL")
var defaultPageSize int32 = 100

func main() {
	logger.SetLevel(shim.LogDebug)
	err := shim.Start(new(glnBillCC))
	if err != nil {
		logger.Error("Error starting sttlbill chaincode : %s", err)
	}
}

func (t *glnBillCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *glnBillCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println()
	logger.Info("Func :", function)
	logger.Info("Args : ", args)

	//  Handle different functions
	if function == "putsttlbill" {
		return t.putBill(stub, args)
	} else if function == "getsttlbill" {
		return t.getBill(stub, args)
	} else if function == "getsttlbillhistory" {
		return t.getBillHistory(stub, args)
	} else if function == "confirmsttlbill" {
		return t.confirmBill(stub, args)
	} else if function == "delstate" {
		return t.deleteState(stub, args)
	} else if function == "delstatehistory" {
		return t.deleteStateHistory(stub, args)
	}
	// } else if function == "setLogLevel" {
	// 	return setLogLevel(args[0])
	// }
	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *glnBillCC) putBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	logger.Info("Put Data Count : ", len(args))

	// Check arguments
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	// Identity Check
	err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	txID := stub.GetTxID()
	var validData [][]byte
	var keyList []string
	evtMap := make(map[string][]string)
	var pyld hEvt

	// Validation loop
	for k := 0; k < len(args); k++ {
		var bill glnbill

		// Json Decoding
		err := json.Unmarshal([]byte(args[k]), &bill)
		if err != nil {
			return shim.Error(errMessage("BCCE0003", err))
		}

		//TX ID
		bill.Txid = txID
		// Empty Value Check
		if isBlank(bill.AdjPblNo) || isBlank(bill.SndrLocalGlnCd) {
			return shim.Error(errMessage("BCCE0005", "Check ADJ_PBL_NO or SndrLocalGlnCd in JSON"))
		}

		// Json Encoding
		glnBillJSONBytes, err := json.Marshal(bill)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// // Add key level
		// var callargs []string
		// callargs = append(callargs, bill.AdjPblNo, endorserMsp, cdToMSP(bill.SndrLocalGlnCd))
		// _, errm := addOrgs(stub, callargs)
		// if errm != "" {
		// 	return shim.Error(errMessage("BCCE0012", errm))
		// }

		// Event JSON
		evtMap[bill.SndrLocalGlnCd] = append(evtMap[bill.SndrLocalGlnCd], bill.AdjPblNo)

		keyList = append(keyList, bill.AdjPblNo)
		validData = append(validData, glnBillJSONBytes)
		pyld.Target = append(pyld.Target, bill.SndrLocalGlnCd)
	}

	// Duplicate Value Check in couchDB
	exs, res, err := isExistByKey(stub, keyList)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if exs {
		return shim.Error(errMessage("BCCE0006", fmt.Sprintf("%s", res)))
	}

	// PutState loop
	evtCheck := false
	for i := 0; i < len(validData); i++ {
		err = stub.PutState(keyList[i], validData[i])
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}
		evtCheck = true
	}

	if evtCheck {
		// Event Payload JSON Encoding
		pyld.Target = rmvDupVal(pyld.Target)
		pyld.Data = evtMap
		dat, e := json.Marshal(pyld)
		if e != nil {
			return shim.Error(errMessage("BCCE0004", e))
		}
		// EVENT!!!
		stub.SetEvent("SETTLEMENT_BILL_SAVED", dat)
		logger.Infof("[EVENT] SETTLEMENT_BILL_SAVED :\n%s", string(dat))
	}

	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *glnBillCC) getBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check arguments
	if len(args) < 1 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}
	var qArgs queryArgs
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Identity Check
	attr, m := checkGlnIntl(stub)
	if m != "" {
		return shim.Error(errMessage("BCCE0005", m))
	}
	if attr {
		// //for international GLN check query argument check
		// if len(checkBlank(qArgs.LcGlnUnqCd)) == 0 {
		// 	return shim.Error(errMessage("BCCE0005", "Check your LOCALGLN_CODE in JSON"))
		// }
	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}
	}
	//ADJ_PBL_NO 의 값이 없으면 
	if isBlank(qArgs.AdjPblNo) {
		return t.getBillHistory(stub, args)
	}

	// Query
	// queryString := fmt.Sprintf(`{"selector": {"ADJ_PBL_NO": "%s", "LOCAL_GLN_CD":"%s"}}`, qArgs.AdjPblNo, qArgs.LcGlnUnqCd)
	// queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, 1, "")
	state, err := stub.GetState(qArgs.AdjPblNo)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if state == nil {
		resp := queryResponseStructMaker(nil, "", 0)
		logger.Info("Query Complete")
		logger.Info(string(resp))
		return shim.Success(resp)
	}

	var resList [][]byte
	resList = append(resList, state)
	queryResults := queryResponseStructMaker(resList, "", 1)

	logger.Info("Query Complete")
	logger.Info(string(queryResults))
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by International GLN
func (t *glnBillCC) getBillHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check arguments
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}
	var qArgs queryArgs
	
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStartTime and ReqEndTime"))
	}
	if len(strings.TrimSpace(qArgs.ReqStartTime)) != 8 || len(strings.TrimSpace(qArgs.ReqEndTime)) != 8 {
		return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDD"`))
	}
	
	qArgs.BpLocalGlnCd = qArgs.LcGlnUnqCd
	
	// Check Identity
	attr, m := checkGlnIntl(stub)
	if m != "" {
		return shim.Error(errMessage("BCCE0005", m))
	}

	qArgs.BpLocalGlnCd = qArgs.LcGlnUnqCd

	if attr {
	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)//LcGlnUnqCd 여기서만 사용 
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match. LclGlnUnqCd: " + qArgs.LcGlnUnqCd + " qArgs.BpLocalGlnCd:"+ qArgs.BpLocalGlnCd))
		}
	}

	// Page Size
	pgs := qArgs.PageSize
	if pgs == 0 {
		pgs = defaultPageSize
	}

	// Query
	var queryString string
	// AS-IS BEFORE 2019.12.27 if qArgs.LcGlnUnqCd == "" 조건 변경 이미 빈값이면 바로 전단계에서 끝남 
	// TO-BE SINCE 2019.12.27 조회 SEL_SSP_CD(gateway 용 필드 ) 빈값이면 SSP 조회가 아니므로 lcoal_gln_cd 가 해당하는 건만 보여준다. 
	if qArgs.SpLocalGlnCd == "" {
		// queryString = fmt.Sprintf(`{"selector": {"$and":[{"ADJ_PBL_DT":{"$gte": "%s"}},{"ADJ_PBL_DT":{"$lte": "%s"}}]}, "use_index":["indexDateDoc", "indexDate"]}`, qArgs.ReqStartTime, qArgs.ReqEndTime)
		queryString = fmt.Sprintf(`{"selector": {"$and":[{"LOCAL_GLN_CD": "%s"},{"ADJ_PBL_DT":{"$gte": "%s"}},{"ADJ_PBL_DT":{"$lte": "%s"}}]}, "use_index":["indexDateLclDoc", "indexDateLcl"]}`, qArgs.BpLocalGlnCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	} else {
		// queryString = fmt.Sprintf(`{"selector": {"$and":[{"LOCAL_GLN_CD": "%s"},{"ADJ_PBL_DT":{"$gte": "%s"}},{"ADJ_PBL_DT":{"$lte": "%s"}}]}, "use_index":["indexDateLclDoc", "indexDateLcl"]}`, qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
		queryString = fmt.Sprintf(`{"selector": {"$and":[{"BP_LOCAL_GLN_CD":"%s"},{"LOCAL_GLN_CD": "%s"},{"ADJ_PBL_DT":{"$gte": "%s"}},{"ADJ_PBL_DT":{"$lte": "%s"}}]}, "use_index":["indexDateLclDoc", "indexDateLcl"]}`, qArgs.BpLocalGlnCd, qArgs.SpLocalGlnCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	}
	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	logger.Info(string(queryResults))
	return shim.Success(queryResults)
}

func (t *glnBillCC) confirmBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check arguments
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}
	var qArgs queryArgs
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	adjPblNo := strings.TrimSpace(qArgs.AdjPblNo)
	lcGlnUnqCd := strings.TrimSpace(qArgs.LcGlnUnqCd)
	if adjPblNo == "" || lcGlnUnqCd == "" {
		shim.Error(errMessage("BCCE0005", "Need arguments ADJ_PBL_NO, LOCAL_GLN_CD in JSON"))
	}

	// Check Identities
	attr, m := checkGlnIntl(stub)
	if m != "" {
		return shim.Error(errMessage("BCCE0005", m))
	}
	if attr {

	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", lcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}
	}

	// Query
	data, err := stub.GetState(adjPblNo)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if data == nil {
		return shim.Error(errMessage("BCCE0008", fmt.Sprintf("No data by %s", adjPblNo)))
	}
	logger.Debug("QueryResponse : ", string(data))

	// JSON Decoding
	var bill glnbill
	err = json.Unmarshal(data, &bill)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if bill.SndrLocalGlnCd != lcGlnUnqCd {
		return shim.Error(errMessage("BCCE0008", fmt.Sprintf("ADJ_PBL_NO[%s]'s LOCAL_GLN_CD is not match %s", adjPblNo, lcGlnUnqCd)))
	}

	// Update Value
	bill.SndrAdjDfnYn = "Y"
	jtx, err := json.Marshal(bill)
	if err != nil {
		return shim.Error(errMessage("BCCE0004", err))
	}
	err = stub.PutState(adjPblNo, jtx)
	if err != nil {
		return shim.Error(errMessage("BCCE0010", err))
	}

	logger.Info("Update Complete")
	logger.Info(string(jtx))
	return shim.Success(nil)
}

func (t *glnBillCC) deleteState(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check Arguments
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}
	var qArgs queryArgs
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	if isBlank(qArgs.AdjPblNo) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string ADJ_PBL_NO"))
	}

	// Check identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Delete State
	err = stub.DelState(qArgs.AdjPblNo)
	if err != nil {
		return shim.Error(errMessage("BCCE0011", err))
	}

	// // Delete Key-Level Endorsement
	// err = delEP(stub, qArgs.AdjPblNo)
	// if err != nil {
	// 	return shim.Error(errMessage("BCCE0011", err))
	// }

	logger.Info("Delete Complete")
	return shim.Success(nil)
}

func (t *glnBillCC) deleteStateHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check Arguments
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}
	var qArgs queryArgs
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}
	if len(strings.TrimSpace(qArgs.ReqStartTime)) != 8 || len(strings.TrimSpace(qArgs.ReqEndTime)) != 8 {
		return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDD"`))
	}

	// Check identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Query
	var queryString string
	if qArgs.LcGlnUnqCd == "" {
		queryString = fmt.Sprintf(`{"selector": {"$and":[{"ADJ_PBL_DT":{"$gte": "%s"}},{"ADJ_PBL_DT":{"$lte": "%s"}}]}, "use_index":["indexDateDoc", "indexDate"]}`, qArgs.ReqStartTime, qArgs.ReqEndTime)
	} else {
		queryString = fmt.Sprintf(`{"selector": {"$and":[{"LOCAL_GLN_CD": "%s"},{"ADJ_PBL_DT":{"$gte": "%s"}},{"ADJ_PBL_DT":{"$lte": "%s"}}]}, "use_index":["indexDateLclDoc", "indexDateLcl"]}`, qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	}
	dataList, err := getDataByQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	// Delete
	count := len(dataList)
	for i := 0; i < count; i++ {
		err = stub.DelState(dataList[i].AdjPblNo)
		if err != nil {
			return shim.Error(errMessage("BCCE0011", err))
		}
		// err = delEP(stub, dataList[i].AdjPblNo)
		// if err != nil {
		// 	return shim.Error(errMessage("BCCE0011", err))
		// }
	}

	logger.Info("Delete Complete -", count)
	return shim.Success(nil)
}
