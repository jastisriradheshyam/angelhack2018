package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

const USER_TYPE_HOSPITAL string = "HOSPITAL"
const USER_TYPE_BANK string = "BLOODBANK"
const USER_LIST string = "USER_LIST"

// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {



	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "sendBottle" { //done
		return t.SendBottle(stub, args)
	} else if function == "requestForBottle" {
		return t.RequestForBottle(stub, args)
	} else if function == "increaseRequiredCount" {
		return t.IncreaseRequiredCount(stub, args)
	} else if function == "generateBottle" {
		return t.GenerateBottle(stub, args)
	} else if function == "displayAllData" {
		return t.DisplayAllData(stub, args)
	} else if function == "getById" {
		return t.GetById(stub, args)
	} else if function == "initUser" {
		return t.InitUser(stub, args)
	} else if function == "respondToRequest" {
		return t.RespondToRequest(stub, args)
	} else {
		return shim.Error("CHAINCODE ERROR NO FUNCTION FOUND !!!!!!!!!!" + function)
	}

}

func (t *SimpleChaincode) InitUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("need 2 args")
	}

	userAddress, err := getAccountAddress(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	User := user{}
	User.UserId = userAddress
	m := make(map[string]stockRequirement)
	var str []string
	m[BLOOD_APLUS] = stockRequirement{BLOOD_APLUS,str,0,MUST_COUNT}
	m[BLOOD_BPLUS] = stockRequirement{BLOOD_BPLUS,str,0,MUST_COUNT}
	m[BLOOD_ABPLUS] = stockRequirement{BLOOD_ABPLUS,str,0,MUST_COUNT}
	m[BLOOD_OPLUS] = stockRequirement{BLOOD_OPLUS,str,0,MUST_COUNT}
	User.CurrentStock = m
	m2 := make(map[string]int)
	User.CurrentRequirement = m2

	email, err := getEmailId(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	User.EmailId = email
	User.ContactPerson = args[0]
	User.Region = GLOBAL_REGION
	User.Type = args[1]

	userAsBytes, err := json.Marshal(User)
	if err != nil {
		return shim.Error("COULDNT REGISTER USER"+ err.Error())
	}

	err  = store(stub,USER_LIST+" "+email,[]byte(userAddress),true)
	if err!=nil{
		return shim.Error(err.Error())
	}

	err = store(stub, userAddress, userAsBytes, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}
