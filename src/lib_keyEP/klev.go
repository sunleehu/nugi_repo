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

var logger = shim.NewLogger("libKlvCC")

func main() {
	err := shim.Start(new(libKlvCC))
	if err != nil {
		fmt.Printf("Error starting endorse chaincode: %s", err)
	}
}

func (t *libKlvCC) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

func (t *libKlvCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	// fmt.Println("invoke is running " + function)
	logger = shim.NewLogger("CC:libKlvCC TX:" + stub.GetTxID())
	logger.Infof("Invoke is running %s", function)
	// Handle different functions

	//args[0] = epBytes
	//args[1:]... = "Org1MSP","Org2MSP"
	if function == "addOrgs" {
		return addOrgs(stub, args)
	} else if function == "delOrgs" {
		return delOrgs(stub, args)
	} else if function == "listOrgs" {
		return listOrgs(stub, args)
	}

	fmt.Println("could not find func: " + function) //error
	errObj := errMessage("BCCE0001", "Received unknown function invocation "+function)

	return shim.Error(errObj)
}

func addOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	epBytes := []byte(args[0])
	fmt.Println("EP Bytes", epBytes)

	ep, err := statebased.NewStateEP(epBytes)

	fmt.Println("EP ", ep)
	if err != nil {
		logger.Error(err)
		return shim.Error(err.Error())
	}

	fmt.Println("addorgs")

	err = ep.AddOrgs(statebased.RoleTypeMember, args[1:]...)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Org List", args[1:])

	epBytes, err = ep.Policy()

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println(epBytes)

	return shim.Success(epBytes)

}

// delOrgs removes the list of MSP IDs from the invocation parameters
// from the state's endorsement policy
func delOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// get the endorsement policy for the key
	epBytes := []byte(args[0])
	var err error

	ep, err := statebased.NewStateEP(epBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// delete organizations from the endorsement policy of that key
	ep.DelOrgs(args[1:]...)

	// epBytes, err = ep.Policy()
	epBytes, err = ep.PolicyOR()
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(epBytes)
}

// listOrgs returns the list of organizations currently part of
// the state's endorsement policy
func listOrgs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// get the endorsement policy for the key
	epBytes := []byte(args[0])
	var err error

	ep, err := statebased.NewStateEP(epBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	// get the list of organizations in the endorsement policy
	orgs := ep.ListOrgs()
	orgsList, err := json.Marshal(orgs)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(orgsList)
}
