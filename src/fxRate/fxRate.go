/*
* This Chaincode is GLN Exchange Rate code
* And it has functions insert and query,
* International GLN can insert Exchage Rate Data.
* Interational and Local GLN can query Data
**/
package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type fxRateCC struct {
}

// Logger
var logger = shim.NewLogger("fxRateCC")

var pageSize int32 = 100

func main() {
	err := shim.Start(new(fxRateCC))
	if err != nil {
		fmt.Printf("Error starting exRate chaincode: %s", err)
	}
}

func (t *fxRateCC) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

func (t *fxRateCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	// fmt.Println("invoke is running " + function)
	logger.Info("Invoke is running", function)

	// Handle different functions
	if function == "putxchrate" {
		return t.putXchRate(stub, args)
	} else if function == "deployxchrate" {
		// return t.deployXchRate(stub, args)
	} else if function == "getxchrate" {
		return t.getXchRate(stub, args)
	} else if function == "getglnxchrate" {
		// return t.getGlnXchRate(stub, args)
	} else if function == "getlatestglnxchr" {
		// return t.getLatestGlnXchr(stub, args)
	} else if function == "getxchratehistory" {
		return t.getXchRateHistory(stub, args)
	} else if function == "getglnxchratehistory" {
		// return t.getGlnXchRateHistory(stub, args)
	} else if function == "delstate" {
		return t.deleteState(stub, args)
	}
	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions. Called by Local GLN
// func (t *fxRateCC) putXchRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	// Emtpy Argument Check
// 	if len(args) == 0 {
// 		return shim.Error(errMessage("BCCE0007", "Args is empty"))
// 	}

// 	var lcXch localExchangeRate
// 	var latestLcXch latestLocalXCHRInfo
// 	txTime, err := stub.GetTxTimestamp()
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0001", "Couldn't get Tx Timestamp"))
// 	}

// 	// JSON Decoding
// 	err = json.Unmarshal([]byte(args[0]), &lcXch)
// 	if err != nil {
// 		// err case: type error, invalid json
// 		return shim.Error(errMessage("BCCE0003", err))
// 	}

// 	//check identity
// 	// err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", lcXch.LcGlnUnqCd)
// 	// if err != nil {
// 	// 	return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
// 	// }

// 	//check Exchange Rate
// 	if lcXch.UsdBidr <= 0 || lcXch.UsdOfferr <= 0 {
// 		return shim.Error(errMessage("BCCE0007", "Exchange rate can not be less than zero"))
// 	}

// 	//timestamp
// 	lcXch.RgDtm = getTimestamp(txTime.Seconds)

// 	//Create Unq code
// 	unqKey := lcXch.LcGlnCd + lcXch.RgDtm
// 	lcXch.LocalGlnXchrInfUnqno = unqKey
// 	// Local Exchange Rate JSON encoding
// 	exRateJSONBytes, err := json.Marshal(lcXch)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0004", err))
// 	}
// 	// Write Ledger Local GLN Exchange rate Info
// 	err = stub.PutState(unqKey, exRateJSONBytes)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0010", err))
// 	}
// 	logger.Info("Insert Complete")

// 	latestLcXch.localExchangeRate = lcXch
// 	latestLcXch.DocType = "latestLocalXchr"
// 	latestUnqKey := "LatestLocal" + latestLcXch.LcGlnCd
// 	// Latest Exchange Rate JSON encoding
// 	latestExRate, err := json.Marshal(latestLcXch)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0004", err))
// 	}

// 	// Renew Latest Local GLN Exchange rate Info
// 	err = stub.PutState(latestUnqKey, latestExRate)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0010", err))
// 	}
// 	logger.Info("Latest Value Renewal Complete")

// 	return shim.Success(nil)
// }

// //Calculate Local Exchange Rate info to GLN Exchange Rate
// func (t *fxRateCC) deployXchRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	// Emtpy Argument Check
// 	if len(args) == 0 {
// 		return shim.Error(errMessage("BCCE0007", "Args is empty"))
// 	}

// 	var lcXch []localExchangeRate
// 	var depArgs deploymentArgs
// 	var depXchr deployedXchRate
// 	var glnXchr glnExchangeRate

// 	// JSON Decoding
// 	err := json.Unmarshal([]byte(args[0]), &depArgs)
// 	if err != nil {
// 		// err case: type error, invalid json
// 		return shim.Error(errMessage("BCCE0003", err))
// 	}

// 	//check identity
// 	// err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
// 	// if err != nil {
// 	// 	return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
// 	// }
// 	queryString := fmt.Sprintf(`{"selector": {"XCHR_UNQNO":"%s"}}`, depArgs.XchrUnqno)
// 	exs, err := isExist(stub, queryString)
// 	if exs {
// 		return shim.Error(errMessage("BCCE0006", fmt.Sprintf("Data %s", args[0])))
// 	}

// 	queryString = fmt.Sprintf(`{"selector": {"DOC_TYPE":"latestLocalXchr"}}`)
// 	// Query
// 	queryResults, err := getQueryResultForQueryString(stub, queryString)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0008", err))
// 	}
// 	//Length 0 error func insert Here
// 	logger.Info("Query Success")

// 	err = json.Unmarshal(queryResults, &lcXch)
// 	if err != nil {
// 		// err case: type error, invalid json
// 		return shim.Error(errMessage("BCCE0003", err))
// 	}
// 	fmt.Println(lcXch)
// 	fmt.Println(string(queryResults))
// 	fmt.Println(depArgs)

// 	sprdList := make(map[string]sprdInfo)
// 	for i := 0; i < len(depArgs.SprdList); i++ {
// 		sprdList[depArgs.SprdList[i].LcGlnCd] = depArgs.SprdList[i]
// 		fmt.Println(depArgs)
// 	}

// 	for j := 0; j < len(lcXch); j++ {
// 		if len(lcXch[j].LcGlnCd) > 0 {
// 			depXchr.localExchangeRate = lcXch[j]
// 			sprd := sprdList[lcXch[j].LcGlnCd]
// 			fmt.Println(sprd)
// 			depXchr.GlnUsdBidr = decimalTrunc(decimalMultiply(depXchr.UsdBidr, decimalSub(1, sprd.RcvrSpr)), 6)
// 			depXchr.GlnUsdOfferr = decimalCeil(decimalMultiply(depXchr.UsdOfferr, decimalAdd(1, sprd.SndrSpr)), 6)
// 			fmt.Println("dex", depXchr)
// 			glnXchr.DeployL = append(glnXchr.DeployL, depXchr)
// 		}
// 	}

// 	//GLN FX
// 	glnXchr.pbldInfo = depArgs.pbldInfo
// 	glnXchrJSONByte, err := json.Marshal(glnXchr)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0004", err))
// 	}
// 	fmt.Println(string(glnXchrJSONByte))
// 	err = stub.PutState(glnXchr.XchrUnqno, glnXchrJSONByte)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0010", err))
// 	}

// 	var latestGlnXchr latestGlnXCHRInfo
// 	latestGlnXchr.DocType = "latestGlnXchr"
// 	latestGlnXchr.glnExchangeRate = glnXchr
// 	fmt.Println("pbldinfo", depArgs.pbldInfo)
// 	fmt.Println(glnXchr)
// 	lateGlnXchrJSONByte, err1 := json.Marshal(latestGlnXchr)
// 	if err1 != nil {
// 		return shim.Error(errMessage("BCCE0004", err1))
// 	}

// 	err = stub.PutState(latestGlnXchr.DocType, lateGlnXchrJSONByte)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0010", err))
// 	}
// 	return shim.Success(nil)
// }

// This Function Performs Query. to get localXchRate  called by all
func (t *fxRateCC) getXchRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type err, invalid json
		return shim.Error(errMessage("BCCE0003", err))
	}
	// Empty Value Check
	if len(checkBlank(qArgs.LocalGlnXchrInfUnqno)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find LOCAL_GLN_XCHR_INF_UNQNO in JSON"))
	}

	queryString := fmt.Sprintf(`{"selector": {"LOCAL_GLN_XCHR_INF_UNQNO": "%s"}}`, qArgs.LocalGlnXchrInfUnqno)
	//Default Size 100
	var pgs int32
	if qArgs.PageSize > 0 {
		pgs = qArgs.PageSize
	} else {
		pgs = pageSize
	}
	// Query
	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// func (t *fxRateCC) getGlnXchRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	var qArgs queryArgs
