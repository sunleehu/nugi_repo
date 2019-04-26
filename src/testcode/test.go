package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type firstCC struct {
	Version string
}

type transaction struct {
	DeSeq    string `json:"de_seq"` // 거래일련번호
	Balance  int    `json:"balance,omitempty"`
	AddField int    `json:"add_field,omitempty"`
}

type payloads struct {
}

type queryArgs struct {
	DeSeq        string
	ReqStartTime string
	ReqEndTime   string
}

var logger = shim.NewLogger("firstCC")

func main() {
	// logger.Info("When?")
	err := shim.Start(new(firstCC))
	if err != nil {
		fmt.Printf("Error starting First chaincode: %s", err)
	}
}

func (t *firstCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("GLN TEST CHAIN CODE")
	welcomeMsg := map[string]string{
		"message": "Welcome to GLN BLOCK CHAIN",
	}
	wmsg, err := json.Marshal(welcomeMsg)
	if err != nil {
		shim.Error(errMessage("BCCE0004", err))
	}
	err = stub.PutState("message", wmsg)

	if err != nil {
		shim.Error(errMessage("BCCE0009", err))
	}

	return shim.Success(nil)
}

func (t *firstCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	// fmt.Println("invoke is running " + function)
	logger = shim.NewLogger("CC:test TX:" + stub.GetTxID())
	logger.Infof("Invoke is running %s", function)
	// Handle different functions
	if function == "putData" {
		return t.putData(stub, args)
	} else if function == "getCertRole" {
		return t.getCertRole(stub)
	} else if function == "getWelcomeMessage" {
		return t.getWelcomeMessage(stub)
	} else if function == "welcomeEvt" {
		return t.welcomeEvt(stub)
	} else if function == "healthCheck" {
		return t.healthCheck(stub)
	} else if function == "getData" {
		return t.getData(stub, args)
	} else if function == "eventStruct" {
		return t.eventStruct(stub, args)
	}
	// fmt.Println("Could not find func: " + function) //error
	errObj := errMessage("BCCE0001", "Received unknown function invocation "+function)

	return shim.Error(errObj)
}
func (t *firstCC) welcomeEvt(stub shim.ChaincodeStubInterface) pb.Response {
	var evtPayload hEvt
	evtPayload.Target = append(evtPayload.Target)
	evtPayload.Data = fmt.Sprintf("Greeting! You've got Event Message!")

	payloadByte, err := json.Marshal(evtPayload)

	err = stub.SetEvent("WELCOME_EVT", payloadByte)

	if err != nil {
		return shim.Error(errMessage("BCCE0004", err))
	}
	return shim.Success(nil)
}

func (t *firstCC) eventStruct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var evtPayload hEvt
	evtMap := make(map[string][]string)
	evtMap[args[0]] = append(evtMap[args[0]], args[0])
	evtPayload.Target = append(evtPayload.Target, args[0])
	evtPayload.Data = evtMap

	payloadByte, err := json.Marshal(evtPayload)

	err = stub.SetEvent("EVT_STRUCT", payloadByte)

	if err != nil {
		return shim.Error(errMessage("BCCE0004", err))
	}
	return shim.Success(nil)
}

func (t *firstCC) healthCheck(stub shim.ChaincodeStubInterface) pb.Response {

	cert := t.getCertRole(stub)
	logger.Info("Are you %s?", cert.GetPayload())

	return shim.Success(nil)
}

func (t *firstCC) getCertRole(stub shim.ChaincodeStubInterface) pb.Response {
	id, err := cid.New(stub)
	if err != nil {
		logger.Error(err)
		return shim.Error(errMessage("BCCE0011", err))
	}
	role, ok, err := id.GetAttributeValue("LCL_UNQ_CD")

	if err != nil {
		logger.Error(err)
		return shim.Error(errMessage("BCCE0011", err))
	}
	if !ok {
		return shim.Error(errMessage("BCCE0011", "Cert Doesn't have LCL_UNQ_CD"))
	}
	return shim.Success([]byte(role))
}

