package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type FirstccChaincode struct {
	Version string
}

type transaction struct {
	DeSeq    string `json:"de_seq"` // 거래일련번호
	Balance  int    `json:"balance,omitempty"`
	AddField int    `json:"add_field,omitempty"`

	// SndrLcGlnUnqCd string  // 회원 Local GLN 고유코드
	// RcvrDepoAmt    float64 // 사용 금액(상품 금액)
	// GlnDeDtm       string
}

type updateData struct {
	// AdjMnDsbReqNo  string
	// SndrLcGlnUnqCd string
	// AdjDtm         string
	// AdjCompYn      string
	Balance string
}
type evtResp struct {
	Filter []string
	Data   payloads
}

type payloads struct {
}

type queryArgs struct {
	DeSeq        string
	ReqStartTime string
	ReqEndTime   string
}

// var errLog *log.Logger

// func errLogger(code int, message interface{}) {
// 	// errLog := log.New(os.Stdout, "ERR: ", log.LstdFlags|log.LUTC)
// 	// m := fmt.Sprintf("%s", message)
// 	errLog.Println(code, message)
// }

var logger = shim.NewLogger("First")

func main() {
	// logger.Info("When?")
	err := shim.Start(new(FirstccChaincode))
	if err != nil {
		fmt.Printf("Error starting First chaincode: %s", err)
	}
}

func (t *FirstccChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// errLog = log.New(os.Stdout, "[ERR]: ", log.LstdFlags|log.LUTC)
	// logger.Info("INIT?")
	return shim.Success(nil)
}

func (t *FirstccChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	// fmt.Println("invoke is running " + function)
	logger = shim.NewLogger("CC:attest TX:" + stub.GetTxID())
	logger.Infof("Invoke is running %s", function)
	// Handle different functions
	if function == "insertTxLog" {
		return t.insertTxLog(stub, args)
	} else if function == "selectTxLog" {
		return t.selectTxLog(stub, args)
	} else if function == "selectPeriodTxLog" {
		return t.selectPeriodTxLog(stub, args)
	} else if function == "testQuery" {
		return t.testQuery(stub, args)
	} else if function == "addTest" {
		return t.addTest(stub, args)
	} else if function == "attrAssert" {
		return t.attrAssert(stub, args)
	} else if function == "checkAttr" {
		return t.checkAttr(stub, args)
	} else if function == "getPageQuery" {
		return t.getPageQuery(stub, args)
	} else if function == "deleteState" {
		return t.deleteState(stub, args)
	} else if function == "addOrgs" {
		return addOrgs(stub)
	} else if function == "getOtherPut" {
		return t.getOtherPut(stub, args)
	} else if function == "getPrivateData" {
		return t.getPrivateData(stub, args)
	} else if function == "putPrivateData" {
		return t.putPrivateData(stub, args)
	} else if function == "putGetData" {
		return t.putGetData(stub, args)
	} else if function == "putTest" {
		return t.putTest(stub, args)
	} else if function == "trans" {
		return t.trans(stub, args)
	} else if function == "putTransientData" {
		return t.putTransientData(stub, args)
	} else if function == "delPrivState" {
		return t.delPrivState(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error
	errObj := errMessage("BCCE0001", "Received unknown function invocation "+function)

	return shim.Error(errObj)
}
func (t *FirstccChaincode) getOtherPut(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("getOtherPut")

	var nargs [][]byte
	nargs = append(nargs, []byte(args[3]), []byte(args[4]), []byte(args[5]))
	resp := invokeCC(stub, args[1], args[2], nargs)
	logger.Info("resp:", string(resp.GetPayload()))
	if resp.GetStatus() != 200 {
		logger.Error(resp.GetMessage())
		return shim.Error(resp.GetMessage())
	}
	logger.Info(resp.GetMessage())
	logger.Info(resp.GetStatus())
	logger.Info(resp.GetPayload())
	err := stub.PutPrivateData(args[0], "getPriv", resp.GetPayload())
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) putTest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.PutPrivateData(args[0], args[1], []byte(args[2]))
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) getPrivateData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	resp, err := getQueryResultForQueryStringPriv(stub, args[0], args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(resp)
	return shim.Success(resp)
}

func (t *FirstccChaincode) putGetData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	resp, err := getQueryResultForQueryStringPriv(stub, args[0], args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(resp)
	var tx transaction
	fmt.Println(args[1])
	err = json.Unmarshal([]byte(args[2]), &tx)
	if err != nil {
		fmt.Println("type Error:", err)
		// errLogger(1, err)
		// logger.Error(err)
		return shim.Error("Check Your Json")

	}
	fmt.Println(":", tx)
	tx.Balance--
	txlogJSONBytes, err := json.Marshal(tx)

	fmt.Println(tx)

	// err = stub.PutState(tx.DeSeq, txlogJSONBytes)
	err = stub.PutPrivateData(args[3], tx.DeSeq, txlogJSONBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) putTransientData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	lis := stub.GetArgs()
	fmt.Println(len(lis))
	mapArgs, err := stub.GetTransient()
	fmt.Println(string(mapArgs["de_seq"]))
	if err != nil {
		return shim.Error(err.Error())
	}

	var tx transaction
	err = json.Unmarshal(mapArgs["de_seq"], &tx)
	if err != nil {
		fmt.Println("type Error:", err)
		// errLogger(1, err)
		// logger.Error(err)
		return shim.Error("Check Your Json")

	}
	fmt.Println(":", tx)
	txlogJSONBytes, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("encoding error:", err)
		// errLogger(1, err)
		// logger.Error(err)
		return shim.Error("encoding error")

	}
	err = stub.PutPrivateData(string(lis[1]), tx.DeSeq, txlogJSONBytes)
	if err != nil {
		fmt.Println("couch Error:", err)
		// errLogger(1, err)
		// logger.Error(err)
		return shim.Error("put private Data Error")

	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) trans(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	lis := stub.GetArgs()
	fmt.Println(len(lis))
	mapArgs, err := stub.GetTransient()
	fmt.Println(string(mapArgs["de_seq"]))
	if err != nil {

	}
	return shim.Success(nil)

}

func (t *FirstccChaincode) putPrivateData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var tx transaction
	fmt.Println(args[0])
	err := json.Unmarshal([]byte(args[0]), &tx)
	if err != nil {
		fmt.Println("type Error:", err)
		// errLogger(1, err)
		// logger.Error(err)
		return shim.Error("Check Your Json")

	}
	fmt.Println(":", tx)
	tx.Balance--
	// exs, e := isExist(stub, queryString)
	// fmt.Println(exs)
	// if exs {

	// 	fmt.Println("?????????????")
	// 	errObj := errMessage("BCCE0006", args[i])
	// 	return shim.Error(errObj)
	// }

	// if e != nil {
	// 	return shim.Error(err.Error())
	// }
	// PrintMemUsage()
	txlogJSONBytes, err := json.Marshal(tx)

	fmt.Println(tx)

	// err = stub.PutState(tx.DeSeq, txlogJSONBytes)
	err = stub.PutPrivateData("collection2", tx.DeSeq, txlogJSONBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) attrAssert(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := cid.AssertAttributeValue(stub, args[0], args[1])
	if err != nil {
		fmt.Println(err)
		return shim.Error("Error")
	}

	return shim.Success(nil)

}

func (t *FirstccChaincode) checkAttr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	is, err := checkIdentity(stub, args[0], args[1])
	if err != nil {
		return shim.Error("Err")
	}
	logger.Info("is", is)
	return shim.Success(nil)
}