// 	// JSON Decoding
// 	err := json.Unmarshal([]byte(args[0]), &qArgs)
// 	if err != nil {
// 		// err case: type err, invalid json
// 		return shim.Error(errMessage("BCCE0003", err))
// 	}
// 	// Empty Value Check
// 	if len(checkBlank(qArgs.LocalGlnXchrInfUnqno)) == 0 {
// 		return shim.Error(errMessage("BCCE0005", "Couldn't find XCHR_UNQNO in JSON"))
// 	}

// 	var queryString string
// 	if len(checkBlank(qArgs.LocalGlnXchrInfUnqno)) > 0 {
// 		queryString = fmt.Sprintf(`{"selector": {"LOCAL_GLN_XCHR_INF_UNQNO": "%s"}}`, qArgs.LocalGlnXchrInfUnqno)
// 	} else if checkAtoi(qArgs.PbldDtm) {
// 		return shim.Error(errMessage("BCCE0007", "You must fill out the string number on PbldDtm"))
// 	} else {
// 		queryString = fmt.Sprintf(`{"selector": {"$and":[{"PBLD_DTM": "%s"},{"PBLD_TN": %d}]}}`, qArgs.PbldDtm, qArgs.PbldTn)
// 	}

// 	//Default Size 100
// 	var pgs int32
// 	if qArgs.PageSize > 0 {
// 		pgs = qArgs.PageSize
// 	} else {
// 		pgs = pageSize
// 	}
// 	// Query
// 	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0008", err))
// 	}
// 	logger.Info("Query Success")

