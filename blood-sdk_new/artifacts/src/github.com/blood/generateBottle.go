package main
import (
"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)


const TYPE_BLOODBANK string = "BLOOD_BANK"
const STATUS_ACTIVE string = "ACTIVE"
const STATUS_INTRANSIT string = "INTRANSIT"
const GLOBAL_REGION string = "South-Delhi"
const RESPONSE_PROVIDE string = "PROVIDE"
const RESPONSE_CANCEL string = "CANCEL"
const BLOOD_APLUS string = "A+"
const BLOOD_BPLUS string ="B+"
const BLOOD_ABPLUS string = "AB+"
const BLOOD_OPLUS string = "O+"

const MUST_COUNT = 5

//RespondToRequest
func (t *SimpleChaincode) RespondToRequest(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	response := args[0]
	index,err := strconv.Atoi(args[1])
	if err!=nil {
		return shim.Error("need an integer"+err.Error())
	}


	userAddress,err := getAccountAddress(stub)
		if err!=nil{
			return shim.Error(err.Error())
		}

	// userAsBytes,err := fetch(stub,fromUser,true)
	// if err!=nil {
	// 	return shim.Error("COULNT GETSTATE USER"+err.Error())
	// }

	User := user{}


	donatorAsBytes,err :=fetch(stub,userAddress,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(donatorAsBytes,&User)

	mystr := User.Giving[index]
		words := strings.Fields(mystr)
		bloodGroup := words[0]
		takerId := words[2]
		bottleCount,err := strconv.Atoi(words[1])
		if err!=nil {
			return shim.Error("need a integer"+err.Error())
		}

		Taker := user{}


	takerAsBytes,err :=fetch(stub,takerId,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(takerAsBytes,&Taker)


	if response == RESPONSE_CANCEL{

	User.Giving = append(User.Giving[:index], User.Giving[index+1:]...)
	Taker.Asking = append(Taker.Asking[:index], Taker.Asking[index+1:]...)

	finalUserAsBytes,err := json.Marshal(User)
	if err!=nil{
		return shim.Error("couldnt marshal user back"+err.Error())
	}
	finalTakerAsBytes,err := json.Marshal(Taker)
	if err!=nil {
		return shim.Error("couldnt marshal taker")
	}

	err = store(stub,userAddress,finalUserAsBytes,true)
	err = store(stub,takerId,finalTakerAsBytes,true)
	if err!=nil{
		return shim.Error("couldnt put state "+err.Error())
	}



	
		


	}else if response == RESPONSE_PROVIDE {

		if User.CurrentStock[bloodGroup].Count > (User.CurrentStock[bloodGroup].MustCount+User.CurrentRequirement[bloodGroup]+bottleCount) {

			for i:=0;i<bottleCount;i++{
				bottleId := User.CurrentStock[bloodGroup].BottleMap[0]
				bottleAsBytes,err := fetch(stub,bottleId,true)
				if err!=nil {
					return shim.Error("NO BOTTLE FOUND ")
				}
				Bottle:= bottle{}
				err = json.Unmarshal(bottleAsBytes,&Bottle)
				if err!=nil {
					return shim.Error("COULDNT MARSHAL BOTTLE"+err.Error())
				}
				Bottle.CurrentOwner = takerId
				Bottle.Trail += "BOTTLE HAS BEEN TRANSFERED TO "+ Taker.ContactPerson+" from "+ User.ContactPerson
				stockLevelUser,ok := User.CurrentStock[bloodGroup]
				if !ok{
					return shim.Error("error here at current stock user"+bloodGroup)
				}

				stockLevelUser.BottleMap = append(stockLevelUser.BottleMap[:0], stockLevelUser.BottleMap[1:]...)
				stockLevelUser.Count--;
				User.CurrentStock[bloodGroup] = stockLevelUser
				
				stockLevelTaker,ok := Taker.CurrentStock[bloodGroup]
				if !ok{
					return shim.Error("error here at current stock taker "+bloodGroup)
				}

				stockLevelTaker.BottleMap = append(stockLevelTaker.BottleMap[:0], stockLevelTaker.BottleMap[1:]...)
				stockLevelTaker.Count++;
				Taker.CurrentStock[bloodGroup] = stockLevelTaker

				// Taker.CurrentStock[bloodGroup].BottleMap = append(Taker.CurrentStock[bloodGroup].BottleMap,bottleId)
				// User.CurrentStock[bloodGroup].Count++
				


			}
			User.Giving = append(User.Giving[:index], User.Giving[index+1:]...)
				Taker.Asking = append(Taker.Asking[:index], Taker.Asking[index+1:]...)

				finalUserAsBytes,err := json.Marshal(User)
	if err!=nil{
		return shim.Error("couldnt marshal user back"+err.Error())
	}
	finalTakerAsBytes,err := json.Marshal(Taker)
	if err!=nil {
		return shim.Error("couldnt marshal taker")
	}

	err = store(stub,userAddress,finalUserAsBytes,true)
	err = store(stub,takerId,finalTakerAsBytes,true)
	if err!=nil{
		return shim.Error("couldnt put state "+err.Error())
	}



		}else{
			User.Giving = append(User.Giving[:index], User.Giving[index+1:]...)
			Taker.Asking = append(Taker.Asking[:index], Taker.Asking[index+1:]...)
		
			finalUserAsBytes,err := json.Marshal(User)
			if err!=nil{
				return shim.Error("couldnt marshal user back"+err.Error())
			}
			finalTakerAsBytes,err := json.Marshal(Taker)
			if err!=nil {
				return shim.Error("couldnt marshal taker")
			}
		
			err = store(stub,userAddress,finalUserAsBytes,true)
			err = store(stub,takerId,finalTakerAsBytes,true)
			if err!=nil{
				return shim.Error("couldnt put state "+err.Error())
			}

			return shim.Success(nil)


		}



	}

	return shim.Success(nil)



		}






func GetAdminCerts(stub shim.ChaincodeStubInterface) (*x509.Certificate, error) {
	var certi *x509.Certificate
	creator, err := stub.GetCreator()
	if err != nil {
		return certi, errors.New("Initsupplier couldn't get creator")
	}
	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)

	if err != nil {
		return certi, errors.New("COULDN UNMARSHAL SUPPLIER")
	}

	block, _ := pem.Decode(id.GetIdBytes())
	// if err !=nil {
	// 	return shim.Error(fmt.Sprintf("couldn decode"));
	// }
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return certi, errors.New("Initsupplier couldn't parse cert")
	}
	return cert, nil
}

//RequestForBottle
func (t *SimpleChaincode) RequestForBottle(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args)!=3{
		return shim.Error("need exact 3 args")
	}	
	
	numberOfBottles, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("chaincode:QueryCropByRange::ERROR01 need integer for count")
	}
		fromUser := args[1]
		BloodGroup := args[2]

		

		userAddress,err := getAccountAddress(stub)
		if err!=nil{
			return shim.Error(err.Error())
		}

	userAsBytes,err := fetch(stub,fromUser,true)
	if err!=nil {
		return shim.Error("COULNT GETSTATE USER"+err.Error())
	}

	User := user{}
	Donator := user{}

	donatorAsBytes,err :=fetch(stub,userAddress,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	err = json.Unmarshal(donatorAsBytes,&Donator)
	
	err = json.Unmarshal(userAsBytes,&User)
	if err!=nil{
		return shim.Error("COULDNT UNMARSHAL reciever")
	}

	User.Asking = append( User.Asking ,BloodGroup+" "+fmt.Sprint(numberOfBottles)+" " + userAddress+ " ")
	
	Donator.Giving = append(Donator.Giving,BloodGroup+" "+fmt.Sprint(numberOfBottles)+" " + userAddress+ " ")

	finalUserAsBytes,err:= json.Marshal(User)
	if err!=nil {
		return shim.Error("COULDNT MARSHAL DEFICIT ONE")
	}

	finalDonatorAsBytes,err := json.Marshal(Donator)
	if err!=nil{
		return shim.Error("Couldnt marshak donator")
	}

	err = store(stub,userAddress,finalUserAsBytes,true)
	err = store(stub,fromUser,finalDonatorAsBytes,true)

	if err!=nil {
		return shim.Error("COULDNT PUTSTATE AT END "+err.Error())
	}


	return shim.Success(nil);


}

