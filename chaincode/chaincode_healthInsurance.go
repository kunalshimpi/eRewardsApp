/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
		 http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	//"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	//"github.com/hyperledger/fabric/core/crypto/primitives"
)

// SimpleHealthChaincode example simple Chaincode implementation
type SimpleHealthChaincode struct {
}
type eReward struct{
	Points string `json:"points"`
	Hash string `json:"hash"`
	Signature int `json:"signature"`
	Tx_ID string `json:"tx_id"`
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	//primitives.SetSecurityLevel("SHA", 256)
	err := shim.Start(new(SimpleHealthChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleHealthChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("**********Inside Init*******");
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	/*adminCert, err := stub.GetCallerMetadata()

	if err!= nil{
		return nil, errors.New("Error Getting Metadata")
	}
	if len(adminCert) == 0 {
		return nil, errors.New("Admin Certificate is Empty!")
	}
	stub.PutState("admin", adminCert)

	fmt.Println("Admin is [%x] : ", adminCert)

	fmt.Println("Assigning Amount for admin!")
	_, err = stub.InsertRow("InsuranceAmount", shim.Row{
		Columns: []*shim.Column {
			&shim.Column{Value: &shim.Column_Bytes{Bytes:[]byte("admin")}},
			&shim.Column{Value: &shim.Column_Int64{Int64:1000}}},
	})
	if err != nil {
		return nil, errors.New("Failed to Assign Amount!")
	}*/
	fmt.Println("Init Finished!")

	return nil, nil
}

func (t *SimpleHealthChaincode) assign(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("assign is running ")

	if len(args) != 3{
		return nil, errors.New("Expected 3 arguments!")
	}

	points, _ := strconv.Atoi(args[1]) //parse to int
	user := args[0]
	Sign_assigner := args[2]

	eRewardAsBytes, err := stub.GetState(user)
	if err != nil {
		return nil, errors.New("Failed to get eReward Object")
	}
	if eRewardAsBytes == nil {
		t.init_eReward(stub, args) //will create key/value with eReward stuct
	}else{
		//update existing eReward struct
	eRewardAsBytes, err := stub.GetState(user)
	if err != nil {
		return nil, errors.New("Failed to get struct")
	}
	res := eReward{}
	json.Unmarshal(eRewardAsBytes, &res)

	oldPoints,_ := strconv.Atoi(res.Points)
	newPoints := oldPoints + points
	res.Points = strconv.Itoa(newPoints)

	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(user, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	_, err1 := stub.InsertRow("ActivityTable", shim.Row{
		Columns: []*shim.Column {
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:strconv.Itoa(points)}},
			&shim.Column{Value: &shim.Column_String_{String_:user}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:Sign_assigner}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:"assign"}},
			},
	})
	if err1 != nil{
		return nil, errors.New("Insert Row failed!")
	}
  }

	/*adminCert, err := stub.GetState("admin")
	if err != nil{
		return nil, errors.New("Failed to get admin Certificate!")
	}

	ok, err := t.isCaller(stub, adminCert)
	if err != nil {
		return nil, errors.New("Failed to Check Certificates!")
	}
	if !ok {
		return nil, errors.New("Only Admin can call Approve function")
	}
*/
	// fmt.Println("Adding Transaction Detail")
	//
	// ok, err = stub.InsertRow("ActivityTable", shim.Row{
	// 	Columns: []*shim.Column {
	// 		&shim.Column{Value: &shim.Column_String_String_:nil}},
	// 		&shim.Column{Value: &shim.Column_String_{String_:nil}},
	// 		&shim.Column{Value: &shim.Column_String_{String_:points}},
	// 		&shim.Column{Value: &shim.Column_String_{String_:user}},
	// 		&shim.Column{Value: &shim.Column_String_{String_:nil}},
	// 		&shim.Column{Value: &shim.Column_String_{String_:Sign_assigner}},
	// 		&shim.Column{Value: &shim.Column_String_{String_:timestamp}},
	// 		},
	// })
	// if !ok && err == nil {
	// 	return nil, errors.New("Failed to insert transaction row!")
	// }

	fmt.Println("Assign Finished")
	return nil, err
}

