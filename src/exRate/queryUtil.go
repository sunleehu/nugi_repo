package main

import (
	"bytes"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

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

	logger.Debug("Query Result:", buffer.String())

	return buffer.Bytes(), nil
}

func getQueryResultForLatest(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	logger.Debug("QueryString:", queryString)

	//Get Query Result
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	//Query Result Iterator
	if resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		buffer.WriteString("{")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
	}
	buffer.WriteString("]")

	defer resultsIterator.Close()

	// fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	logger.Debug("Query Result:", buffer.String())

	return buffer.Bytes(), nil

}
