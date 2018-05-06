package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func fetch(stub shim.ChaincodeStubInterface, key string, existanceCheck bool) ([]byte, error) {
	res, err := stub.GetState(key)
	if err != nil {
		return nil, returnError("GetState Error", "Couldnt GetState for key- "+key)
	}
	if len(res) == 0 && existanceCheck {
		return nil, returnError("GetState Error", "No bytes exist for key- "+key)
	}
	return res, nil
}

func store(stub shim.ChaincodeStubInterface, key string, data []byte, overwrite bool) error {

	var err error
	if overwrite == false {
		res, err := stub.GetState(key)
		if err != nil {
			return returnError("PutState Overwrite Error", "Can't GetState data for key- "+key)
		}
		if len(res) != 0 {
			return returnError("PutState Overwrite Error", "Can't Overwrite data for key- "+key)
		}

	}
	err = stub.PutState(key, data)
	if err != nil {
		return returnError("PutState Error", "Couldn't PutState data for key- "+key)
	}

	return nil
}

func remove(stub shim.ChaincodeStubInterface, key string, checkExistance bool) error {
	if checkExistance {
		res, err := stub.GetState(key)
		if err != nil {
			return returnError("DelState Error", "Couldnt GetState for key- "+key)
		}
		if len(res) == 0 {
			return returnError("DelState Error", "No bytes exist for key- "+key)
		}
	}

	err := stub.DelState(key)
	if err != nil {
		return returnError("Delete-00", "Couldnt Delete State for key- "+key)
	}
	return nil
}

func returnError(code string, message string) error {

	errorObj := errorResponse{}
	errorObj.Status = code
	errorObj.Message = message

	errorAsbytes, err := json.Marshal(errorObj)
	if err != nil {
		return errors.New("ERROR WHILE BUILDING ERROR MESSAGE")
	}
	err = errors.New(string(errorAsbytes))

	return err
}

func checkUser(stub shim.ChaincodeStubInterface, userAddress string, Type string) (bool, user, error) {

	UserAcc := user{}
	userAsbytes, err := fetch(stub, userAddress, true)
	if err != nil {
		return false, UserAcc, err
	}
	err = json.Unmarshal(userAsbytes, &UserAcc)
	if err != nil {
		return false, UserAcc, returnError("checkUser-00", "Couldn't Umarshal User")
	}

	if UserAcc.Type != Type {
		return false, UserAcc, nil
	}

	return true, UserAcc, nil

}
