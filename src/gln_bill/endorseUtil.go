package main

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	pb "github.com/hyperledger/fabric/protos/peer"
)

/*
EndorsementCC is an example chaincode that uses state-based endorsement.
In the init function, it creates two KVS states, one public, one private, that
can then be modified through chaincode functions that use the state-based
endorsement chaincode convenience layer. The following chaincode functions
are provided:
-) "addorgs": supply a list of MSP IDs that will be added to the
   state's endorsement policy
-) "delorgs": supply a list of MSP IDs that will be removed from
   the state's endorsement policy
-) "delep": delete the key-level endorsement policy for the state altogether
-) "listorgs": list the orgs included in the state's endorsement policy
*/

// addOrgs adds the list of MSP IDs from the invocation parameters
// to the state's endorsement policy

func addOrgs(stub shim.ChaincodeStubInterface, args []string) (string, string) {
	if len(args) < 2 {
		return "", "No orgs to add specified"
	}

	// get the endorsement policy for the key
	var epBytes []byte
	var err error
	var nargs [][]byte

	epBytes, err = stub.GetStateValidationParameter(args[0])
	fmt.Println("EP Bytes", epBytes)
	nargs = append(nargs, []byte("addOrgs"), epBytes, []byte(args[1]), []byte(args[2]))
	fmt.Println("nargs", nargs)

	resp := _invokeCC(stub, channelID, libEp, nargs)
	if err != nil {
		logger.Error(err)
		return "", err.Error()
	}
	fmt.Println("resp:", string(resp.GetPayload()))
	fmt.Println("resp payload:", string(resp.Payload))

	// set the modified endorsement policy for the key
	err = stub.SetStateValidationParameter(args[0], resp.GetPayload())

	if err != nil {
		return "", err.Error()
	}

	return "", ""
}

// delOrgs removes the list of MSP IDs from the invocation parameters
// from the state's endorsement policy
func delOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) < 2 {
		return shim.Error("No orgs to delete specified")
	}

	// get the endorsement policy for the key
	var epBytes []byte
	var err error

	epBytes, err = stub.GetStateValidationParameter(parameters[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(epBytes)
	return shim.Success([]byte{})
}

// listOrgs returns the list of organizations currently part of
// the state's endorsement policy
func listOrgs(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) < 1 {
		return shim.Error("No key specified")
	}

	// get the endorsement policy for the key
	// var epBytes []byte
	// var err error

	// epBytes, err = stub.GetStateValidationParameter(parameters[0])

	// ep, err := statebased.NewStateEP(epBytes)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// get the list of organizations in the endorsement policy
	// orgs := ep.ListOrgs()
	// orgsList, err := json.Marshal(orgs)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	return shim.Success(nil)
}

// delEP deletes the state-based endorsement policy for the key altogether
func delEP(stub shim.ChaincodeStubInterface) pb.Response {
	_, parameters := stub.GetFunctionAndParameters()
	if len(parameters) < 1 {
		return shim.Error("No key specified")
	}

	// set the modified endorsement policy for the key to nil
	var err error

	err = stub.SetStateValidationParameter(parameters[0], nil)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte{})
}

// invokeCC is used for chaincode to chaincode invocation of a given cc on another channel
func _invokeCC(stub shim.ChaincodeStubInterface, channel, cc string, args [][]byte) pb.Response {
	ch := channel
	ccName := cc
	resp := stub.InvokeChaincode(ccName, args, ch)
	return resp
}