func getAccountAddress(stub shim.ChaincodeStubInterface) (string, error) {

	cert, err := GetAdminCerts(stub)
	if err != nil {
		return "", errors.New("getAccountAddress couldn't get  creator certs")
	}

	userHash := sha256.Sum256([]byte(cert.Subject.CommonName + cert.Issuer.CommonName))
	userAddress := hex.EncodeToString(userHash[:])
	return userAddress, nil

}

func getEmailId(stub shim.ChaincodeStubInterface) (string,error){
	cert, err := GetAdminCerts(stub)
	if err != nil {
		return "", errors.New("getAccountAddress couldn't get  creator certs")
	}

	return cert.Subject.CommonName,nil

}

func (t *SimpleChaincode) SendBottle(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	bottleId := args[0]
	sentTo := args[1]

	bottleAsBytes,err := fetch(stub,bottleId,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	Bottle := bottle{}
	err = json.Unmarshal(bottleAsBytes,&Bottle)
	if err!=nil {
		return shim.Error("COULDNT MARSHAL BOTTLE")
	}

	recieverAsBytes,err := fetch(stub,sentTo,true)
	if err!=nil {
		return shim.Error("couldnt get reciever "+err.Error())
	}


	userAddress,err := getAccountAddress(stub)
	if err!=nil{
		return shim.Error(err.Error())
	}

	userAsBytes,err := fetch(stub,userAddress,true)
	if err!=nil {
		return shim.Error("COULNT GETSTATE USER"+err.Error())
	}

	User := user{}
	reciever := user{}
	err = json.Unmarshal(recieverAsBytes,&reciever)
	if err!=nil{
		return shim.Error("COULDNT UNMARSHAL reciever")
	}
	err = json.Unmarshal(userAsBytes,&User)
	if err!=nil{
		return shim.Error("COULDNT UNMARSHAL BANKs")
	}

	userStock := User.CurrentStock[Bottle.BloodGroup]

	bottleArray := userStock.BottleMap
	for i:=0;i<len(bottleArray);i++{
		if bottleArray[i]==bottleId {
			bottleArray = append(bottleArray[:i], bottleArray[i+1:]...)
			break;
		}
	}

	userStock.BottleMap = bottleArray
	userStock.Count --;
	User.CurrentStock[Bottle.BloodGroup] = userStock
	


	recieverStock := reciever.CurrentStock[Bottle.BloodGroup]
	
	recieverStock.Count +=1
	recieverStock.BottleMap = append(recieverStock.BottleMap,bottleId)
	User.CurrentStock[Bottle.BloodGroup]=recieverStock
	Bottle.Trail+= " BOTTLE TRANFERED FROM " +User.ContactPerson +" TO " + reciever.ContactPerson + "AS PER TXID"+stub.GetTxID()
	Bottle.CurrentOwner = reciever.ContactPerson

	finalRecieverAsBytes,err:= json.Marshal(reciever)
	if err!=nil{
		return shim.Error("COULDNT MARSHAL RECIEVER"+err.Error())
	}
	err = store(stub,sentTo,finalRecieverAsBytes,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	finalBankAsBytes,err := json.Marshal(User)
	if err!=nil{
		return shim.Error("COULDNT MARSHAL BANK"+err.Error())
	}
	err =  store(stub,userAddress,finalBankAsBytes,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	finalBottleAsBytes,err := json.Marshal(Bottle)
	if err!=nil{
		return shim.Error("COULDNT MARSHAL BOTTLE"+err.Error())
	}
	err =  store(stub,bottleId,finalBottleAsBytes,true)
	if err!=nil {
		return shim.Error(err.Error())
	}

	
return shim.Success(nil)

}

//GetById
func (t *SimpleChaincode) GetById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args)!=1{
		return shim.Error("need 1 args")
	}
	id := args[0]
	dataBytes,err := fetch(stub,id,true)
	if err!=nil{
		return shim.Error(err.Error())
	}
	
	return shim.Success(dataBytes)


}


func (t *SimpleChaincode) DisplayAllData(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	resultsIterator, err := stub.GetStateByRange(USER_LIST+"", USER_LIST+"~")
	if err != nil {
		return shim.Error(returnError("chaincode:QueryPolicyByRange::ERROR02", "couldnt marshal districtArray").Error())
	}
	var buffer bytes.Buffer
	buffer.WriteString(`[`)
	requiredCount := 4
	bArrayMemberAlreadyWritten := false
	count := 0
	for resultsIterator.HasNext() {
		if count >= requiredCount {
			break
		}
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		userAsBytes, err := fetch(stub, string(queryResponse.Value), true)
		if err != nil || len(userAsBytes) == 0 {
			return shim.Error(err.Error())
		}
		// buffer.WriteString("{\"key\":\"" + queryResponse.Key + "\"},")
		buffer.Write(userAsBytes)
		bArrayMemberAlreadyWritten = true
		count++

	}





	// var buffer bytes.Buffer
	// buffer.WriteString("[")
	// var mybool bool = false
	// for i:=0;i<len(args);i++{

	// 	userData,err := fetch(stub,args[0],false)
	// 	if err!=nil {
	// 		return shim.Error("user not found"+err.Error())
	// 	}
	// 	buffer.Write(userData)
	// 	mybool=true
	// 	if mybool{
	// 		buffer.WriteString(",")
	// 	}


	// }
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())

}


