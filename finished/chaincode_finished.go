package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"	
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Bid struct {
	Status					string			`json:"status"`
	BidTime                 string          `json:"BidTime"`
	Amount					string			`json:"amount"`
	Content					string			`json:"content"`
	Username				string			`json:"username"`
	TxID                    string          `json:"txid"`
}
var bidIndexStr = "_bids"
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "register" {
		return t.register(stub, args)
	}  
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {
		return t.read(stub, args)
	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// register (aVoteToken)  - register a vote token
func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running register()")
	var  txid string
	var err error
	var jsonBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. a vote token to be registered.")
	}
	txid = stub.GetTxID()
	val := Bid{Status: "NEW", Content: "", BidTime: "", Username: "", Amount: "", TxID: txid}
	jsonBytes, _ = json.Marshal(val)
	id, err:= t.append_id(stub, bidIndexStr, args[0], false)
	err = stub.PutState(string(id), jsonBytes) // add the vote token into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) bid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running vote()")
	var key,  txid string
	var jsonBytes []byte
	val := new(Bid)

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 (votetoken, candidateid, timestamp, ipaddr, ua).")
	}
	key = args[3]
	jsonBytes, _ = t.read(stub, []string{args[3]})
	json.Unmarshal(jsonBytes, val)
	txid = stub.GetTxID()
	if val.Status == "NEW" {//val := Bid{Status: "NEW", Content: "", BidTime: "", Username: "", Amount: "", TxID: txid}
		// the vote token exists and has not been voted.
		val.Status = "BID"
		val.Content = args[1]
		val.BidTime = args[2]
		val.Username = args[3]
		val.Amount = args[4]
		val.TxID = txid
		jsonBytes, _ = json.Marshal(val)
		stub.PutState(key, jsonBytes)
		
	} else if val.Status == "BID" {
		stub.PutState("failure_"+args[3], []byte(txid))
		return nil, errors.New("DUPLICATED: the BID has already Done this UserId.")
	} else {
		stub.PutState("failure_"+args[3], []byte(txid))
		return nil, errors.New("ERROR: USER NOT REGISTER")
	}
	return nil, nil
}

// read - query function to read a vote token status
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running read()")
	var key string
	var err error
	var jsonBytes []byte
	
	key = args[0]
	jsonBytes, err = stub.GetState(key)
	if err != nil {
		return nil, errors.New("Error occurred when getting state of " + key)
	}
	
	return jsonBytes, nil
}
func (t *SimpleChaincode) append_id(stub shim.ChaincodeStubInterface, indexStr string, id string, create bool) ([]byte, error) {

	indexAsBytes, err := stub.GetState(indexStr)
	if err != nil {
		return nil, errors.New("Failed to get " + indexStr)
	}
	fmt.Println(indexStr + " retrieved")

	// Unmarshal the index
	var tmpIndex []string
	json.Unmarshal(indexAsBytes, &tmpIndex)
	fmt.Println(indexStr + " unmarshalled")

	// Create new id
	var newId = id
	if create {
		newId += strconv.Itoa(len(tmpIndex) + 1)
	}

	// append the new id to the index
	tmpIndex = append(tmpIndex, newId)
	jsonAsBytes, _ := json.Marshal(tmpIndex)
	err = stub.PutState(indexStr, jsonAsBytes)
	if err != nil {
		return nil, errors.New("Error storing new " + indexStr + " into ledger")
	}

	return []byte(newId), nil

}
func (t *SimpleChaincode) get_all_bids(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	bidIndexBytes, err := stub.GetState(bidIndexStr)
	if err != nil { return nil, errors.New("Failed to get bids index")}

	var bidIndex []string
	err = json.Unmarshal(bidIndexBytes, &bidIndex)
	if err != nil { return nil, errors.New("Could not marshal bid indexes") }

	var bids []Bid
	for _, bidId := range bidIndex {
		bytes, err := stub.GetState(bidId)
		if err != nil { return nil, errors.New("Not able to get bid") }

		var b Bid
		err = json.Unmarshal(bytes, &b)
		bids = append(bids, b)
	}

	bidsJson, err := json.Marshal(bids)
	if err != nil { return nil, errors.New("Failed to marshal bids to JSON")}

	return bidsJson, nil

}