// 	return shim.Success(queryResults)
// }

// func (t *fxRateCC) getLatestLocalXchr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	var qArgs queryArgs
// 	err := json.Unmarshal([]byte(args[0]), &qArgs)
// 	if err != nil {
// 		// err case: type err, invalid json
// 		return shim.Error(errMessage("BCCE0003", err))
// 	}
// 	var queryString string
// 	if len(checkBlank(qArgs.LocalGlnCd)) == 0 {
// 		queryString = fmt.Sprintf(`{"selector": {"DOC_TYPE":"latestLocalXchr"}}`)
// 	} else {
// 		queryString = fmt.Sprintf(`{"selector": {"DOC_TYPE":"latestLocalXchr","LC_GLN_CD":"%s"}}`, qArgs.LocalGlnCd)
// 	}
// 	//Default Size 100
// 	var pgs int32
// 	if qArgs.PageSize > 0 {
// 		pgs = qArgs.PageSize
// 	} else {
// 		pgs = pageSize
// 	}
// 	// Query
// 	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)
// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0008", err))
// 	}
// 	logger.Info("Query Success")

// 	return shim.Success(queryResults)
// }

func (t *fxRateCC) getLatestGlnXchr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type err, invalid json
		return shim.Error(errMessage("BCCE0003", err))
	}
	var queryString string
	queryString = fmt.Sprintf(`{"selector": {"DOC_TYPE":"latestGlnXchr"}}`)

	//Default Size 100
	var pgs int32
	if qArgs.PageSize > 0 {
		pgs = qArgs.PageSize
	} else {
		pgs = pageSize
	}
	// Query
	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs insertions. Called by International GLN
func (t *fxRateCC) putXchRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Emtpy Argument Check
	// logger.Info(args[0])
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	//check identity
	// err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	// if err != nil {
	// 	return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	// }
	txID := stub.GetTxID()
	var validData [][]byte
	var keyList []string

	for k := 0; k < len(args); k++ {
		var lcXch exchangeRate
		// Json Decoding
		err := json.Unmarshal([]byte(args[k]), &lcXch)
		if err != nil {
			return shim.Error(errMessage("BCCE0003", err))
		}

		lcXch.Dtm = lcXch.XchrPbldDt + lcXch.XchrPbldHr

		if len(checkBlank(lcXch.Dtm)) != 14 || checkAtoi(lcXch.Dtm) {
			return shim.Error(errMessage("BCCE0007", "Check the string number XCHR_PBLD_DT and XCHR_PBLD_HR"))
		}

		//TX ID
		// Empty Value Check
		if len(checkBlank(lcXch.LocalGlnXchrInfUnqno)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Check LOCAL_GLN_XCHR_INF_UNQNO in JSON"))
		}
		lcXch.TxID = txID

		// Json Encoding
		lcXchJSONBytes, err := json.Marshal(lcXch)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		//Event JSON
		keyList = append(keyList, lcXch.LocalGlnXchrInfUnqno)
		validData = append(validData, lcXchJSONBytes)
	}
	// Duplicate Value Check in couchDB
	mulQuery := multiQueryMaker("LOCAL_GLN_XCHR_INF_UNQNO", keyList)
	queryString := fmt.Sprintf(`{"selector":{%s}, "fields":[%s]}`, mulQuery, `"LOCAL_GLN_XCHR_INF_UNQNO","TX_ID"`)

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

	for i := 0; i < len(validData); i++ {
		// Write Ledger Local GLN Exchange rate Info
		err := stub.PutState(keyList[i], validData[i])

		if err != nil {
			return shim.Error(errMessage("BCCE0010", err))
		}
		logger.Info("Insert Complete")

	}
	return shim.Success(nil)
}

// This Function Performs insertions. Called by International GLN
// func (t *fxRateCC) deployXchRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	// Emtpy Argument Check
// 	if len(args) == 0 {
// 		return shim.Error(errMessage("BCCE0007", "Args is empty"))
// 	}

