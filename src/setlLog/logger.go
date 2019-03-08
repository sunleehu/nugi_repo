package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type errorObj struct {
	ErrCd      string      `json:"errCd"`
	Message    string      `json:"errMsg"`
	AddMessage interface{} `json:"addMsg,string"`
}

func errMessage(code string, message interface{}) string {
	errList := map[string]string{
		"BCCE0001": "Invoke Error",
		"BCCE0002": "Permission Error",
		"BCCE0003": "JSON Decoding Error",
		"BCCE0004": "JSON Encoding Error",
		"BCCE0005": "Missing Required Value Error",
		"BCCE0006": "Duplicate Value Exists Error",
		"BCCE0007": "Invalid Arguments Error",
		"BCCE0008": "Query Error",
		"BCCE0009": "CouchDB Insert Error",
		"BCCE0010": "CouchDB Update Error",
	}
	eobj := errorObj{code, errList[code], message}
	logger.Error(code, errList[code], message)

	ee, _ := json.Marshal(eobj)
	return string(ee)
}

func setLogLevel(str string) pb.Response {
	lv, err := shim.LogLevel(str)
	if err != nil {
		return shim.Error("ERR")
	}
	logger.SetLevel(lv)
	logger.Notice("Now Log Mode:", str)
	return shim.Success(nil)
}