func (t *SimpleChaincode) GenerateBottle(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args)!=2 {
		return shim.Error("need exact 2 args")
	}
	bloodGroup := args[0]
	//count := args[1]


	requiredCount, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("chaincode:QueryCropByRange::ERROR01 need integer for count")
	}

	baseId:= 	stub.GetTxID()
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return shim.Error(returnError("InitCrops", "couldnt get timestamp for transaction").Error())
	}
	var millis int64
	millis = int64((timestamp.Seconds)*1000 + int64(timestamp.Nanos/1000000))
	epochTime := time.Unix(millis, 0)

	userAddress,err := getAccountAddress(stub)
	if err!=nil{
		return shim.Error(err.Error())
	}

	userAsBytes,err := fetch(stub,userAddress,true)
	if err!=nil {
		return shim.Error("COULNT GETSTATE USER"+err.Error())
	}

	bloodBank := user{}
	err = json.Unmarshal(userAsBytes,&bloodBank)
	if err!=nil{
		return shim.Error("COULDNT UNMARSHAL BANKs")
	}

	for i:=0;i<requiredCount;i++ {

		newBottle:= bottle{}
		newBottle.BottleId = baseId+fmt.Sprint(i)
		newBottle.BloodGroup = bloodGroup
		newBottle.DateOfPacking = epochTime
		newBottle.CurrentOwner = TYPE_BLOODBANK
		newBottle.Status = STATUS_ACTIVE
		newBottle.Trail = ""
		newBottle.Trail += "BOTTLE OF TYPE "+bloodGroup + "CREATED AT "+fmt.Sprint(epochTime.Date())+ "AS PER TXID "+ stub.GetTxID()

		newBottleAsBytes,err := json.Marshal(newBottle)
		if err!=nil {
			return shim.Error("COULDN MARSHAL NEW BOTTLE")
		}

		err = store(stub,newBottle.BottleId,newBottleAsBytes,false)
		if err!=nil{
			return shim.Error("Couldnt putstate"+fmt.Sprint(i))
		}

		userStorage := bloodBank.CurrentStock[bloodGroup]
		userStorage.Count +=1
		userStorage.BottleMap = append(userStorage.BottleMap,newBottle.BottleId)
		bloodBank.CurrentStock[bloodGroup]=userStorage
		

	}

	bloodBankAsBytes,err := json.Marshal(bloodBank)
	if err!=nil {
		shim.Error("COULDNT MARSHAL BLOOD BANK S BYTES")
	}

	err = store(stub,userAddress,bloodBankAsBytes,true)
	if err!=nil {
		shim.Error("COULDNT MARSHAL BLOOD BANK S BYTES")
	}

	return shim.Success(nil)



}



func (t *SimpleChaincode) IncreaseRequiredCount(stub shim.ChaincodeStubInterface, args[]string) pb.Response {

	if len(args)!=1{
		return shim.Error("need 1 args as blood group type")
	}

	bloodGroup := args[0]
	userAddress,err := getAccountAddress(stub)
	if err!=nil{
		return shim.Error(err.Error())
	}

	userAsBytes,err := fetch(stub,userAddress,true)
	if err!=nil {
		return shim.Error("COULNT GETSTATE USER"+err.Error())
	}

	User := user{}
	err = json.Unmarshal(userAsBytes,&User)
	if err!=nil{
		return shim.Error("COULDNT UNMARSHAL BANKs")
	}

	User.CurrentRequirement[bloodGroup]+=1;

	finalUserAsBytes,err := json.Marshal(User)
	if err!=nil{
		return shim.Error("CULNT MARSHAL USER")
	}

	err = store(stub,userAddress,finalUserAsBytes,true)
	if err!=nil {
		return shim.Error("COULDNT PUTSTATE "+err.Error())
	}

	return shim.Success(nil)

}