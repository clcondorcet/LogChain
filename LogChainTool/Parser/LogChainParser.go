package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	org1CfgPath = "../network-config.yaml"
	//org1CfgPath = "../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json"
	Org1Cert = "../../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
	Org1Key  = "../../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
)

var contract *gateway.Contract

func main() {
	filePath := flag.String("file", "", "Chemin du fichier de journal")
	hostname := flag.String("hostname", "", "Hostname de l'ordinateur")
	readAll := flag.Bool("readall", false, "Insérer tout le fichier.")
	flag.Parse()
	if *filePath == "" {
		fmt.Println("Veuillez spécifier le chemin du fichier de journal en utilisant l'option -file")
		return
	}
	if *hostname == "" {
		fmt.Println("Veuillez spécifier le hostname en utilisant l'option -hostname")
		return
	}

	// Setup Hyperledger client
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

	fmt.Println("Hyperledger linked! Now starting parser ...")

	// Setup file reader
	lines := make(chan string)
	go tailFile(*filePath, lines, *readAll)

	for line := range lines {
		currentTime := time.Now()
		timestampSeconds := currentTime.Unix()
		_, err := contract.SubmitTransaction("AddAsset", *hostname, line, strconv.FormatInt(timestampSeconds, 10))
		if err != nil {
			fmt.Printf("Failed to submit transaction: %s\n", err)
		}
	}
}

func tailFile(filePath string, lines chan string, readAll bool) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if readAll {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
	} else {
		_, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			fmt.Println("Error seeking to end of file:", err)
			return
		}
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		fmt.Println("Error adding file to watcher:", err)
		return
	}

	scanner := bufio.NewScanner(file)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Lorsque le fichier est modifié, lisez les nouvelles lignes
				for scanner.Scan() {
					lines <- scanner.Text()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error watching file:", err)
		}
	}
}
