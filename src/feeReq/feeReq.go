/*
* This Chaincode is GLN Commission fee code
* And it has functions insert, update and query
* International GLN can insert Commission Fee Settlement Data
* Interational and Local GLN can query Data
* Local GLN can update the Data after finishing the settlement. and make update event
* International GLN is possible to Make Event When Data Insertion
**/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type feeReqChaincode struct {
}

var logger = shim.NewLogger("feeReqChaincode")

func main() {
	err := shim.Start(new(feeReqChaincode))
	if err != nil {
		fmt.Printf("Error starting feeReq chaincode: %s", err)
	}
}

func (t *feeReqChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *feeReqChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)
	// Handle different functions
	if function == "insert_fee_req" {
		return t.insertFeeReq(stub, args)
	} else if function == "select_fee_req_log" {
		return t.selectFeeReqLog(stub, args)
	} else if function == "select_period_fee_req_log" {
		return t.selectPeriodFeeReqLog(stub, args)
	} else if function == "select_sender_fee_req_log" {
		return t.selectSenderFeeReqLog(stub, args)
	} else if function == "select_period_sender_fee_req_log" {
		return t.selectPeriodSenderFeeReqLog(stub, args)
	} else if function == "update_sender_fee_pmt_result" {
		return t.updateSenderFeePmtResult(stub, args)
	} else if function == "setLogLevel" {
		return setLogLevel(args[0])
	}

	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *feeReqChaincode) insertFeeReq(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Empty Argument Check
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	// Check Identity
	err := cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Event Payload variables
	var pyld hEvt
	var idata iEvt
	evtCheck := false

	// Insert Loop
	for i := 0; i < len(args); i++ {
		var pr feeReq

		// JSON Decoding
		err := json.Unmarshal([]byte(args[i]), &pr)
		if err != nil {
			return shim.Error(errMessage("BCCE0003", err))
		}

		// Empty Value Check
		if len(checkBlank(pr.AdjMnDsbReqNo)) == 0 {
			return shim.Error(errMessage("BCCE0005", "Couldn't find ADJ_MN_DSB_REQ_NO in JSON"))
		}

		// Duplicate Value Check in couchDB
		queryString := fmt.Sprintf("{\"selector\": {\"ADJ_MN_DSB_REQ_NO\": \"%s\"},\"fields\":[\"ADJ_MN_DSB_REQ_NO\"]}", pr.AdjMnDsbReqNo)
		exs, err := isExist(stub, queryString)
		if exs {
			return shim.Error(errMessage("BCCE0006", fmt.Sprintf("Data %s", args[i])))
		}

		// JSON Encoding
		feeReqJSONBytes, err := json.Marshal(pr)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}
		// Write couchDB
		err = stub.PutState(pr.AdjMnDsbReqNo, feeReqJSONBytes)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}

		// Event Payload
		pyld.Target = append(pyld.Target, pr.LcGlnUnqCd)
		idata.AdjDtm = pr.AdjDtm
		pyld.Data = idata
		evtCheck = true
	}

	if evtCheck {
		pyld.Target = rmvDupVal(pyld.Target)
		dat, e := json.Marshal(pyld)

		if e != nil {
			return shim.Error(errMessage("BCCE0004", e))
		}

		logger.Info("EVENT_FEE_SETTLMENT_START")
		logger.Debug("EVENT_FEE_SETTLMENT_START", string(dat))
		// EVENT!!
		stub.SetEvent("EVENT_FEE_SETTLMENT_START", dat)
	}
	logger.Info("Insert Complete")

	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *feeReqChaincode) selectFeeReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs
	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
		// err case: invalid json(empty args), type error
		return shim.Error(errMessage("BCCE0003", err))
	}

	// Check Identity
	err = cid.AssertAttributeValue(stub, "ACC_ROLE", "INT")
	if err != nil {
		return shim.Error(errMessage("BCCE0002", "This function Only for INT GLN"))
	}

	// Empty Value Check
	if len(checkBlank(qArgs.AdjMnDsbReqNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find ADJ_MN_DSB_REQ_NO in JSON"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"ADJ_MN_DSB_REQ_NO\": \"%s\"}}", qArgs.AdjMnDsbReqNo)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by International GLN
func (t *feeReqChaincode) selectPeriodFeeReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"LC_GLN_UNQ_CD\": \"%s\"},{\"ADJ_DTM\":{\"$gte\": \"%s\"}},{\"ADJ_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *feeReqChaincode) selectSenderFeeReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	if len(checkBlank(qArgs.AdjMnDsbReqNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find ADJ_MN_DSB_REQ_NO in JSON"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"ADJ_MN_DSB_REQ_NO\": \"%s\",\"LC_GLN_UNQ_CD\": \"%s\"}}", qArgs.AdjMnDsbReqNo, qArgs.LcGlnUnqCd)
	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *feeReqChaincode) selectPeriodSenderFeeReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	// Valid Check Time String
	if checkAtoi(qArgs.ReqStartTime) || checkAtoi(qArgs.ReqEndTime) {
		return shim.Error(errMessage("BCCE0007", "You must fill out the string number ReqStratTime and ReqEndTime"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"LC_GLN_UNQ_CD\": \"%s\"},{\"ADJ_DTM\":{\"$gte\": \"%s\"}},{\"ADJ_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)
	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs update. called by Local GLN
func (t *feeReqChaincode) updateSenderFeePmtResult(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Empty Argument Check
	if len(args) == 0 {
		return shim.Error(errMessage("BCCE0007", "Args is empty"))
	}

	// Event Payload Variable
	var pyld hEvt
	var udata uEvt
	evtCheck := false
	// Update Loop
	for i := 0; i < len(args); i++ {
		var uArgs resultArgs

		err := json.Unmarshal([]byte(args[i]), &uArgs)
		if err != nil {
			return shim.Error(errMessage("BCCE0003", err))
		}

		// Check Identity
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", uArgs.LcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}

		queryString := fmt.Sprintf("{\"selector\":{\"LC_GLN_UNQ_CD\":\"%s\", \"ADJ_MN_DSB_REQ_NO\": \"%s\"}}", uArgs.LcGlnUnqCd, uArgs.AdjMnDsbReqNo)

		// Get Data
		resultsIterator, err := stub.GetQueryResult(queryString)
		if err != nil {
			return shim.Error(errMessage("BCCE0008", err))
		}
		defer resultsIterator.Close()
		if !resultsIterator.HasNext() {
			return shim.Error(errMessage("BCCE0010", fmt.Sprintf("Data %s", args[i])))
		}
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			logger.Debug("queryResponse:", queryResponse)
			if err != nil {
				return shim.Error(errMessage("BCCE0008", err))
			}
			var feeRes feeReq

			json.Unmarshal(queryResponse.Value, &feeRes)

			// Update Value
			if feeRes.LcGlnUnqCd == uArgs.LcGlnUnqCd {
				feeRes.DsbCompYn = uArgs.DsbCompYn
			}
			jtx, err := json.Marshal(feeRes)
			if err != nil {
				return shim.Error(errMessage("BCCE0004", err))
			}

			// Update CouchDB
			err = stub.PutState(queryResponse.Key, jtx)
			if err != nil {
				return shim.Error(errMessage("BCCE0010", err))
			}

			// Event Payload
			udata.LcGlnUnqCd = uArgs.LcGlnUnqCd
			udata.AdjMnDsbReqNo = append(udata.AdjMnDsbReqNo, uArgs.AdjMnDsbReqNo)
			evtCheck = true
		}
		defer resultsIterator.Close()

	}

	//EVENT Block
	if evtCheck {
		// EVENT Payload JSON Encoding
		udata.AdjMnDsbReqNo = rmvDupVal(udata.AdjMnDsbReqNo)
		pyld.Data = udata
		dat, e := json.Marshal(pyld)
		if e != nil {
			return shim.Error(errMessage("BCCE0004", e))
		}
		logger.Info("EVENT_FEE_SETTLEMENT_COMPLETE")
		logger.Debug("EVENT_FEE_SETTLEMENT_COMPLETE", string(dat))
		// EVENT!!!!
		stub.SetEvent("EVENT_FEE_SETTLEMENT_COMPLETE", dat)
	}
	logger.Info("Update Complete")
	return shim.Success(nil)
}