func (t *firstCC) getWelcomeMessage(stub shim.ChaincodeStubInterface) pb.Response {
	queryString := fmt.Sprintf(`{"selector":{"message":"Welcome to GLN BLOCK CHAIN"}}`)
	result, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(result)
}
func (t *firstCC) deleteState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("del err")
	}
	return shim.Success(nil)
}

func (t *firstCC) getData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	queryString := args[0]
	queryResult, err := getQueryResultForQueryStringWithPagination(stub, queryString, int32(2), args[1])
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}
	logger.Info(queryResult)

	return shim.Success(queryResult)
}

func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{\"BC_RES_DATA\":[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("],")

	return &buffer, nil
}

func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {

	buffer.WriteString("\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}")
	return buffer
}

func (t *firstCC) putData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ret []byte

	for i := 0; i < len(args); i++ {
		var tx transaction
		fmt.Println(args[i])
		err := json.Unmarshal([]byte(args[i]), &tx)
		if err != nil {
			fmt.Println("type Error:", err)
			return shim.Error("Check Your Json")

		}
		fmt.Println(":", tx)
		txlogJSONBytes, err := json.Marshal(tx)

		fmt.Println(tx)

		err = stub.PutState(tx.DeSeq, txlogJSONBytes)
		// err = stub.PutPrivateData("privCollection", tx.DeSeq, txlogJSONBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		ret = txlogJSONBytes
	}

	return shim.Success(ret)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	fmt.Println("ResultsIterator: ", resultsIterator)

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		fmt.Println(queryResponse)

		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))

		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func getQueryResultForQueryStringPriv(stub shim.ChaincodeStubInterface, collection, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetPrivateDataQueryResult(collection, queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	fmt.Println("ResultsIterator: ", resultsIterator)

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		fmt.Println(queryResponse)

		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// buffer.WriteString("\"{\"")
		// buffer.WriteString("{\"Key\":")
		// buffer.WriteString("\"")
		// buffer.WriteString(queryResponse.Key)
		// buffer.WriteString("\"")

		// buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		// buffer.WriteString("\"}\"")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func getQueryResultForArray(stub shim.ChaincodeStubInterface, queryString string, elem []string, receiver string) ([]byte, error) {

	// fmt.Printf("- getQueryResultForArray: queryString: \n%s\n", queryString)
	// id := "gcoinuser1"

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		// buffer.WriteString(string(queryResponse.Value)

		var datmap map[string]interface{}
		var arrmap []map[string]interface{}
		val := queryResponse.Value
		json.Unmarshal([]byte(val), &datmap)
		arr, err := json.Marshal(datmap[elem[0]])
		json.Unmarshal([]byte(arr), &arrmap)
		for i := 0; i < len(arrmap); i++ {
			if arrmap[i][elem[1]] != receiver {
				arrmap = append(arrmap[:i], arrmap[i+1:]...)
				i--
			}
		}
		datmap[elem[0]] = arrmap

		co, err := json.Marshal(datmap)
		buffer.WriteString(string(co))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	defer resultsIterator.Close()

	return buffer.Bytes(), nil
}
func isExist(stub shim.ChaincodeStubInterface, queryString string) (bool, error) {
	existence := false
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return false, err
	}
	// fmt.Println("result: ", queryResults)

	fmt.Println(string(queryResults))
	fmt.Println(len(string(queryResults)))

	if len(string(queryResults)) > 2 {
		fmt.Println("here")
		existence = true
	}
	return existence, nil

}

func isExistPriv(stub shim.ChaincodeStubInterface, collection, queryString string) (bool, error) {
	existence := false
	queryResults, err := getQueryResultForQueryStringPriv(stub, collection, queryString)
	if err != nil {
		return false, err
	}
	// fmt.Println("result: ", queryResults)

	fmt.Println(string(queryResults))
	fmt.Println(len(string(queryResults)))

	if len(string(queryResults)) > 2 {
		fmt.Println("here")
		existence = true
	}
	return existence, nil
}
