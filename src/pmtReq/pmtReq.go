/*
* This Chaincode is GLN Payment code
* And it has functions insert, update and query
* International GLN can insert Payment Request Data
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

type pmtReqChaincode struct {
}

var logger = shim.NewLogger("pmtReqChaincode")

func main() {
	err := shim.Start(new(pmtReqChaincode))
	if err != nil {
		fmt.Printf("Error starting pmtReq chaincode: %s", err)
	}
}

func (t *pmtReqChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *pmtReqChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Info("Invoke is running", function)

	// Handle different functions
	if function == "insert_pmt_req" {
		return t.insertPmtReq(stub, args)
	} else if function == "select_pmt_req_log" {
		return t.selectPmtReqLog(stub, args)
	} else if function == "select_period_pmt_req_log" {
		return t.selectPeriodPmtReqLog(stub, args)
	} else if function == "select_sender_pmt_req_log" {
		return t.selectSenderPmtReqLog(stub, args)
	} else if function == "select_receiver_pmt_req_log" {
		return t.selectReceiverPmtReqLog(stub, args)
	} else if function == "select_period_sender_pmt_req_log" {
		return t.selectPeriodSenderPmtReqLog(stub, args)
	} else if function == "select_period_receiver_pmt_req_log" {
		return t.selectPeriodReceiverPmtReqLog(stub, args)
	} else if function == "update_receiver_pmt_result" {
		return t.updateReceiverPmtResult(stub, args)
	} else if function == "setLogLevel" {
		return setLogLevel(args[0])
	}

	return shim.Error(errMessage("BCCE0001", "Received unknown function invocation "+function))
}

// This Function Performs insertions and Generating settlement start event. Called by International GLN
func (t *pmtReqChaincode) insertPmtReq(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		var pr pmtReq

		// JSON Decoding with
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
		pmtReqJSONBytes, err := json.Marshal(pr)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}

		// Write couchDB
		err = stub.PutState(pr.AdjMnDsbReqNo, pmtReqJSONBytes)
		if err != nil {
			return shim.Error(errMessage("BCCE0009", err))
		}

		// Event Payload
		pyld.Target = append(pyld.Target, pr.RcvrLcGlnUnqCd, pr.SndrLcGlnUnqCd)
		idata.AdjDtm = pr.AdjDtm
		pyld.Data = idata
		evtCheck = true
	}

	// EVENT BLOCK
	if evtCheck {
		// Event Payload JSON Encoding
		pyld.Target = rmvDupVal(pyld.Target)
		dat, e := json.Marshal(pyld)

		if e != nil {
			return shim.Error(errMessage("BCCE0004", e))
		}

		logger.Info("EVENT_PAYMENT_SETTLEMENT_START")
		logger.Debug("EVENT_PAYMENT_SETTLEMENT_START", string(dat))
		// EVENT!!!
		stub.SetEvent("EVENT_PAYMENT_SETTLEMENT_START", dat)
	}

	logger.Info("Insert Complete")
	return shim.Success(nil)
}

// This Function Performs Query. called by International GLN
func (t *pmtReqChaincode) selectPmtReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
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
func (t *pmtReqChaincode) selectPeriodPmtReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
	err := json.Unmarshal([]byte(args[0]), &qArgs)
	if err != nil {
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

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"SNDR_LC_GLN_UNQ_CD\": \"%s\"},{\"ADJ_DTM\":{\"$gte\": \"%s\"}},{\"ADJ_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *pmtReqChaincode) selectSenderPmtReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var qArgs queryArgs

	// JSON Decoding
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
	if len(checkBlank(qArgs.AdjMnDsbReqNo)) == 0 {
		return shim.Error(errMessage("BCCE0005", "Couldn't find ADJ_MN_DSB_REQ_NO in JSON"))
	}

	queryString := fmt.Sprintf("{\"selector\": {\"ADJ_MN_DSB_REQ_NO\": \"%s\",\"SNDR_LC_GLN_UNQ_CD\": \"%s\"}}", qArgs.AdjMnDsbReqNo, qArgs.LcGlnUnqCd)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs Query. called by Local GLN
func (t *pmtReqChaincode) selectReceiverPmtReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	queryString := fmt.Sprintf("{ \"selector\": {\"$and\": [{\"ADJ_MN_DSB_REQ_NO\": \"%s\" }, { \"RCVR_LC_GLN_UNQ_CD\": \"%s\" }]}}", qArgs.AdjMnDsbReqNo, qArgs.LcGlnUnqCd)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")

	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *pmtReqChaincode) selectPeriodSenderPmtReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"SNDR_LC_GLN_UNQ_CD\": \"%s\"},{\"ADJ_DTM\":{\"$gte\": \"%s\"}},{\"ADJ_DTM\":{\"$lte\": \"%s\"}}]}}", qArgs.LcGlnUnqCd, qArgs.ReqStartTime, qArgs.ReqEndTime)

	// Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}

	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs Periodic Query. called by Local GLN
func (t *pmtReqChaincode) selectPeriodReceiverPmtReqLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	queryString := fmt.Sprintf("{\"selector\": {\"$and\": [{\"ADJ_DTM\":{\"$gte\": \"%s\"}}, {\"ADJ_DTM\":{\"$lte\": \"%s\"}},{ \"RCVR_LC_GLN_UNQ_CD\": \"%s\" }]}}", qArgs.ReqStartTime, qArgs.ReqEndTime, qArgs.LcGlnUnqCd)

	//Query
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(errMessage("BCCE0008", err))
	}
	logger.Info("Query Success")
	return shim.Success(queryResults)
}

// This Function Performs update. called by Local GLN
func (t *pmtReqChaincode) updateReceiverPmtResult(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
		err = cid.AssertAttributeValue(stub, "LCL_UNQ_CD", uArgs.RcvrLcGlnUnqCd)
		if err != nil {
			return shim.Error(errMessage("BCCE0002", "Tx Maker and LclGlnUnqCd does not match"))
		}

		// Get Data
		queryString := fmt.Sprintf("{\"selector\":{\"SNDR_LC_GLN_UNQ_CD\":\"%s\", \"RCVR_LC_GLN_UNQ_CD\":\"%s\", \"ADJ_MN_DSB_REQ_NO\": \"%s\"}}", uArgs.SndrLcGlnUnqCd, uArgs.RcvrLcGlnUnqCd, uArgs.AdjMnDsbReqNo)
		logger.Debug("QueryString:", queryString)

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

			if err != nil {
				return shim.Error(errMessage("BCCE0008", err))
			}

			var pmtRes pmtReq

			json.Unmarshal(queryResponse.Value, &pmtRes)
			logger.Debug("QueryResponse:", queryResponse)
			// Update Value
			if pmtRes.RcvrLcGlnUnqCd == uArgs.RcvrLcGlnUnqCd {
				pmtRes.AdjCompYn = uArgs.AdjCompYn
			}

			jtx, err := json.Marshal(pmtRes)

			if err != nil {
				return shim.Error(errMessage("BCCE0004", err))
			}

			// Update CouchDB
			err = stub.PutState(queryResponse.Key, jtx)
			if err != nil {
				return shim.Error(errMessage("BCCE0010", err))
			}

			// Event Payload
			udata.LcGlnUnqCd = uArgs.RcvrLcGlnUnqCd
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
		dat, err := json.Marshal(pyld)
		if err != nil {
			return shim.Error(errMessage("BCCE0004", err))
		}
		logger.Info("EVENT_PAYMENT_UPDATE_COMPLETE")
		logger.Debug("EVENT_PAYMENT_UPDATE_COMPLETE", string(dat))
		// EVENT!!!!
		stub.SetEvent("EVENT_PAYMENT_UPDATE_COMPLETE", dat)
	}

	logger.Info("Update Complete")
	return shim.Success(nil)
}
