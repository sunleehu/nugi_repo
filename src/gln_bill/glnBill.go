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

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type glnBillCC struct {
}

var pageSize int32 = 100

var logger = shim.NewLogger("gln_billChaincode")

func main() {
	err := shim.Start(new(glnBillCC))
	if err != nil {
		fmt.Printf("Error starting setlLog chaincode: %s", err)
	}
}

func (t *glnBillCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *glnBillCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)

	//  Handle different functions
	if function == "putsttlbill" {
		return t.putBill(stub, args)
	} else if function == "getsttlbill" {
		return t.getBill(stub, args)
	} else if function == "getstllbillhistory" {
		return t.getBillHistory(stub, args)
	} else if function == "confirmsttlbill" {
		return t.confirmBill(stub, args)
	} else if function == "setLogLevel" {
		return setLogLevel(args[0])
	}
	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *glnBillCC) putBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Empty Argument Check
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}
	//event payload data map
	evtMap := make(map[string][]string)
	var pyld hEvt
	evtCheck := false

	// Identity Check
	err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}
	txID := stub.GetTxID()
	var validData [][]byte
	var keyList []string

	// validation loop
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
		if len(checkBlank(bill.AdjReqNo)) == 0 || len(checkBlank(bill.SndrLocalGlnCd)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Check ADJ_REQ_NO or SndrLocalGlnCd in JSON"))
		}

		// Json Encoding
		glnBillJSONBytes, err := json.Marshal(bill)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}
		//add key level
		var callargs []string
		callargs = append(callargs, bill.AdjReqNo, endorserMsp, cdToMSP(bill.SndrLocalGlnCd))
		_, errm := addOrgs(stub, callargs)
		if errm != "" {
			return shim.Error(errMessage("BCCE0011", errm))
		}

		//Event JSON
		evtMap[bill.SndrLocalGlnCd] = append(evtMap[bill.SndrLocalGlnCd], bill.AdjReqNo)

		keyList = append(keyList, bill.AdjReqNo)
		validData = append(validData, glnBillJSONBytes)
		pyld.Target = append(pyld.Target, bill.SndrLocalGlnCd)
	}

	// Duplicate Value Check in couchDB
	mulQuery := multiQueryMaker("ADJ_REQ_NO", keyList)
	queryString := fmt.Sprintf(`{"selector":{%s}, "fields":[%s]}`, mulQuery, `"ADJ_REQ_NO","TX_ID"`)

	fmt.Println(queryString)
	exs, res, err := isExist(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if exs {
		if err != nil {
			return shim.Error(errMessage("BCCE0008", err))
		}
		return shim.Error(errMessage("BCCE0006", fmt.Sprintf("%s", res)))
	}

	// putState loop
	for i := 0; i < len(validData); i++ {
		// Write couchDB
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

		logger.Info("SETTLEMENT_BILL_SAVED")
		logger.Debug("SETTLEMENT_BILL_SAVED", string(dat))
		// EVENT!!!
		stub.SetEvent("SETTLEMENT_BILL_SAVED", dat)

	}

	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *glnBillCC) getBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// Json Decoding
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
		//for international GLN check query argument check
		if len(checkBlank(qArgs.LcGlnUnqCd)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Check your LOCALGLN_CODE in JSON"))
		}

	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}
	}

	// Empty Value Check
	if len(checkBlank(qArgs.AdjReqNo)) == 0 {
		return t.getBillHistory(stub, args)
	}

	//Default Size 100
	// var pgs int32
	// if qArgs.PageSize > 0 {
	// 	pgs = qArgs.PageSize
	// } else {
	// 	pgs = pageSize
	// }

	// Query
	queryString := fmt.Sprintf(`{"selector": {"ADJ_REQ_NO": "%s", "LOCAL_GLN_CD":"%s"}}`, qArgs.AdjReqNo, qArgs.LcGlnUnqCd)
	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, 1, "")
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by International GLN
func (t *glnBillCC) getBillHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type error, invalid json(empty args)
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	attr, m := checkGlnIntl(stub)
	if m != "" {
		return shim.Error(errMessage("BCCE0005", m))
	}
	if attr {
	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}

	}

	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStartTime and ReqEndTime"))
	}
	//Default Size 100
	var pgs int32
	if qArgs.PageSize > 0 {
		pgs = qArgs.PageSize
	} else {
		pgs = pageSize
	}

	// Query
	queryString := fmt.Sprintf(`{"selector": {"$and":[{"LOCAL_GLN_CD": "%s"},{"ADJ_DT":{"$gte": "%s"}},{"ADJ_DT":{"$lte": "%s"}}]}}`, qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

func (t *glnBillCC) confirmBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// Json Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	attr, m := checkGlnIntl(stub)
	if m != "" {
		return shim.Error(errMessage("BCCE0005", m))
	}
	if attr {

	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}
	}

	queryString := fmt.Sprintf(`{"selector": {"ADJ_REQ_NO": "%s","LOCAL_GLN_CD":"%s"}}`, qArgs.AdjReqNo, qArgs.LcGlnUnqCd)
	logger.Debug("QueryString:", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if !resultsIterator.HasNext() {
		return shim.Error(errMessage("BCCE0010", fmt.Sprintf("Data %s", args[0])))
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return shim.Error(errMessage("BCCE0008", err))
		}

		var bill glnbill

		json.Unmarshal(queryResponse.Value, &bill)
		logger.Debug("QueryResponse:", queryResponse)
		// Update Value
		bill.SndrAdjDfnYn = "Y"

		jtx, err := json.Marshal(bill)

		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// Update CouchDB
		err = stub.PutState(queryResponse.Key, jtx)
		if err != nil {
			return shim.Error(errMessage("BCCE0010", err))
		}
		defer resultsIterator.Close()
	}
	return shim.Success(nil)
}
