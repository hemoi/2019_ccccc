/* changed hyperledger fabric private DB tutorial */

package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Chaincode struct {
}

// for all of blockchain network user
type secret struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`  // data information
	Hash       []byte `json:"hash"`  // hash of original data
	Date       string `json:"date"`  // timestamp
	Owner      string `json:"owner"` // who has right of this data
}

// people who are permissioned
type secretOriginal struct {
	ObjectType string `json:"docType"`
	Name       string `json:"name"`
	Original   string `json:"original"` // original data
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode:%s", err)
	}
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	switch function {
	case "initSecret":
		return t.initSecret(stub, args)
	case "readSecret":
		return t.readSecret(stub, args)
	case "readOriginal":
		return t.readOriginal(stub, args)
	default:
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}

}

func (t *Chaincode) initSecret(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type secretTransientInput struct {
		Name     string `json:"name"`
		Hash     []byte `json:"hash"`
		Date     string `json:"date"`
		Owner    string `json:"owner"`
		Original string `json:"original"`
	}

	fmt.Println("= start init business secret protection system")

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private secret data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	if _, ok := transMap["secret"]; !ok {
		return shim.Error("secret must be a key in the transient map")
	}

	if len(transMap["secret"]) == 0 {
		return shim.Error("secret value in the transient map must be a non-empty JSON string")
	}

	var secretInput secretTransientInput
	err = json.Unmarshal(transMap["secret"], &secretInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["secret"]))
	}

	secretInput.Date = time.Now().Format("20060102150405")

	sha := sha512.New()
	sha.Write([]byte(secretInput.Original))
	sha.Write([]byte(secretInput.Owner))
	sha.Write([]byte(secretInput.Name))
	secretInput.Hash = sha.Sum(nil)

	if len(secretInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(secretInput.Hash) == 0 {
		return shim.Error("hash field must be a non-empty string")
	}
	if len(secretInput.Date) <= 0 {
		return shim.Error("date field must be a positive integer")
	}
	if len(secretInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if len(secretInput.Original) <= 0 {
		return shim.Error("original field must be a non-empty string")
	}

	secretAsBytes, err := stub.GetPrivateData("collectionSecret", secretInput.Name)
	if err != nil {
		return shim.Error("Failed to get secret: " + err.Error())
	} else if secretAsBytes != nil {
		fmt.Println("This secret already exists: " + secretInput.Name)
		return shim.Error("This secret already exists: " + secretInput.Name)
	}
	secret := &secret{
		ObjectType: "secret",
		Name:       secretInput.Name,
		Hash:       secretInput.Hash,
		Date:       secretInput.Date,
		Owner:      secretInput.Owner,
	}
	secretJSONasBytes, err := json.Marshal(secret)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Save secret to state
	err = stub.PutPrivateData("collectionSecret", secretInput.Name, secretJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Create secret private details (secret original) object with original data, marshal to JSON, and save to state
	secretOriginal := &secretOriginal{
		ObjectType: "secretOriginal",
		Name:       secretInput.Name,
		Original:   secretInput.Original,
	}
	secretOriginalBytes, err := json.Marshal(secretOriginal)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionSecretDetails", secretInput.Name, secretOriginalBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "owner~Hash"
	hashNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{secret.Owner, secret.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	stub.PutPrivateData("collectionSecret", hashNameIndexKey, value)

	fmt.Println("- end init secret")
	return shim.Success(nil)
}

func (t *Chaincode) readSecret(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the secret to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionSecret", name) //get the secret from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Secret does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *Chaincode) readOriginal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the secret to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionSecretDetails", name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Secret private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}
