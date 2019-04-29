/*
* This Chaincode is GLN Tranasction Log code
* And it has functions insert and query,
* International GLN can insert Transaction Log.
* Interational and Local GLN can query Data
**/
package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("TXDATA")
var DEFAULT_PAGE_SIZE int32 = 100

type txDataCC struct {
}

func main() {
	err := shim.Start(new(txDataCC))
	if err != nil {
		logger.Error("Error starting txdata chaincode : %s", err)
	}
}

func (t *txDataCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *txDataCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)
	logger.Info("Args: ", args)

	// Handle different functions
	if function == "puttxdata" {
		return t.putTxData(stub, args)
	} else if function == "gettxdata" {
		return t.getTxData(stub, args)
	} else if function == "gettxdatahistory" {
		return t.getTxDataHistory(stub, args)
	}
	// } else if function == "setLogLevel" {
	// 	return setLogLevel(args[0])
	// }
	//else if function == "updateTxLog" {
	// 	return t.updateTxLog(stub, args)
	// }

	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *txDataCC) putTxData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check arguments
	privData, err := stub.GetTransient()
	if err != nil {
		return shim.Error(errMessage("BCCE0005", "GET transient Data Error"))
	}
	logger.Debug("Transient Args : ", string(privData["args"]))
	var txdata []transaction
	err = json.Unmarshal(privData["args"], &txdata)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	if len(txdata) < 1 {
		return shim.Error(errMessage("BCCE0007", "Args are empty"))
	}
	logger.Info("Put Data Count : ", len(txdata))

	// Check Identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	//event payload data map,
	evtMap := make(map[string][]string)
	var pyld hEvt
	evtCheck := false
	txID := stub.GetTxID()
	// keyMap := make(map[string]string)
	// var keyList []string

	//validation loop
	// for k := 0; k < len(txdata); k++ {
	// 	// Empty Value Check
	// 	if len(checkBlank(txdata[k].GlnTxNo)) == 0 {
	// 		return shim.Error(errMessage("BCCE0005", "Check GLN_TX_NO in JSON"))
	// 	}

	// 	// hash := sha256Hash(txdata[k].GlnTxNo)
	// 	// keyMap[hash] = txdata[k].GlnTxNo
	// 	// keyList = append(keyList, hash)
	// }

	// Duplicate Value Check in couchDB
	// var duplList []string
	// mulQuery := multiQueryMaker("GLN_TX_HASH", keyList)
	// queryString := fmt.Sprintf(`{"selector":%s, "fields":[%s]}`, mulQuery, `"GLN_TX_HASH","TX_ID"`)
	// fmt.Println(queryString)
	// exs, res, err := isExist(stub, queryString)
	// if err != nil {
	// 	return shim.Error(errMessage("BCCE0008", err))
	// }
	// if exs {
	// 	if err != nil {
	// 		return shim.Error(errMessage("BCCE0008", err))
	// 	}
	// 	var qResp []pubData
	// 	json.Unmarshal(res, &qResp)
	// 	for j := 0; j < len(qResp); j++ {
	// 		var respJ respStruct
	// 		respJ.BcTxID = qResp[j].BcTxID
	// 		respJ.GlnTxNo = keyMap[qResp[j].GlnTxHash]
	// 		respJSONBytes, err := json.Marshal(respJ)
	// 		if err != nil {
	// 			return shim.Error(errMessage("BCCE0004", err))
	// 		}
	// 		duplList = append(duplList, string(respJSONBytes))
	// 	}
	// 	return shim.Error(errMessage("BCCE0006", fmt.Sprintf("[%s]", strings.Join(duplList, ","))))
	// }

	// Insert Loop
	for i := 0; i < len(txdata); i++ {

		// Value Check
		if isBlank(txdata[i].GlnTxNo) {
			return shim.Error(errMessage("BCCE0007", "Check GLN_TX_NO in JSON"))
		}
		if len(strings.TrimSpace(txdata[i].UtcTxDtm)) != 14 || checkAtoi(txdata[i].UtcTxDtm) {
			return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDDhhmmss"`))
		}
		if len(strings.TrimSpace(txdata[i].SndrLocalGlnCd)) != 6 || len(strings.TrimSpace(txdata[i].RcvrLocalGlnCd)) != 6 {
			return shim.Error(errMessage("BCCE0007", "Check GLN_TX_NO in JSON"))
		}

		hash := sha256Hash(txdata[i].GlnTxNo)
		txdata[i].GlnTxHash = hash
		txdata[i].TxID = txID

		// Make collection name
		colName := collectionMaker(txdata[i].SndrLocalGlnCd, txdata[i].RcvrLocalGlnCd)

		// Public data
		pData := new(pubData)
		pData.GlnTxHash = hash
		pData.Date = txdata[i].UtcTxDtm
		pData.From = txdata[i].SndrLocalGlnCd
		pData.To = txdata[i].RcvrLocalGlnCd
		pData.BcTxID = txdata[i].TxID
		pdd, err := json.Marshal(pData)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// Write public data
		err = stub.PutState(hash, pdd)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}

		// Write private data
		txlogJSONBytes, err := json.Marshal(txdata[i])
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}
		err = stub.PutPrivateData(colName, hash, txlogJSONBytes)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}

		pyld.Target = append(pyld.Target, txdata[i].RcvrLocalGlnCd, txdata[i].SndrLocalGlnCd)
		evtMap[txdata[i].SndrLocalGlnCd] = append(evtMap[txdata[i].SndrLocalGlnCd], hash)
		evtMap[txdata[i].RcvrLocalGlnCd] = append(evtMap[txdata[i].RcvrLocalGlnCd], hash)
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

		logger.Info("TRANSACTION_DATA_SAVED")
		//logger.Debug("SAVED_DATA : ", string(dat))
		// EVENT!!!
		stub.SetEvent("TRANSACTION_DATA_SAVED", dat)
	}

	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *txDataCC) getTxData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check arguments
	var qArgs queryArgs
	if len(args) < 1 {
		return shim.Error(errMessage("BCCE0007", "empty args"))
	}
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

	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LOCALGLN_CODE does not match"))
		}
	}

	// Empty Value Check
	var hash string
	if isBlank(qArgs.GlnTxNo) {
		return t.getTxDataHistory(stub, args)
	} else {
		hash = sha256Hash(qArgs.GlnTxNo)
	}

	pubQueryRes, err := stub.GetState(hash)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if pubQueryRes == nil {
		resp := queryResponseStructMaker(nil, "", 0)
		return shim.Success(resp)
	}

	var pData pubData
	err = json.Unmarshal(pubQueryRes, &pData)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// if len(checkBlank(qArgs.LcGlnUnqCd)) > 0 {
	// 	if qArgs.LcGlnUnqCd != pData.From && qArgs.LcGlnUnqCd != pData.To {
	// 		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	// 	}
	// }

	colName := collectionMaker(pData.From, pData.To)
	queryResult, err := stub.GetPrivateData(colName, hash)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if queryResult == nil {
		resp := queryResponseStructMaker(nil, "", 0)
		return shim.Success(resp)
	}

	var resList [][]byte
	resList = append(resList, queryResult)
	result := queryResponseStructMaker(resList, "", 1)

	logger.Info("Query Success")
	return shim.Success(result)
}

