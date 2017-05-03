package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type VoteStateValue struct {
	Status      string `json:"status"`
	Candidateid string `json:"candidateid"`
	Timestamp   string `json:"timestamp"`
	Ipaddr      string `json:"ipaddr"`
	Ua          string `json:"ua"`
	TxID        string `json:"txid"`
}

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
	}  else if function == "bid" {
		return t.bid(stub, args)
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
	var key, txid string
	var err error
	//var jsonBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. a vote token to be registered.")
	}

	key = "vt_" + args[0] // vt_<voteToken>
	txid = stub.GetTxID()
	val := VoteStateValue{Status: "NEW", Candidateid: "", Timestamp: "", Ipaddr: "", Ua: "", TxID: txid}
	jsonBytes, _ := json.Marshal(val)
	err = stub.PutState(key, jsonBytes) // add the vote token into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


// vote (aVoteToken, aCandidateId, timestamp, ipaddr, ua)  - vote action
func (t *SimpleChaincode) bid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running vote()")
	var key,  txid string
	var jsonBytes []byte
	val := new(VoteStateValue)

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 (votetoken, candidateid, timestamp, ipaddr, ua).")
	}

	
	key = args[0]
	
	txid = stub.GetTxID()

	
		// the vote token exists and has not been voted.
		val.Status = "VOTED"
		val.Candidateid = args[1]
		val.Timestamp = args[2]
		val.Ipaddr = args[3]
		val.Ua = args[4]
		val.TxID = txid
		jsonBytes, _ = json.Marshal(val)
		stub.PutState(key, jsonBytes)
		

	
	return nil, nil
}


// read - query function to read a vote token status
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("running read()")
	var key, ret string
	var err error
	var jsonBytes []byte
	
	key = args[0]
	jsonBytes, err = stub.GetState(key)
	if err != nil {
		return nil, errors.New("Error occurred when getting state of " + key)
	}
	ret = string(jsonBytes)
	return []byte(ret), nil
}
