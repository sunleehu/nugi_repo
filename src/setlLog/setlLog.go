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

type setlLogChaincode struct {
}

var logger = shim.NewLogger("setlLogChaincode")

func main() {
	err := shim.Start(new(setlLogChaincode))
	if err != nil {
		fmt.Printf("Error starting setlLog chaincode: %s", err)
	}
}

func (t *setlLogChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *setlLogChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)

	//  Handle different functions
	if function == "insert_setl_log" {
		return t.insertSetlLog(stub, args)
	} else if function == "select_setl_log" {
		return t.selectSetlLog(stub, args)
	} else if function == "select_period_setl_log" {
		return t.selectPeriodSetlLog(stub, args)
	} else if function == "select_receiver_setl_log" {
		return t.selectReceiverSetlLog(stub, args)
	} else if function == "select_sender_setl_log" {
		return t.selectSenderSetlLog(stub, args)
	} else if function == "select_period_receiver_setl_log" {
		return t.selectPeriodReceiverSetlLog(stub, args)
	} else if function == "select_period_sender_setl_log" {
		return t.selectPeriodSenderSetlLog(stub, args)
	} else if function == "update_setl_log" {
		return t.updateSetlLog(stub, args)
	} else if function == "setLogLevel" {
		return setLogLevel(args[0])
	}

	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *setlLogChaincode) insertSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Empty Argument Check
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	// Identity Check
	err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	for i := 0; i < len(args); i++ {
		var setl settlmentData

		// Json Decoding
		err := json.Unmarshal([]byte(args[i]), &setl)
		if err != nil {
			return shim.Error(errMessage("BCCE0003", err))
		}

		// Empty Value Check
		if len(checkBlank(setl.GlnDeNo)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
		}

		// Duplicate Value Check in couchDB
		queryString := fmt.Sprintf("{\"selector\":{\"GLN_DE_NO\": \"%s\"},\"fields\":[\"GLN_DE_NO\"]}", setl.GlnDeNo)
		exs, err := isExist(stub, queryString)
		if exs {
			return shim.Error(errMessage("BCCE0006", fmt.Sprintf("Data %s", args[i])))
		}

		// Json Encoding
		setlLogJSONBytes, err := json.Marshal(setl)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// Write couchDB
		err = stub.PutState(setl.GlnDeNo, setlLogJSONBytes)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}
	}
	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *setlLogChaincode) selectSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArg queryArgs
	// Json Decoding
	err := json.Unmarshal([]byte(args[0]), &qArg)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		// err case: type error, invalid json(empty args)
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Empty Value Check
	if len(checkBlank(qArg.GlnDeNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
	}

	// Query
	queryString := fmt.Sprintf("{\"selector\": {\"GLN_DE_NO\": \"%s\"}}", qArg.GlnDeNo)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by International GLN
func (t *setlLogChaincode) selectPeriodSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArg queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArg)
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
	if checkAtoi(qArg.ReqStartTime) || checkAtoi(qArg.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	// Query
	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"GLN_DE_DTM\":{\"$gte\": \"%s\"}},{\"GLN_DE_DTM\":{\"$lte\": \"%s\"}}, {\"RCVR_LC_GLN_UNQ_CD\": \"%s\"}]}}", qArg.ReqStartTime, qArg.ReqEndTime, qArg.LcGlnUnqCd)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *setlLogChaincode) selectReceiverSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	// Empty Value Check
	if len(checkBlank(qArg.GlnDeNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
	}

	// Query
	queryString := fmt.Sprintf("{ \"selector\": {\"$and\": [{\"GLN_DE_NO\": \"%s\" }, { \"RCVR_LC_GLN_UNQ_CD\": \"%s\" }]}, \"fields\": [%s, %s]}", qArg.GlnDeNo, qArg.LcGlnUnqCd, commonField, rcvrField)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *setlLogChaincode) selectPeriodReceiverSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	queryString := fmt.Sprintf("{ \"selector\": {\"$and\": [{\"GLN_DE_DTM\":{\"$gte\": \"%s\"}}, {\"GLN_DE_DTM\":{\"$lte\": \"%s\"}}, { \"RCVR_LC_GLN_UNQ_CD\":  \"%s\"}]}, \"fields\": [%s, %s]}", qArg.ReqStartTime, qArg.ReqEndTime, qArg.LcGlnUnqCd, commonField, rcvrField)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *setlLogChaincode) selectSenderSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	// Empty Value Check
	if len(checkBlank(qArg.GlnDeNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_DE_NO in JSON"))
	}

	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"GLN_DE_NO\": \"%s\"}, {\"SNDR_LC_GLN_UNQ_CD\":\"%s\"}]},\"fields\":[%s, %s]}", qArg.GlnDeNo, qArg.LcGlnUnqCd, commonField, sndrField)

	// Query
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResult)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *setlLogChaincode) selectPeriodSenderSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	queryString := fmt.Sprintf("{\"selector\": {\"$and\": [{\"GLN_DE_DTM\":{\"$gte\": \"%s\"}}, {\"GLN_DE_DTM\":{\"$lte\": \"%s\"}}, { \"SNDR_LC_GLN_UNQ_CD\": \"%s\" }]}, \"fields\": [%s, %s]}", qArg.ReqStartTime, qArg.ReqEndTime, qArg.LcGlnUnqCd, commonField, sndrField)
	// Query
	queryResult, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResult)
}

func (t *setlLogChaincode) updateSetlLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}
