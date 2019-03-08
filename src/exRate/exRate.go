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
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type exRateChaincode struct {
}

// Logger
var logger = shim.NewLogger("exRateChaincode")

func main() {
	err := shim.Start(new(exRateChaincode))
	if err != nil {
		fmt.Printf("Error starting exRate chaincode: %s", err)
	}
}

func (t *exRateChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

func (t *exRateChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	// fmt.Println("invoke is running " + function)
	logger.Info("Invoke is running", function)

	// Handle different functions
	if function == "insert_exchange_rate" {
		return t.insertExRate(stub, args)
	} else if function == "select_exchange_rate_log" {
		return t.selectExRateLog(stub, args)
	} else if function == "select_period_exchange_rate_log" {
		return t.selectPeriodExRateLog(stub, args)
	} else if function == "select_period_nat_exchage_rate_log" {
		return t.selectPeriodNatExRateLog(stub, args)
	} else if function == "select_latest_exchange_rate_log" {
		return t.selectLatestExRateLog(stub, args)
	} else if function == "setLogLevel" {
		return setLogLevel(args[0])
	}

	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions. Called by International GLN
func (t *exRateChaincode) insertExRate(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Emtpy Argument Check
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	// Check Identity
	err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Insert loop
	for i := 0; i < len(args); i++ {
		var ex exchangeRate

		// JSON Decoding
		err := json.Unmarshal([]byte(args[i]), &ex)
		if err != nil {
			// err case: type error, invalid json
			return shim.Error(errMessage("BCCE0003", err))
		}
		// Empty Value Check
		if len(checkBlank(ex.GlnFxNo)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_FX_NO in JSON"))
		}
		// Query string for Duplicate check in couchDB
		queryString := fmt.Sprintf("{\"selector\": {\"GLN_FX_NO\": \"%s\"},\"fields\":[\"GLN_FX_NO\"]}", ex.GlnFxNo)
		exs, err := isExist(stub, queryString)
		if exs {
			return shim.Error(errMessage("BCCE0006", fmt.Sprintf("Data %s", args[i])))
		}

		// JSON encoding
		exRateJSONBytes, err := json.Marshal(ex)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}
		// Write Ledger
		err = stub.PutState(ex.GlnFxNo, exRateJSONBytes)
		if err != nil {
			return shim.Error(errMessage("BCCE0010", err))
		}
	}
	logger.Info("Insert Complete")

	return shim.Success(nil)
}

// This Function Performs Query. called by all
func (t *exRateChaincode) selectExRateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type err, invalid json
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Empty Value Check
	if len(checkBlank(qArgs.GlnFxNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find GLN_FX_NO in JSON"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"GLN_FX_NO\": \"%s\"}}", qArgs.GlnFxNo)
	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by all
func (t *exRateChaincode) selectPeriodExRateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: type error, invalid json
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"UP_DTM\":{\"$gte\": \"%s\"}},{\"UP_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.ReqStartTime, qArgs.ReqEndTime)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query with national code. called by all
func (t *exRateChaincode) selectPeriodNatExRateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"NAT_CD\": \"%s\"},{\"UP_DTM\":{\"$gte\": \"%s\"}},{\"UP_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.NatCd, qArgs.ReqStartTime, qArgs.ReqEndTime)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	return shim.Success(queryResults)
}

// This Function Performs Query for get Latest Value. called by all
func (t *exRateChaincode) selectLatestExRateLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		return shim.Error(errMessage("BCCE0003", err))
	}

	// For Safety: Query Time from PbldDtm - 1
	safety, err := strconv.Atoi(qArgs.PbldDtm)
	if err != nil {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number PbldDtm"))
	}
	safety--

	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"NAT_CD\":\"%s\"},{\"PBLD_DTM\":{\"$gte\":\"%s\"}}]},\"sort\":[{\"UP_DTM\":\"desc\"}]}", qArgs.NatCd, strconv.Itoa(safety))

	// Query
	queryResults, err := getQueryResultForLatest(stub, queryString)

	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}
