package main

import (
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkAtoi(str string) bool {
	_, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return true
	}
	return false
}

func isExist(stub shim.ChaincodeStubInterface, queryString string) (bool, error) {
	existence := false
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return false, err
	}

	if len(string(queryResults)) > 2 {
		existence = true
	}
	return existence, nil
}

func checkBlank(str string) string {
	return strings.TrimSpace(str)
}
