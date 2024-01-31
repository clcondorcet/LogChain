package chaincode

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ID        string `json:"ID"`
	Hostname  string `json:"hostname"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

const collection_name = "logs"

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// assets := []Asset{
	// 	{Hostname: "test.test", Message: "This is the first log", Timestamp: "1098737410837"},
	// 	{Hostname: "test.test", Message: "This is the second log", Timestamp: "109873749876"},
	// }

	// for _, asset := range assets {
	// 	assetJSON, err := json.Marshal(asset)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	err = ctx.GetStub().PutState(asset.ID, assetJSON)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to put to world state. %v", err)
	// 	}
	// }

	return nil
}

// func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
// 	assetJSON, err := ctx.GetStub().GetState(id)

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read from world state: %v", err)
// 	}
// 	if assetJSON == nil {
// 		return nil, fmt.Errorf("the asset %s does not exist", id)
// 	}

// 	var asset Asset
// 	err = json.Unmarshal(assetJSON, &asset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &asset, nil
// }

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetPrivateData(collection_name, id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
// 	assetJSON, err := ctx.GetStub().GetState(id)

// 	if err != nil {
// 		return false, fmt.Errorf("failed to read from world state: %v", err)
// 	}

// 	return assetJSON != nil, nil
// }

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetPrivateData(collection_name, id)

	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
// func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all assets in the chaincode namespace.
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var assets []*Asset
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var asset Asset
// 		err = json.Unmarshal(queryResponse.Value, &asset)
// 		if err != nil {
// 			return nil, err
// 		}
// 		assets = append(assets, &asset)
// 	}

// 	return assets, nil
// }

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collection_name, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

/*
{
	"selector": {
		"$and": [
      {
        "timestamp": {"$gte": low}
      },
      {
        "timestamp": {"$lte": high}
      }
    ]
	},
	"sort": [{"timestamp": "asc"}],
	"use_index": ["_design/timestampIndexDoc","timestampIndex"]
}

{
	"Args":[
		"QueryAssets",
		"{
			\"selector\":{
				\"docType\":\"asset\",
				\"owner\":\"tom\"
			},
			\"use_index\":[
				\"_design/indexOwnerDoc\",
				\"indexOwner\"
			]
		}"
	]
}

{"selector":{"$and":[{"timestamp":{"$gte": low}},{"timestamp":{"$lte":high}}]},"sort":[{"timestamp": "asc"}],"use_index":["timestampIndexDoc","timestampIndex"]}
,\"fields\": [\"_id\",\"hostname\",\"message\",\"timestamp\"]
*/

func (s *SmartContract) GetAssetByRange(ctx contractapi.TransactionContextInterface, low string, high string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collection_name, "{\"selector\":{\"$and\":[{\"timestamp\":{\"$gte\":\""+low+"\"}},{\"timestamp\":{\"$lte\":\""+high+"\"}}]}}")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the asset %s does not exist", id)
// 	}

// 	return ctx.GetStub().DelState(id)
// }

func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().PurgePrivateData(collection_name, id)
}

// func (s *SmartContract) AddAsset(ctx contractapi.TransactionContextInterface, id string, hostname string, message string, timestamp string) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		return fmt.Errorf("the asset %s already exists", id)
// 	}

// 	asset := Asset{
// 		ID:        id,
// 		Hostname:  hostname,
// 		Message:   message,
// 		Timestamp: timestamp,
// 	}
// 	assetJSON, err := json.Marshal(asset)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(id, assetJSON)
// }

func (s *SmartContract) AddAsset(ctx contractapi.TransactionContextInterface, hostname string, message string, timestamp string) error {
	data_id := []byte(timestamp + hostname + message)
	sha := sha1.Sum(data_id)
	id := fmt.Sprintf("%x", sha[:])

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:        id,
		Hostname:  hostname,
		Message:   message,
		Timestamp: timestamp,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutPrivateData(collection_name, id, assetJSON)
}

type Input_log struct {
	Hostname  string `json:"hostname"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func (s *SmartContract) AddAssets(ctx contractapi.TransactionContextInterface, data_in string) error {
	var data []Input_log
	err := json.Unmarshal([]byte(strings.ReplaceAll(data_in, "'", "\"")), &data)
	if err != nil {
		return err
	}

	for _, d := range data {
		err = s.AddAsset(ctx, d.Hostname, d.Message, d.Timestamp)
		if err != nil {
			return err
		}
	}

	return nil
}
