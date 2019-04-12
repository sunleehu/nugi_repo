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

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("txDataChaincode")
var pageSize int32 = 100

type txDataCC struct {
}

func main() {
	err := shim.Start(new(txDataCC))
	if err != nil {
		fmt.Printf("Error starting txData chaincode: %s", err)
	}
}

func (t *txDataCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *txDataCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)
	fmt.Println("args:", args)
	// Handle different functions
	if function == "puttxdata" {
		return t.putTxData(stub, args)
	} else if function == "gettxdata" {
		return t.getTxData(stub, args)
	} else if function == "gettxdatahistory" {
		return t.getTxDataHistory(stub, args)
	} else if function == "setLogLevel" {
		return setLogLevel(args[0])
	}
	//else if function == "updateTxLog" {
	// 	return t.updateTxLog(stub, args)
	// }
	errM := errMessage("BCCE0001", "Received unknown function invocation "+function)
	return shim.Error(errM)
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *txDataCC) putTxData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	privData, err := stub.GetTransient()
	if err != nil {
		return shim.Error(errMessage("BCCE0005", "GET transient Data Error"))
	}
	// Check Identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	fmt.Println(string(privData["args"]))

	var txdata []transaction

	err = json.Unmarshal(privData["args"], &txdata)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	if len(txdata) < 1 {
		return shim.Error(errMessage("BCCE0007", "Args are empty"))
	}
	//event payload data map,
	evtMap := make(map[string][]string)
	var pyld hEvt
	evtCheck := false
	txID := stub.GetTxID()
	keyMap := make(map[string]string)
	var keyList []string
	var duplList []string

	//validation loop
	for k := 0; k < len(txdata); k++ {
		// Empty Value Check
		if len(checkBlank(txdata[k].GlnTxNo)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Check GLN_TX_NO in JSON"))
		}
		hash := sha256Hash(txdata[k].GlnTxNo)
		keyMap[hash] = txdata[k].GlnTxNo
		keyList = append(keyList, hash)
	}

	// Duplicate Value Check in couchDB
	// mulQuery := multiQueryMaker("GLN_TX_HASH", keyList)
	// queryString := fmt.Sprintf(`{"selector":%s, "fields":[%s]}`, mulQuery, `"GLN_TX_HASH","TX_ID"`)
	// fmt.Println(queryString)

	// exs, res, err := isExist(stub, queryString)
	// if err != nil {
	// 	return shim.Error(errMessage("BCCE0008", err))
	// }
	// if exs {
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
		pData := new(pubData)
		hash := sha256Hash(txdata[i].GlnTxNo)
		txdata[i].GlnTxHash = hash
		pData.BcTxID = txID
		txdata[i].TxID = txID

		// Json Encoding
		txlogJSONBytes, err := json.Marshal(txdata[i])
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}
		// Duplicate Value Check in couchDB
		if len(checkBlank(txdata[i].UtcTxDtm)) != 14 || checkAtoi(txdata[i].UtcTxDtm) {
			return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDDhhmmss"`))
		}

		pData.GlnTxHash = hash
		pData.From = txdata[i].SndrLocalGlnCd
		pData.To = txdata[i].RcvrLocalGlnCd
		pData.Date = txdata[i].UtcTxDtm

		pdd, err := json.Marshal(pData)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// Write couchDB
		err = stub.PutState(hash, pdd)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}

		// Due to collection name error
		// if len(checkBlank(txdata[i].SndrLocalGlnCd)) != 6 || len(checkBlank(txdata[i].RcvrLocalGlnCd)) != 6 {
		// 	return shim.Error(errMessage("BCCE0005", "Check GLN_TX_NO in JSON"))
		// }

		// collection name sorting
		colName := collectionMaker(txdata[i].SndrLocalGlnCd, txdata[i].RcvrLocalGlnCd)

		// Write couchDB
		err = stub.PutPrivateData(colName, txdata[i].GlnTxNo, txlogJSONBytes)
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
		logger.Debug("TRANSACTION_DATA_SAVED", string(dat))
		// EVENT!!!
		stub.SetEvent("TRANSACTION_DATA_SAVED", dat)
	}

	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *txDataCC) getTxData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	if len(args) < 1 {
		return shim.Error(errMessage("BCCE0007", "empty args"))
	}
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

	} else {
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}

	}

	var hash string
	// Empty Value Check
	if len(checkBlank(qArgs.GlnTxNo)) > 0 {
		hash = sha256Hash(qArgs.GlnTxNo)
	} else if len(checkBlank(qArgs.GlnTxHash)) == 64 {
		hash = qArgs.GlnTxHash
	} else {
		return t.getTxDataHistory(stub, args)
	}

	pubQueryRes, err := stub.GetState(hash)

	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	if len(pubQueryRes) < 2 {
		resp := queryResponseStructMaker(nil, "", 0)
		return shim.Success(resp)
	}

	var pData pubData

	err = json.Unmarshal(pubQueryRes, &pData)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	if !(pData.From == qArgs.LcGlnUnqCd || pData.To == qArgs.LcGlnUnqCd) {
		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	}
	colName := collectionMaker(pData.From, pData.To)

	queryResult, err := stub.GetPrivateData(colName, qArgs.GlnTxNo)
	var resList [][]byte
	resList = append(resList, queryResult)
	result := queryResponseStructMaker(resList, "", 1)

	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(result)
}

