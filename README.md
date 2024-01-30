# LogChain
 Blockchain for logs. Engineering project

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
- `-n logContract`

## All functions 

### Querry

Base command :
```bash
peer chaincode query -C $CHANNEL_NAME -n logContract
```

> Example : `-c '{"Args":["GetAllAssets"]}'`

> Example : `-c '{"Args":["ReadAsset","LOG1"]}'`

> Example : `-c '{"Args":["AssetExists","LOG1"]}'`

### Invoke

Base command :
```bash
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA --waitForEvent -C $CHANNEL_NAME -n logContract --peerAddresses $CORE_PEER0_ORG1_ADDRESS --tlsRootCertFiles $PEER0_ORG1_CA --peerAddresses $CORE_PEER0_ORG2_ADDRESS --tlsRootCertFiles $PEER0_ORG2_CA
```

> Example : `-c '{"Args":["AddAsset","LOG1","test.com","This is an important log !!","109877891"]}'`

> Example : `-c '{"Args":["DeleteAsset","LOG1"]}'`
