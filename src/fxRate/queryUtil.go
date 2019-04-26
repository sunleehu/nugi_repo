package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getDataByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]exchangeRate, error) {

	logger.Debug("QueryString:", queryString)
	// Get Query Result
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// Query Result Iterator
	var dataList []exchangeRate
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var data exchangeRate
		err = json.Unmarshal([]byte(queryResponse.Value), &data)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, data)
	}
	return dataList, nil
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	logger.Debug("QueryString:", queryString)

	// Get Query Result
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		if comma == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		comma = true
	}
	buffer.WriteString("]")

	logger.Debug("Query Result:", buffer.String())
	return buffer.Bytes(), nil
}

func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {

	logger.Debug("QueryString :\n%s\n", queryString)

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

	logger.Debug("QueryResults :\n%s\n", bufferWithPaginationInfo.String())
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

// $or 는 Full Query 이므로 성능저하 가능성 있음
// func multiQueryMaker(key string, data []string) string {
// 	var selectKey string
// 	comma := false
// 	for i := 0; i < len(data); i++ {
// 		if comma {
// 			selectKey = selectKey + ", "
// 		}
// 		selectKey = selectKey + fmt.Sprintf(`{"%s":"%s"}`, key, data[i])
// 		comma = true
// 	}
// 	selector := fmt.Sprintf(`"$or":[%s]`, selectKey)
// 	return selector
// }

func queryResponseStructMaker(result [][]byte, bookmark string, recordCnt int32) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("{\"BC_RES_DATA\":[")
	bArrayMemberAlreadyWritten := false
	for i := 0; i < len(result); i++ {
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.Write(result[i])
		bArrayMemberAlreadyWritten = true
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