// This Function Performs Periodic Query. called by International GLN
func (t *txDataCC) getTxDataHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	divcd := ""
	// var colList []string

	if len(args) < 1 {
		return shim.Error(errMessage("BCCE0007", "empty args"))
	}
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

	if qArgs.DivCd == "02" {
		divcd = "TO"
	} else if qArgs.DivCd == "01" {
		divcd = "FROM"
	} else {
		return shim.Error(errMessage("BCCE0005", "You must fill out DIV_Cd"))
	}

	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDDhhmmss"`))
	}
	if len(checkBlank(qArgs.ReqStartTime)) != 14 || len(checkBlank(qArgs.ReqEndTime)) != 14 {
		return shim.Error(errMessage("BCCE0007", `You should fill out date data "YYYYMMDDhhmmss"`))
	}

	//Default Size 100
	var pgs int32
	if qArgs.PageSize > 0 {
		pgs = qArgs.PageSize
	} else if qArgs.PageSize > pageSize {
		pgs = pageSize
	} else {
		pgs = pageSize
	}

	var pData []pubData
	// Query
	queryString := fmt.Sprintf(`{"selector": {"$and":[{"%s": "%s"},{"DATE":{"$gte": "%s"}},{"DATE":{"$lte": "%s"}}]}}`, divcd, qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	queryResults, bookmark, recordCnt, err := getResultForPublicData(stub, queryString, qArgs.BookMark, pgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	err = json.Unmarshal(queryResults, &pData)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	kvList := make(map[string][]string)
	//for sorting collection for query
	for i := 0; i < len(pData); i++ {
		colName := collectionMaker(pData[i].From, pData[i].To)
		kvList[colName] = append(kvList[colName], pData[i].GlnTxHash)
	}
	colNameList := reflect.ValueOf(kvList).MapKeys()
	var respList [][]byte
	for j := 0; j < len(colNameList); j++ {
		colName := colNameList[j].String()
		getPrivSelector := multiSelector("GLN_TX_HASH", kvList[colName])
		qResp, err := getPrivQueryResultForQueryString(stub, colName, getPrivSelector)
		if err != nil {
			return shim.Error(errMessage("BCCE0008", err))
		}
		respList = append(respList, qResp)
	}
	privResp := queryResponseStructMaker(respList, bookmark, recordCnt)
	logger.Info("Query Success")
	return shim.Success(privResp)
}

func (t *txDataCC) updateTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
