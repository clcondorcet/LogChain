# LogChain

This project is a proof of concept.

The goal was to test if a "LogChain" (a way to store logs in a blockchain) is possible and viable.
We used HyperLedger Fabric to support the blockchain.

The SmartContract ("ChainCode") is coded under the `./ChainCode` folder.

This project is not intended for production, the informations below help to install and test the hyperledger network locally.

# Inital Setup

**You must have docker and docker-compose V2 !!**

```bash
make init # To initialize everything.
or
make init-git # To initialize the git submodules and some parameters.
make init-project # To install fabric binaries and docker images.
```

You can also install `JQ` to make sure that you can run the test network.

```bash
sudo apt-get install jq # Debian and Ubuntu, go to https://jqlang.github.io/jq/download/ for others
```

# Running the test network.

## Setup the environement variables

Take a look at the `env.vars` file.

It can look like this:
```
GOPATH=/mnt/c/Users/cleme/go
CHANNEL_NAME=logchannel
CHAINCODE_PATH=../../ChainCode
COLLECTION_CONFIG_PATH=../../collection_config.json
```

All relative paths must be considered from the `fabric-samples/test-network` folder.

If you have installed golang, to find the binary folder run the following:

```bash
echo "$PATH" | grep "go/bin"
```

You should now put the absolute path in the env vars.

If you don't have go installed you can probably install it via apt with the package name `golang` or follow the link: https://go.dev/dl/ .

Warning the channel name cannot contains uppercase character.

## Run the test network

Make sure to have set the proper environement variables.

```bash
make run-test-network
```

## Stop the test network

```bash
make stop-test-network
```

wait for the mention of `Network Ready !!`.

# Using the test network

Before trying to use the peer command you shoud run this (in the root folder):

```bash
. ./env.vars ; \
export PATH=${PWD}/fabric-samples/bin:$PATH ; \
export FABRIC_CFG_PATH=${PWD}/fabric-samples/config/ ; \
export $(./setEnv.sh Org1 | xargs)
```

To switch to Org2 peer (or swith back to Org1) use this command

```bash
export $(./setEnv.sh Org2 | xargs)
```

When using the peer command,
Always use those parameters:
- `-C $CHANNEL_NAME`
- `-n logcontract`

## All functions 

### Query

Base command :
```bash
peer chaincode query -C $CHANNEL_NAME -n logcontract
```

> Example : `-c '{"Args":["GetAllAssets"]}'`

> Example : `-c '{"Args":["GetAssetByRange","1706704300","1706704522"]}'`

> Example : `-c '{"Args":["ReadAsset","87256cd1e75c4a60dc6f569742c6302a3bfaf011"]}'`

> Example : `-c '{"Args":["AssetExists","87256cd1e75c4a60dc6f569742c6302a3bfaf011"]}'`

### Invoke

Base command :
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --waitForEvent -C $CHANNEL_NAME -n logcontract --peerAddresses $CORE_PEER0_ORG1_ADDRESS --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses $CORE_PEER0_ORG2_ADDRESS --tlsRootCertFiles $PEER0_ORG2_CA
```

> Example : `-c '{"Args":["AddAsset","test.com","This is an important log !!","1706704322"]}'`

> Example : `-c '{"Args":["AddAssets","[{\"hostname\":\"test.com\",\"message\":\"This is an important log !!\",\"timestamp\":\"1706704322\"},{\"hostname\":\"test.com\",\"message\":\"This is another important log !!\",\"timestamp\":\"1706704325\"}]"]}'`

> Example : `-c '{"Args":["DeleteAsset","87256cd1e75c4a60dc6f569742c6302a3bfaf011"]}'`

## LogChainTools

For the moment, no automation is available to setup the configuration to every environment.
To use the tools you must change every paths in the `./LogChainTool/network-config.yaml`.

The linux binaries are already built but you can build your own versions with go build command.

### HTTP-Server

The HTTP-Server is a http server that listen to port 5000 locally.
There are two routes available:
- `/querry` to use the querry features on the network.
- `/invoke` to use the invoke features on the network.

### Parser

This tool allows to tail a file of your choice and insert new lines directly in the LogChain.
You can run it like this:
```shell
./LogChainParser -file <path> -hostname <Your hostname> [-readall]
```

The `-readall` flag will allow the parser to read the whole file and insert each lines one by one in the LogChain (not recommended, use the tests python scripts to have better performances).

## Tests

The tests python scripts allow you to test the LogChain.
The scripts uses the HTTP-Server Tool to work.

Available scripts:
- `GetAllLogs.py` Read all the logs in the LogChain and output them in the `output.log` file.
- `InsertLogFileTest.py` Insert file line by line in the LogChain (you may need to modify the timestamp reading function to your log format), the insertion is divided into workers to have better performances.
- `InsertLogFileTestV2.py` Same as the previous script but uses another function of the SmartContract that can accept multiples lines. You can define the number of workers and lines to be inserted (better performances than v1 and allow more transactions per seconds in the LogChain).
- `RemoveAllLogs.py` Remove all the logs in the LogChain. This is really innefficient, recreating the network is less time consuming when having a lot of logs in the LogChain.
