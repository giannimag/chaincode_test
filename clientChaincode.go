
package main

import (
 "errors"
 "fmt"

 "encoding/json"
 "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Person struct{
	Id, Name, Surname string
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 if function == "addNewPerson" {
  var dati Person
  var jsonAsBytes []byte 
  var err error
  
  if len(args) != 3 {
       return nil, errors.New("Incorrect number of arguments. Expecting 3")
  }
  
  // Initialize the chaincode
  dati.Id = args[0]
  dati.Name = args[1]
  dati.Surname = args[2]
  
  jsonAsBytes, _ = json.Marshal(dati)
  
  //fmt.Printf("idClient = %c, name = %c\n surname = %c", idClient, dati.name, dati.surname)
  fmt.Printf("********************************** -> %c", string(jsonAsBytes))
  
  
  // Write the state to the ledger
  err = stub.PutState(dati.Id, jsonAsBytes)
  if err != nil {
       return nil, err
  }
  
  return jsonAsBytes, nil
 }

 return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 if function != "query" {
      return nil, errors.New("Invalid query function name. Expecting \"query\"")
 }
 var idClient string
 //var dati Person
 var err error

 if len(args) != 1 {
      return nil, errors.New("Incorrect number of arguments. Expecting idClient to query")
 }

 idClient = args[0]

 // Get the state from the ledger
 Avalbytes, err := stub.GetState(idClient)
 if err != nil {
      jsonResp := "{\"Error\":\"Failed to get state for " + idClient + "\"}"
      return nil, errors.New(jsonResp)
 }

 if Avalbytes == nil {
      jsonResp := "{\"Error\":\"Nil value for " + idClient + "\"}"
      return nil, errors.New(jsonResp)
 }

 //dati = Person(Avalbytes)
 
 //jsonResp := "{\"idCliente\":\"" + idClient + "\",\"Name\":\"" + dati.name + "\",\"Surname\":\"" + dati.surname + "\"}"
 jsonResp := string(Avalbytes)
 fmt.Printf("Query Response:%s\n", jsonResp)
 return Avalbytes, nil
}

func main() {
 err := shim.Start(new(SimpleChaincode))
 if err != nil {
      fmt.Printf("Error starting Simple chaincode: %s", err)
 }
}