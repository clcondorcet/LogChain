/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"errors"
	"log"
	"strings"

	"crhack.com/logchain/chaincode/chaincode"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var actions = map[string]map[string]bool{
	"Org1MSP": {"ReadAsset": true, "AssetExists": true, "GetAllAssets": true, "InitLedger": true, "DeleteAsset": true, "AddAsset": true},
	"Org2MSP": {"ReadAsset": true, "AssetExists": true, "GetAllAssets": true, "InitLedger": true, "DeleteAsset": true, "AddAsset": true},
}

func CheckACL(ctx contractapi.TransactionContextInterface) error {
	// Read incoming data from stub
	stub := ctx.GetStub()

	// Extract operation name
	operation, _ := stub.GetFunctionAndParameters()
	operationSplitted := strings.Split(operation, ":")
	operationName := operationSplitted[len(operationSplitted)-1]

	// Get requestor info from stub
	mspID, _ := cid.GetMSPID(stub)
	//userID, _ := cid.GetID(stub)
	//value, found, _ := cid.GetAttributeValue(stub, "role")

	user_actions := actions[mspID]

	if !user_actions[operationName] {
		return errors.New("You are not allowed to do this operation.")
	}

	// Operation allowed
	return nil
}

func main() {
	logChainContract := new(chaincode.SmartContract)
	logChainContract.BeforeTransaction = CheckACL
	assetChaincode, err := contractapi.NewChaincode(logChainContract)
	if err != nil {
		log.Panicf("Error creating logchain chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting logchain chaincode: %v", err)
	}
}