// This Function Performs Periodic Query. called by International GLN
func (t *txDataCC) getTxDataHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check arguments
	if len(args) < 1 {
		return shim.Error(errMessage("BCCE0007", "empty args"))
	}
	var qArgs queryArgs
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDDhhmmss"`))
	}
	if len(strings.TrimSpace(qArgs.ReqStartTime)) != 14 || len(strings.TrimSpace(qArgs.ReqEndTime)) != 14 {
		return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDDhhmmss"`))
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

	// Valid DIV_CODE
	divcd := ""
	if qArgs.LcGlnUnqCd != "" {
		if qArgs.DivCd == "02" {
			divcd = "TO"
		} else if qArgs.DivCd == "01" {
			divcd = "FROM"
		} else {
			return shim.Error(errMessage("BCCE0005", "You must fill out DIV_CD"))
		}
	}

	//Default Size 100
	pgs := qArgs.PageSize
	if pgs == 0 {
		pgs = DEFAULT_PAGE_SIZE
	}

	// Query
	var pData []pubData
	queryString := ""
	if divcd == "" {
		queryString = fmt.Sprintf(`{"selector": {"$and":[{"DATE":{"$gte": "%s"}},{"DATE":{"$lte": "%s"}}]}, "use_index":["indexDateDoc", "indexDate"]}`, qArgs.ReqStartTime, qArgs.ReqEndTime)
	} else {
		queryString = fmt.Sprintf(`{"selector": {"$and":[{"%s": "%s"},{"DATE":{"$gte": "%s"}},{"DATE":{"$lte": "%s"}}]}, "use_index":["indexDate%sDoc", "indexDate%s"]}`, divcd, qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime, divcd, divcd)
	}
	queryResults, bookmark, recordCnt, err := getResultForPublicData(stub, queryString, qArgs.BookMark, pgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	// To json
	err = json.Unmarshal(queryResults, &pData)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// for sorting collection for query
	kvList := make(map[string][]string)
	for i := 0; i < len(pData); i++ {
		colName := collectionMaker(pData[i].From, pData[i].To)
		kvList[colName] = append(kvList[colName], pData[i].GlnTxHash)
	}
	colNameList := reflect.ValueOf(kvList).MapKeys()
	var respList [][]byte
	for j := 0; j < len(colNameList); j++ {
		colName := colNameList[j].String()
		// getPrivSelector := multiSelector("GLN_TX_HASH", kvList[colName])
		// qResp, err := getPrivQueryResultForQueryString(stub, colName, getPrivSelector)
		qResp, err := getPrivateDataForKeys(stub, colName, kvList[colName])
		if err != nil {
			return shim.Error(errMessage("BCCE0008", err))
		}
		if len(qResp) > 0 {
			respList = append(respList, qResp)
		}
	}
	privResp := queryResponseStructMaker(respList, bookmark, recordCnt)
	logger.Info("Query Success")
	return shim.Success(privResp)
}

// func (t *txDataCC) updateTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	return shim.Success(nil)
// }
