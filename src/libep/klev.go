package main

/*This Chaincode is library. Using for Key Level Endorsement Policy Make
* Library can not get any state.
* Only return Policy Bytes
 */
import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/statebased"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type libKlvCC struct {
}

var logger = shim.NewLogger("LIBEP")

func main() {
	err := shim.Start(new(libKlvCC))
	if err != nil {
		logger.Error("Error starting libep chaincode : %s", err)
	}
}

func (t *libKlvCC) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

func (t *libKlvCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println()
	logger.Info("Func :", function)
	logger.Info("Args :", args)

	//args[0] = epBytes
	//args[1:]... = "Org1MSP","Org2MSP"
	if function == "addOrgs" {
		return addOrgs(stub, args)
	} else if function == "delOrgs" {
		return delOrgs(stub, args)
	} else if function == "listOrgs" {
		return listOrgs(stub, args)
	}

	errM := errMessage("BCCE0001", "Received unknown function invocation "+function)
	return shim.Error(errM)
}

func addOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	epBytes := []byte(args[0])
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	err = ep.AddOrgs(statebased.RoleTypeMember, args[1:]...)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	resBytes, err := ep.Policy()
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	logger.Info("AddOrgs Complete")
	return shim.Success(resBytes)
}

// delOrgs removes the list of MSP IDs from the invocation parameters
// from the state's endorsement policy
func delOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// get the endorsement policy for the key
	epBytes := []byte(args[0])
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	// delete organizations from the endorsement policy of that key
	ep.DelOrgs(args[1:]...)

	resBytes, err := ep.Policy()
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	logger.Info("DelOrgs Complete")
	return shim.Success(resBytes)
}

// listOrgs returns the list of organizations currently part of
// the state's endorsement policy
func listOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// get the endorsement policy for the key
	epBytes := []byte(args[0])
	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	// get the list of organizations in the endorsement policy
	orgs := ep.ListOrgs()
	orgsList, err := json.Marshal(orgs)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	return shim.Success(orgsList)
}
