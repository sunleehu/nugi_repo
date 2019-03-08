package main

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkAtoi(str string) bool {
	_, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return true
	}
	return false
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

func getTimestamp(seconds int64) string {
	ts := time.Unix(seconds, 0)
	timestamp := ts.UTC().Format("20060102150405")
	return timestamp
}

func decimalMultiply(f1 float64, f2 float64) float64 {
	v1 := decimal.NewFromFloat(f1)
	v2 := decimal.NewFromFloat(f2)
	res, _ := v1.Mul(v2).Float64()
	return res
}

func decimalAdd(f1, f2 float64) float64 {
	v1 := decimal.NewFromFloat(f1)
	v2 := decimal.NewFromFloat(f2)
	res, _ := v1.Add(v2).Float64()
	return res
}

func decimalSub(f1, f2 float64) float64 {
	v1 := decimal.NewFromFloat(f1)
	v2 := decimal.NewFromFloat(f2)
	res, _ := v1.Sub(v2).Float64()
	return res
}

func decimalCeil(f1 float64, digit int) float64 {
	dig := math.Pow10(digit)
	v1 := decimal.NewFromFloat(f1)
	v2 := decimal.NewFromFloat(dig)
	res := v1.Mul(v2)
	res1 := res.Ceil()
	res2, _ := res1.Div(v2).Float64()
	return res2

}

func decimalTrunc(f1 float64, digit int32) float64 {
	v1 := decimal.NewFromFloat(f1)
	res, _ := v1.Truncate(digit).Float64()
	return res

}
