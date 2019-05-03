package main

/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import (
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
	nargs = append(nargs, []byte("addOrgs"), epBytes, []byte(args[1]), []byte(args[2]))

	resp := _invokeCC(stub, channelID, libEp, nargs)
	if resp.GetStatus() != 200 {
		logger.Info("Invoke Response Payload:", string(resp.GetPayload()))
		logger.Info("Invoke Response status", resp.GetStatus())
		return "", string(resp.GetPayload())
	}

	// set the modified endorsement policy for the key
	err = stub.SetStateValidationParameter(args[0], resp.GetPayload())
	if err != nil {
		return "", err.Error()
	}

	return "", ""
}

// delEP deletes the state-based endorsement policy for the key altogether
func delEP(stub shim.ChaincodeStubInterface, key string) error {

	// set the modified endorsement policy for the key to nil
	err := stub.SetStateValidationParameter(key, nil)
	return err
}

// invokeCC is used for chaincode to chaincode invocation of a given cc on another channel
func _invokeCC(stub shim.ChaincodeStubInterface, channel, cc string, args [][]byte) pb.Response {
	ch := channel
	ccName := cc
	resp := stub.InvokeChaincode(ccName, args, ch)
	return resp
}
