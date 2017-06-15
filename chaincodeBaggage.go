package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Position struct{
        BaggageLastPositionDescription, BaggageLastPositionLatitude, BaggageLastPositionLongitude, BaggageFlightID, BaggageFlightDestination string
        IsBoarded bool
}

type RefundInfo struct{
        EtheriumClientAddress, EtheriumAirlineAddress string
        RefundAmount float32
        ToRefund bool
}

type Baggage struct{
        Id, FlightID, SensorID, BaggageAirlineID, ClientID, TicketID, Destination string
		BaggageWeight float32
        IsFlightLanded, IsBaggageDelivered bool
        Tracking []Position
        Refund RefundInfo
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

// Transaction manage the baggages info
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var bag Baggage
	var jsonAsBytes []byte 
	var err error
	var id string
	
	if function == "addBaggage" {
		if len(args) != 14 {
			return nil, errors.New("Incorrect number of arguments. Expecting 14")
		}
	  
		// Initialize the chaincode
		bag.Id = args[0]
		bag.FlightID = args[1]
		bag.SensorID = args[2]
		bag.BaggageAirlineID = args[3]
		bag.ClientID = args[4]
		bag.TicketID = args[5]
		bag.Destination = args[6]
		w64,_ := strconv.ParseFloat(args[7],32)
		bag.BaggageWeight = float32(w64)


		tr := make([]Position, 1, 1)
		bag.Tracking = tr
		tr[0].BaggageLastPositionDescription = args[8]
		tr[0].BaggageLastPositionLatitude = args[9]
		tr[0].BaggageLastPositionLongitude  = args[10]

		var rf RefundInfo
		rf.EtheriumClientAddress = args[11]
		rf.EtheriumAirlineAddress = args[12]
		r64,_ := strconv.ParseFloat(args[13],32)
		rf.RefundAmount = float32(r64)
		bag.Refund = rf
		
		jsonAsBytes, _ = json.Marshal(bag)

		// Write the state to the ledger
		err = stub.PutState(bag.Id, jsonAsBytes)
		if err != nil {
		   return nil, err
		}

		return jsonAsBytes, nil
	}

	if function == "addBaggagePosition" {
		if len(args) != 7 {
			return nil, errors.New("Incorrect number of arguments. Expecting 7")
		}

		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		json.Unmarshal(Avalbytes, &bag)
		
		var lastPosition Position
		
		lastPosition.BaggageLastPositionDescription = args[1]
		lastPosition.BaggageLastPositionLatitude = args[2]
		lastPosition.BaggageLastPositionLongitude  = args[3]
		lastPosition.BaggageFlightID = args[4]
		lastPosition.BaggageFlightDestination = args[5]
		lastPosition.IsBoarded, _ = strconv.ParseBool(args[6])
		
		tr := &bag.Tracking
		*tr = append (*tr,lastPosition)
		
		jsonAsBytes, _ = json.Marshal(bag)
		
		err = stub.PutState(id, jsonAsBytes)
		if err != nil {
			return nil, err
		}
		
		return jsonAsBytes, nil  
	}
	
	if function == "airLandedEvent" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting Id")
		}

		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		json.Unmarshal(Avalbytes, &bag)
		
		bag.IsFlightLanded = true
		
		jsonAsBytes, _ = json.Marshal(bag)
		
		err = stub.PutState(id, jsonAsBytes)
		if err != nil {
			return nil, err
		}
		
		return jsonAsBytes, nil
	}
	
	if function == "baggageDeliveredEvent" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting Id")
		}

		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		json.Unmarshal(Avalbytes, &bag)
		
		bag.IsBaggageDelivered = true
		
		jsonAsBytes, _ = json.Marshal(bag)
		
		err = stub.PutState(id, jsonAsBytes)
		if err != nil {
			return nil, err
		}
		
		return jsonAsBytes, nil
	}
	if function == "updateRefundCondition" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting Id")
		}

		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		json.Unmarshal(Avalbytes, &bag)
		
		if ( bag.IsFlightLanded == true && bag.IsBaggageDelivered == false && bag.Destination != bag.Tracking[len(bag.Tracking)-1].BaggageFlightDestination ){
			bag.Refund.ToRefund = true
			
			jsonAsBytes, _ = json.Marshal(bag)
		
			err = stub.PutState(id, jsonAsBytes)
			if err != nil {
				return nil, err
			}
		}
		
		return nil, nil
	}
	
	errmsg := "Invalid Invoke function name. Expecting addBaggage or addBaggagePosition or airLandedEvent or baggageDeliveredEvent or updateRefundCondition"
	return []byte(errmsg), errors.New(errmsg)

}


// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting the baggageId to query")
	}
	
	if function == "getBaggageInfo" {
		var id string
		var err error
		
		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		jsonResp := string(Avalbytes)
		fmt.Printf("Query Response:%s\n", jsonResp)
		return Avalbytes, nil
	}	
	
	if function == "getBaggageLastPosition" {
		var id string
		var err error
		var bag Baggage 
  
		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		json.Unmarshal(Avalbytes,&bag)		
		
		jsonResp := "{\"BaggageLastPositionDescription\":\"" + bag.Tracking[len(bag.Tracking)-1].BaggageLastPositionDescription + 
		"\",\"BaggageLastPositionLatitude\":\"" + bag.Tracking[len(bag.Tracking)-1].BaggageLastPositionLatitude + 
		"\",\"BaggageLastPositionLongitude\":\"" + bag.Tracking[len(bag.Tracking)-1].BaggageLastPositionLongitude + 
		"\",\"BaggageFlightID\":\"" + bag.Tracking[len(bag.Tracking)-1].BaggageFlightID + 
		"\",\"BaggageFlightDestination\":\"" + bag.Tracking[len(bag.Tracking)-1].BaggageFlightDestination + 
		"\",\"isBoarded\":\"" + strconv.FormatBool(bag.Tracking[len(bag.Tracking)-1].IsBoarded) + 		
		"\"}"
		
		
		fmt.Printf("Query Response:%s\n", jsonResp)
		
		return []byte(jsonResp), nil	 
	}	
 	if function == "getRefundCondition" {
		var id string
		var err error
		var bag Baggage 
  
		id = args[0]
		// Get the state from the ledger
		Avalbytes, err := stub.GetState(id)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		if Avalbytes == nil {
			jsonResp := "{\"Error\":\"Nil value for " + id + "\"}"
			return nil, errors.New(jsonResp)
		}
		
		json.Unmarshal(Avalbytes,&bag)		
		
		jsonResp := "{\"refundCondition\":\"" + strconv.FormatBool(bag.Refund.ToRefund) + 
		"\",\"EtheriumClientAddress\":\"" + bag.Refund.EtheriumClientAddress + 
		"\",\"EtheriumAirlineAddress\":\"" + bag.Refund.EtheriumAirlineAddress + 
		"\",\"RefundAmount\":\"" + strconv.FormatFloat(float64(bag.Refund.RefundAmount), 'f', -1, 32) + 
		"\"}"
		
		fmt.Printf("Query Response:%s\n", jsonResp)
		
		return []byte(jsonResp), nil
	}
	
	errmsg := "Invalid query function name. Expecting getBaggageInfo or getBaggageLastPosition or getRefundCondition"
	return []byte(errmsg), errors.New(errmsg)
}


func main() {
 err := shim.Start(new(SimpleChaincode))
 if err != nil {
      fmt.Printf("Error starting Simple chaincode: %s", err)
 }
}