func (t *SimpleHealthChaincode)init_eReward(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	fmt.Println("init_eReward is running ")

	if len(args) != 3{
		return nil, errors.New("Expected 3 arguments!")
	}

	points, _ := strconv.Atoi(args[1]) //parse to int
	user := args[0]
	Sign_assigner := args[2]

	err := stub.CreateTable("ActivityTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name:"Tx_ID",Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name:"From",Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name:"RewardPoint",Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name:"To",Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name:"Sign_rcvr",Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name:"Sign_assigner",Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name:"timestamp",Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name:"Reason",Type: shim.ColumnDefinition_STRING, Key: false},
	})

	obj := `{"points": "` + strconv.Itoa(points) + `", "hash": "` + "nil" + `", "signature": ` + "nil" + `, "tx_id": "` + "nil" + `"}`
	err = stub.PutState("struct1", []byte(obj))
	if err != nil {
		return nil, err
  }

	_, err1 := stub.InsertRow("ActivityTable", shim.Row{
		Columns: []*shim.Column {
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:strconv.Itoa(points)}},
			&shim.Column{Value: &shim.Column_String_{String_:user}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:Sign_assigner}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:"assign init"}},
			},
	})
	if err1 != nil {
		return nil, errors.New("InsertRow failed in init")
	}
	return nil,nil
}

func (t *SimpleHealthChaincode) redeem(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	if len(args) != 3{
		return nil, errors.New("Expected 3 arguments!")
	}
	b_entity := args[0]
	user := args[1]
	redeemPoints,_ := strconv.Atoi(args[2])

	/*
//check caller and Owner
	ok, err := t.isCaller(stub, sender)
	if err != nil {
		return nil, errors.New("Failed checking sender & caller identity")
	}
	if !ok {
		return nil, errors.New("The caller is not the owner of the amount")
	}*/

//change assets  of sender
	eRewardAsBytes, err := stub.GetState(user)
	if err != nil {
 		return nil, errors.New("Failed to get struct")
	}
	res := eReward{}
	json.Unmarshal(eRewardAsBytes, &res)

	oldPoints,_ := strconv.Atoi(res.Points)
	newPoints := oldPoints - redeemPoints
	res.Points = strconv.Itoa(newPoints)

	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(user, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	_, err1 := stub.InsertRow("ActivityTable", shim.Row{
		Columns: []*shim.Column {
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:strconv.Itoa(redeemPoints)}},
			&shim.Column{Value: &shim.Column_String_{String_:user}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:b_entity}},
			&shim.Column{Value: &shim.Column_String_{String_:"test"}},
			&shim.Column{Value: &shim.Column_String_{String_:"Redeem"}},
			},
	})
	if err1 != nil{
		return nil, errors.New("Insert Row failed!")
	}

	return nil, nil
}

func (t *SimpleHealthChaincode) isCaller(stub shim.ChaincodeStubInterface, certificate []byte) (bool, error) {
	// Verify \sigma=Sign(certificate.sk, tx.Payload||tx.Binding) against certificate.vk
	fmt.Println("isCaller is Running!")

	sigma, err := stub.GetCallerMetadata()
	if err != nil {
		return false, errors.New("Failed to get Metadata")
	}
	payload, err := stub.GetPayload()
	if err != nil {
		return false, errors.New("Failed to get payload")
	}
	binding, err := stub.GetBinding()
	if err != nil {
		return false, errors.New("Failed to get binding")
	}

	fmt.Println("Certificate : [%x]", certificate)
	fmt.Println("Sigma : [%x]", sigma)
	fmt.Println("Payload : [%x]", payload)
	fmt.Println("Binding : [%x]", binding)

	ok, err := stub.VerifySignature(
		certificate,
		sigma,
		append(payload, binding...),
	)
	if err != nil {
		return ok, errors.New("Failed Verifying signatures")
	}
	if !ok {
		fmt.Println("Signatures Does Not Match!")
	}
	fmt.Println("finished isCaller")

	return ok, err
}


// Invoke is our entry point to invoke a chaincode function
func (t *SimpleHealthChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "assign" {													//initialize the chaincode state, used as reset
		return t.assign(stub, args)
	} else if function == "redeem"{
		return t.redeem(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleHealthChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleHealthChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

	if len(args) != 1 {
		return nil, errors.New("Expected 1 argument!")
	}
	user := args[0]
	//fmt.Println("Finding [%x]",string(applicant))
/*
	var columns []shim.Column
	col := shim.Column{Value: &shim.Column_String{String: user}}
	columns = append(columns,col)

	row, err := stub.GetRow("InsuranceAmount",columns)
	while(row != nil)
	if err != nil {
		return nil, errors.New("Cannot retrieve Rows")
	}
*/
	valAsbytes, err := stub.GetState(user)									//get the var from chaincode state
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + user + "\"}"
		return nil, errors.New(jsonResp)
	}

	if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Nil struct for " + user + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + user + "\",\"Amount\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	fmt.Println("Finished Query function")
	return valAsbytes, nil


	//rowString := fmt.Sprintf("%s", row)
	//return []byte(rowString), nil
	//return row.Columns[0].GetBytes(), nil

}