func (t *FirstccChaincode) delPrivState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.DelPrivateData(args[0], args[1])
	if err != nil {
		return shim.Error("del err")
	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) deleteState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("del err")
	}
	return shim.Success(nil)
}

func (t *FirstccChaincode) getPageQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// is, err := checkIdentity(stub, args[0], args[1])
	// if err != nil {
	// 	return shim.Error("Err")
	// }
	// logger.Info("is", is)

	queryString := fmt.Sprintf("{\"selector\":{\"de_seq\":{\"$gte\": \"%s\"}}}", args[0])
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
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// ===========================================================================================
// addPaginationMetadataToQueryResults adds QueryResponseMetadata, which contains pagination
// info, to the constructed query results
// ===========================================================================================
func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {

	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}]")

	return buffer
}

func (t *FirstccChaincode) addTest(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var tx transaction
	err := json.Unmarshal([]byte(args[0]), &tx)
	if err != nil {
		fmt.Println("type Error:", err)
		// errLogger(1, err)
		// logger.Error(err)
		return shim.Error("Check Your Json")
	}
	msp, e := cid.GetMSPID(stub)
	if e != nil {
		fmt.Println("msp err", e)
	}
	fmt.Println("mspmsp", msp)
	queryString := fmt.Sprintf("{\"selector\":{\"DeSeq\": \"%s\"}}", tx.DeSeq)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("err")
	}
	if resultsIterator.HasNext() {
		val, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("query result Err")
		}
		err = json.Unmarshal(val.Value, &tx)
		tx.Balance = tx.Balance + 2
		enc, _ := json.Marshal(tx)
		err = stub.PutState(tx.DeSeq, enc)
	}

	return shim.Success(nil)
}

