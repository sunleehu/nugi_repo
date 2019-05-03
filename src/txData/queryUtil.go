package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getDataByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]pubData, error) {

	logger.Debug("QueryString :", queryString)
	// Get Query Result
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// Query Result Iterator
	var dataList []pubData
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var data pubData
		err = json.Unmarshal([]byte(queryResponse.Value), &data)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	logger.Debug("QueryString :", queryString)
	// Get Query Result
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	// Query Result Iterator
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString, bookmark string, pageSize int32) ([]byte, error) {

	logger.Debug("QueryString :", queryString)

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
	return bufferWithPaginationInfo.Bytes(), nil
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

	buffer.WriteString("\"PAGE_COUNT\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"PAGE_NEXT_ID\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}")

	return buffer
}

func queryResponseStructMaker(result [][]byte, bookmark string, recordCnt int32) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("{\"BC_RES_DATA\":[")
	comma := false
	for i := 0; i < len(result); i++ {
		if comma == true {
			buffer.WriteString(",")
		}
		buffer.Write(result[i])
		comma = true
	}
	buffer.WriteString("],")
	buffer.WriteString("\"PAGE_NEXT_ID\":")
	buffer.WriteString("\"")
	buffer.WriteString(bookmark)
	buffer.WriteString("\",")
	buffer.WriteString("\"PAGE_COUNT\":")
	buffer.WriteString(fmt.Sprintf("%v", recordCnt))
	buffer.WriteString("}")

	return buffer.Bytes()
}

func multiSelector(key string, data []string) string {
	var selectKey string
	comma := false
	for i := 0; i < len(data); i++ {
		if comma {
			selectKey = selectKey + ", "
		}
		selectKey = selectKey + "\"" + data[i] + "\""
		comma = true
	}
	selector := fmt.Sprintf(`{"selector":{"%s":{"$in":[%s]}}, "use_index":["indexKeyDoc", "indexKey"]}`, key, selectKey)
	return selector
}

func multiQueryMaker(key string, data []string) string {
	var selectKey string
	comma := false
	for i := 0; i < len(data); i++ {
		if comma {
			selectKey = selectKey + ", "
		}
		selectKey = selectKey + fmt.Sprintf(`{"%s":"%s"}`, key, data[i])
		comma = true
	}
	selector := fmt.Sprintf(`{"$or":[%s]}`, selectKey)
	return selector
}

func getResultForPublicData(stub shim.ChaincodeStubInterface, queryString, bookmark string, pageSize int32) ([]byte, string, int32, error) {

	logger.Debug("QueryString :", queryString)
	// Get Query Result
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, "", 0, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	comma := false
	// Query Result Iterator
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, "", 0, err
		}
		// Add a comma before array members, suppress it for the first array member
		if comma == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		comma = true
	}
	buffer.WriteString("]")

	newBookmark := responseMetadata.GetBookmark()
	recordCnt := responseMetadata.GetFetchedRecordsCount()
	return buffer.Bytes(), newBookmark, recordCnt, nil
}

func getPrivQueryResultForQueryString(stub shim.ChaincodeStubInterface, collection, queryString string) ([]byte, error) {

	logger.Debug("QueryString :", queryString)

	// Get Query Result
	resultsIterator, err := stub.GetPrivateDataQueryResult(collection, queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	comma := false
	// Query Result Iterator
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if comma {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		comma = true
	}

	return buffer.Bytes(), nil
}

func getPrivateDataForKeys(stub shim.ChaincodeStubInterface, collection string, keys []string) ([]byte, error) {

	var buffer bytes.Buffer
	comma := false
	for i := 0; i < len(keys); i++ {
		data, err := stub.GetPrivateData(collection, keys[i])
		if err != nil {
			return nil, err
		}
		if data != nil {
			if comma {
				buffer.WriteString(",")
			}
			buffer.Write(data)
			comma = true
		}
	}

	return buffer.Bytes(), nil
}

func deletePrivateDataForKeys(stub shim.ChaincodeStubInterface, collection string, keys []string) error {

	for i := 0; i < len(keys); i++ {
		err := stub.DelPrivateData(collection, keys[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// func getPrivQueryResultAndKeys(stub shim.ChaincodeStubInterface, collection, queryString string, keys []string) ([]byte, error) {
// 	logger.Debug("QueryString:", queryString)

// 	// Get Query Result
// 	resultsIterator, err := stub.GetPrivateDataQueryResult(collection, queryString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	// buffer is a JSON array containing QueryRecords
// 	var buffer bytes.Buffer

// 	comma := false
// 	// Query Result Iterator
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		// JSON Decoding
// 		var tx transaction
// 		err = json.Unmarshal(queryResponse.Value, &tx)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for i := len(keys) - 1; i > 0; i-- {
// 			if tx.GlnTxHash == keys[i] {
// 				// Add a comma before array members, suppress it for the first array member
// 				if comma {
// 					buffer.WriteString(", ")
// 				}
// 				buffer.WriteString(string(queryResponse.Value))
// 				comma = true

// 				removeIndex(keys, i)
// 				break
// 			}
// 		}
// 	}

// 	logger.Debug("Query Result:", buffer.String())

// 	return buffer.Bytes(), nil
// }
