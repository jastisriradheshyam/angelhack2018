package main 


type bloodChaincode struct {
}


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
	}else if function == "requestForBottle"{
		return t.RequestForBottle(stub,args)
	}else if function == "increaseRequiredCount"{
		return t.IncreaseRequiredCount(stub,args)
	}else if function == "generateBottle"{
		return t.GenerateBottle(stub,args)
	}else if function == "displayAllData"{
		return t.DisplayAllData(stub,args)
	}else if function == "getById"{
		return t.GetById(stub,args)
	}else if function =="initUser"{
		return t.InitUser(stub,args)
	}else if function == "respondToRequest"{
		return t.RespondToRequest(stub,args)
	} else{
		return shim.Error("CHAINCODE ERROR NO FUNCTION FOUND !!!!!!!!!!"+function)
	}




}

func (t *SimpleChaincode) InitUser(stub shim.ChaincodeStubInterface) pb.Response {

	if len(args)!=2 {
		return shim.Error("need 2 args")
	}

	userAddress,err := getAccountAddress(stub)
	if err!=nil{
		return shim.Error(err.Error())
	}

	User:= user{}
	User.UserId = userAddress
	m:= make(map[string]stockRequirement)
	User.CurrentStock = m
	m2 := make(map[string]int)
	User.CurrentRequirement = m2
	
	email,err := getEmailId
	if err!=nil {
		return shim.Error(err.Error())
	}
	User.EmailId = email
	User.ContactPerson = args[0]
	User.Region = GLOBAL_REGION
	User.Type = args[1]

	userAsBytes,err = json.Marshal(User)
	if err!=nil {
		return shim.Error("COULDNT REGISTER USER",err.Error())
	}

	err = store(stub,userAddress,userAsBytes,false)
	if err!=nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)



}