func (t *FirstccChaincode) insertTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ret []byte

	for i := 0; i < len(args); i++ {
		var tx transaction
		fmt.Println(args[i])
		err := json.Unmarshal([]byte(args[i]), &tx)
		if err != nil {
			fmt.Println("type Error:", err)
			// errLogger(1, err)
			// logger.Error(err)
			return shim.Error("Check Your Json")

		}
		fmt.Println(":", tx)
		// exs, e := isExist(stub, queryString)
		// fmt.Println(exs)
		// if exs {

		// 	fmt.Println("?????????????")
		// 	errObj := errMessage("BCCE0006", args[i])
		// 	return shim.Error(errObj)
		// }

		// if e != nil {
		// 	return shim.Error(err.Error())
		// }
		// PrintMemUsage()
		txlogJSONBytes, err := json.Marshal(tx)

		fmt.Println(tx)

		err = stub.PutState(tx.DeSeq, txlogJSONBytes)
		// err = stub.PutPrivateData("privCollection", tx.DeSeq, txlogJSONBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		ret = txlogJSONBytes
	}
	// PrintMemUsage()

	return shim.Success(ret)
}

func (t *FirstccChaincode) testQuery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	pagesize, _ := strconv.Atoi(args[1])
	bookmark := args[2]
	iterator, repmeta, err := stub.GetQueryResultWithPagination(args[0], int32(pagesize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("iterator", iterator)
	fmt.Println("metadata", repmeta.XXX_Size())
	fmt.Println("met", repmeta)
	a, b := repmeta.Descriptor()
	fmt.Println("asdasd", a, b)
	fmt.Println("fetched:", repmeta.GetFetchedRecordsCount(), "::::", iterator.HasNext())

	return shim.Success(nil)
}

func (t *FirstccChaincode) selectTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// i, err := strconv.Atoi(args[1])
	// if err != nil {
	// 	return shim.Error("not int")
	// }

	// queryString := fmt.Sprintf("{\"selector\":{\"DeSeq\": \"%s\"}}", args[0])
	queryResults, err := getQueryResultForQueryString(stub, args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(string(queryResults))
	// var amap []map[string]interface{}
	// var datmap map[string]interface{}
	// var res result

	// json.Unmarshal(queryResults, &amap)
	// ma, err := json.Marshal(amap[0])

	// fmt.Println("asdklansdkashfkahfjkahsdjkfhkjasdfhkjadsfhkjasdfhjkda:     ", string(ma))

	// json.Unmarshal(queryResults, &res)

	// fmt.Println("asdakljflaksdfjklasdjfklasdjflkasjflkasjdfklsajdqweqw  :: : :: : :: : :: ", res.Record)

	return shim.Success(queryResults)
}

func (t *FirstccChaincode) selectPeriodTxLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// var tx transaction

	// err := json.Unmarshal([]byte(args[0]), &tx)
	// if err != nil {
	// 	fmt.Println("type Error:", err)
	// 	return shim.Error("Check Your Json")

	// }

	// fmt.Println("args[0]:", args[0])
	// fmt.Println("type: ", reflect.TypeOf(args[0]))
	// if len(args[0]) == 0 {
	// 	return shim.Error("args are null")

	// }
	// start, err := strconv.ParseUint(args[1], 10, 64)
	// end, err := strconv.ParseUint(args[2], 10, 64)
	var qArgs queryArgs

	json.Unmarshal([]byte(args[0]), &qArgs)

	queryString := fmt.Sprintf("{\"selector\": {\"$and\":[{\"GlnDeDtm\":{\"$gte\": \"%s\"}},{\"GlnDeDtm\":{\"$lte\": \"%s\"}}]}}", qArgs.ReqStartTime, qArgs.ReqEndTime)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

// func checkIdentity(stub shim.ChaincodeStubInterface, gln string) (string, error) {

// 	id, err := cid.New(stub)
// 	if err != nil {
// 		return "getID Error", err
// 	}
// 	mspid, err := id.GetMSPID()
// 	if err != nil {
// 		return "getMSPID Error", err
// 	}
// 	uid, err := cid.GetID(stub)
// 	fmt.Println("id: ", uid)

// 	fmt.Println("mspid: ", mspid)
// 	val, ok, err := id.GetAttributeValue("role")

// 	fmt.Println("role: ", val)
// 	if err != nil {
// 		// There was an error trying to retrieve the attribute
// 	}
// 	if !ok {
// 		// The client identity does not possess the attribute
// 	}
// 	// Do something with the value of 'val'

// 	return "", nil
// }

func checkIdentity(stub shim.ChaincodeStubInterface, attr string, code string) (bool, error) {
	id, err := cid.New(stub)
	if err != nil {
		return false, err
	}

	isRight := false

	val, ok, err := id.GetAttributeValue(attr)
	if err != nil {
		return false, err // There was an error trying to retrieve the attribute
	}
	if !ok {
		errMessage := fmt.Sprintf("There is no attribute: %s in cert", attr)
		return ok, errors.New(errMessage) // The client identity does not possess the attribute
	}
	fmt.Println("OK? val?", ok, val)
	// Do something with the value of 'val'
	logger.Info("VAL", val)

	if val == code {
		isRight = true
	}
	if (val == "ADMIN") || (val == "INT") {
		isRight = true
	}

	return isRight, nil
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
