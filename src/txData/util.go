package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkAtoi(str string) bool {
	_, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return true
	}
	return false
}

func checkGlnIntl(stub shim.ChaincodeStubInterface) (bool, string) {
	gln := false
	attr, exs, err := cid.GetAttributeValue(stub, "ACC_ROLE")
	if !exs {
		return false, "Certification does not have Attribute"
	} else if err != nil {
		logger.Error(err)
		return false, "Certification Error"
	}

	logger.Info("ACC_ROLE : ", attr)

	if attr == "INT" {
		gln = true
	} else {
		gln = false
	}
	return gln, ""
}

func isExist(stub shim.ChaincodeStubInterface, queryString string) (bool, []byte, error) {
	existence := false
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return false, nil, err
	}

	if len(string(queryResults)) > 2 {
		existence = true
	}
	return existence, queryResults, nil
}

func checkBlank(str string) string {
	return strings.TrimSpace(str)
}

func rmvDupVal(arr []string) []string {
	strmap := map[string]bool{}
	for _, elem := range arr {
		strmap[elem] = true
	}

	keys := []string{}

	for key := range strmap {
		keys = append(keys, key)

	}
	return keys
}

func sha256Hash(value string) string {
	sha := sha256.New()
	sha.Write([]byte(value))
	hash := hex.EncodeToString(sha.Sum(nil))
	return hash
}

func collectionMaker(var1, var2 string) string {
	var collectionName string
	if var1 < var2 {
		collectionName = var1 + var2
	} else {
		collectionName = var2 + var1
	}
	return collectionName
}
