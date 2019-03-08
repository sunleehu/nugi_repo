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
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("txLogChaincode")

type txLogChaincode struct {
}

func main() {
	err := shim.Start(new(txLogChaincode))
	if err != nil {
		fmt.Printf("Error starting txLog chaincode: %s", err)
	}
}

func (t *txLogChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *txLogChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)

	// Handle different functions
	if function == "insert_tx_log" {
		return t.insertTxLog(stub, args)
	} else if function == "select_tx_log" {
		return t.selectTxLog(stub, args)
	} else if function == "select_period_tx_log" {
		return t.selectPeriodTxLog(stub, args)
	} else if function == "select_receiver_tx_log" {
		return t.selectReceiverTxLog(stub, args)
	} else if function == "select_sender_tx_log" {
		return t.selectSenderTxLog(stub, args)
	} else if function == "select_period_receiver_tx_log" {
		return t.selectPeriodReceiverTxLog(stub, args)
	} else if function == "select_period_sender_tx_log" {
		return t.selectPeriodSenderTxLog(stub, args)
	} else if function == "select_user_tx_log" {
		return t.selectUserTxLog(stub, args)
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
func (t *txLogChaincode) insertTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Empty Argument Check
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	// Identity Check
	err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Insert Loop
	for i := 0; i < len(args); i++ {
		var tx transaction

		// Json Decoding
		err := json.Unmarshal([]byte(args[i]), &tx)
		if err != nil {
			return shim.Error(errMessage("BCCE0003", err))
		}

		// Empty Value Check
		if len(checkBlank(tx.GlnDeNo)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
		}

		// Duplicate Value Check in couchDB
		queryString := fmt.Sprintf("{\"selector\":{\"GLN_DE_NO\": \"%s\", \"Seq\": %d},\"fields\":[\"GLN_DE_NO\"]}", tx.GlnDeNo, tx.Seq)
		exs, err := isExist(stub, queryString)
		if exs {
			return shim.Error(errMessage("BCCE0006", fmt.Sprintf("Data %s", args[i])))
		}

		// Json Encoding
		txlogJSONBytes, err := json.Marshal(tx)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// Write couchDB
		err = stub.PutState(tx.GlnDeNo+strconv.FormatUint(tx.Seq, 10), txlogJSONBytes)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}
	}

	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *txLogChaincode) selectTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// Json Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Empty Value Check
	if len(checkBlank(qArgs.GlnDeNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
	}

	// Query
	queryString := fmt.Sprintf("{\"selector\": {\"GLN_DE_NO\": \"%s\"}}", qArgs.GlnDeNo)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by International GLN
func (t *txLogChaincode) selectPeriodTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type error, invalid json(empty args)
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error("You must fill out the string number ReqStratTime and ReqEndTime")
	}

	// Query
	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"SNDR_LC_GLN_UNQ_CD\": \"%s\"},{\"GLN_DE_DTM\":{\"$gte\": \"%s\"}},{\"GLN_DE_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *txLogChaincode) selectReceiverTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	}

	// Empty Value Check
	if len(checkBlank(qArgs.GlnDeNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
	}

	// Query
	queryString := fmt.Sprintf("{\"selector\": { \"$and\" : [{\"GLN_DE_NO\": \"%s\"}, {\"RCVR_LC_GLN_UNQ_CD\" : \"%s\"}]}}", qArgs.GlnDeNo, qArgs.LcGlnUnqCd)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *txLogChaincode) selectPeriodReceiverTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArg queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArg)
	if err != nil {

		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArg.LcGlnUnqCd)
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	}

	// Valid Check Time String
	if checkAtoi(qArg.ReqStartTime) || checkAtoi(qArg.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf("{ \"selector\": {\"$and\": [{\"GLN_DE_DTM\":{\"$gte\": \"%s\"}}, {\"GLN_DE_DTM\":{\"$lte\": \"%s\"}}, { \"RCVR_LC_GLN_UNQ_CD\":  \"%s\" }]}}", qArg.ReqStartTime, qArg.ReqEndTime, qArg.LcGlnUnqCd)
	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *txLogChaincode) selectSenderTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {

		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	}

	// Empty Value Check
	if len(checkBlank(qArgs.GlnDeNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
	}

	// Query
	queryString := fmt.Sprintf("{\"selector\": { \"$and\" : [{\"GLN_DE_NO\": \"%s\"}, {\"SNDR_LC_GLN_UNQ_CD\" : \"%s\"}]}}", qArgs.GlnDeNo, qArgs.LcGlnUnqCd)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *txLogChaincode) selectPeriodSenderTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArg queryArgs
	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArg)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArg.LcGlnUnqCd)
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	}

	// Valid Check Time String
	if checkAtoi(qArg.ReqStartTime) || checkAtoi(qArg.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf("{ \"selector\": {\"$and\": [{\"GLN_DE_DTM\":{\"$gte\": \"%s\"}}, {\"GLN_DE_DTM\":{\"$lte\": \"%s\"}}, { \"SNDR_LC_GLN_UNQ_CD\": \"%s\"}]}}", qArg.ReqStartTime, qArg.ReqEndTime, qArg.LcGlnUnqCd)
	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *txLogChaincode) selectUserTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {

		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", qArgs.LcGlnUnqCd)
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
	}

	// Empty Value Check
	if len(checkBlank(qArgs.GlnMbrUnqk)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_MBR_UNQ_K in JSON"))
	}

	queryString := fmt.Sprintf("{\"selector\": { \"$and\" : [{\"GlnMbrUnqk\": \"%s\"}, {\"SNDR_LC_GLN_UNQ_CD\" : \"%s\"}]}}", qArgs.GlnMbrUnqk, qArgs.LcGlnUnqCd)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

func (t *txLogChaincode) updateTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