// 	var glnXchr glnExchangeRate
// 	var latestGlnXchr latestGlnXCHRInfo

//check identity
// err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
// if err != nil {
// 	return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
// }

// 	for i := 0; i < len(args); i++ {
// 		// JSON Decoding
// 		err := json.Unmarshal([]byte(args[i]), &glnXchr)
// 		if err != nil {
// 			// err case: type error, invalid json
// 			return shim.Error(errMessage("BCCE0003", err))
// 		}

// 		queryString := fmt.Sprintf(`{"selector": {"XCHR_UNQCD":"%s"}}`, glnXchr.XchrUnqno)
// 		exs, err := isExist(stub, queryString)
// 		if exs {
// 			return shim.Error(errMessage("BCCE0006", fmt.Sprintf("Data %s", args[i])))
// 		}

// 		// Local Exchange Rate JSON encoding
// 		exRateJSONBytes, err := json.Marshal(glnXchr)
// 		if err != nil {
// 			return shim.Error(errMessage("BCCE0004", err))
// 		}
// 		// Write Ledger Local GLN Exchange rate Info
// 		err = stub.PutState(glnXchr.XchrUnqno, exRateJSONBytes)
// 		if err != nil {
// 			return shim.Error(errMessage("BCCE0010", err))
// 		}
// 		logger.Info("Insert Complete")

// 		queryString = fmt.Sprintf(`{"selector": {"$and":[{"PBLD_DTM":{"$gte": "%s"}},{"PBLD_TN":{"$gte": %d}}]}}`, glnXchr.PbldDtm, glnXchr.PbldTn)
// 		exs, err = isExist(stub, queryString)
// 		if !exs {
// 			latestGlnXchr.DocType = "latestGlnXchr"
// 			latestGlnXchr.glnExchangeRate = glnXchr
// 			lateGlnXchrJSONByte, err1 := json.Marshal(latestGlnXchr)
// 			if err1 != nil {
// 				return shim.Error(errMessage("BCCE0004", err1))
// 			}

// 			err = stub.PutState(latestGlnXchr.DocType, lateGlnXchrJSONByte)
// 			if err != nil {
// 				return shim.Error(errMessage("BCCE0010", err))
// 			}
// 		}
// 		if err != nil {
// 			return shim.Error(errMessage("BCCE0008", err))
// 		}
// 	}
// 	return shim.Success(nil)
// }

// This Function Performs Query. to get localXchRate  called by all
func (t *fxRateCC) getXchRateHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type err, invalid json
		return shim.Error(errMessage("BCCE0003", err))
	}
	// Empty Value Check
	if len(checkBlank(qArgs.LocalGlnXchrInfUnqno)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find XCHR_INF_UNQNO in JSON"))
	}

	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf(`{"selector": {"$and": [{"LOCAL_GLN_CD":"%s"},{"DTM":{"$gte": "%s"}}, {"DTM":{"$lte": "%s"}}]}}`, qArgs.LocalGlnCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	var pgs int32
	if qArgs.PageSize > 0 {
		pgs = qArgs.PageSize
	} else {
		pgs = pageSize
	}

	// Query
	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)

	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// func (t *fxRateCC) getGlnXchRateHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	var qArgs queryArgs
// 	// JSON Decoding
// 	err := json.Unmarshal([]byte(args[0]), &qArgs)
// 	if err != nil {
// 		// err case: type err, invalid json
// 		return shim.Error(errMessage("BCCE0003", err))
// 	}
// 	// Empty Value Check
// 	if len(checkBlank(qArgs.XchrUnqno)) == 0 {
// 		return shim.Error(errMessage("BCCE0005", "Couldn't find XCHR_UNQNO in JSON"))
// 	}

// 	// Valid Check Time String
// 	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
// 		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
// 	}

// 	//Default Size 100
// 	var pgs int32
// 	if qArgs.PageSize > 0 {
// 		pgs = qArgs.PageSize
// 	} else {
// 		pgs = pageSize
// 	}

// 	queryString := fmt.Sprintf(`{"selector": {"$and": [{"RNL_DTM":{"$gte": "%s"}}, {"RNL_DTM":{"$lte": "%s"}}]}}`, qArgs.ReqStartTime, qArgs.ReqEndTime)

// 	// Query
// 	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, pgs, qArgs.BookMark)

// 	if err != nil {
// 		return shim.Error(errMessage("BCCE0008", err))
// 	}
// 	logger.Info("Query Success")

// 	return shim.Success(queryResults)
// }

func (t *fxRateCC) deleteState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("del err")
	}
	return shim.Success(nil)
}
