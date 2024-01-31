package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	org1CfgPath = "./network-config.yaml"
	//org1CfgPath = "../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json"
	Org1Cert = "../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
	Org1Key  = "../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
)

var contract *gateway.Contract

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "TRUE")

	certBytes, err := os.ReadFile(Org1Cert)
	if err != nil {
		fmt.Printf("unable to read file: %v", err)
		os.Exit(1)
	}

	keyBytes, err := os.ReadFile(Org1Key)
	if err != nil {
		fmt.Printf("unable to read file: %v", err)
		os.Exit(1)
	}

	wallet := gateway.NewInMemoryWallet()
	wallet.Put("Admin", gateway.NewX509Identity("Org1MSP", string(certBytes), string(keyBytes)))

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(org1CfgPath))),
		gateway.WithIdentity(wallet, "Admin"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to gateway: %s\n", err)
		os.Exit(1)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("logchannel")
	if err != nil {
		fmt.Printf("Failed to get network: %s\n", err)
		os.Exit(1)
	}

	contract = network.GetContract("logcontract")

	fmt.Println("Hyperledger linked! Now starting http server ...")

	mux := http.NewServeMux()
	mux.HandleFunc("/querry", querry)
	mux.HandleFunc("/invoke", invoke)

	err = http.ListenAndServe(":3333", mux)
}

// HTTP routes

type data struct {
	Function string   `json:"function"`
	Args     []string `json:"args"`
}

func querry(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d data = data{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := contract.EvaluateTransaction(d.Function, d.Args...)
	if err != nil {
		fmt.Printf("Failed to evaluate transaction: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Println("querry called and run !!")
	//fmt.Println(string(result))
	w.Write(result)
}

func invoke(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d data = data{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := contract.SubmitTransaction(d.Function, d.Args...)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//fmt.Println("invoke called and run !!")
	//fmt.Println(string(result))
	w.Write(result)
}
