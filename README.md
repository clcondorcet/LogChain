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

wait for the mention of `Network Ready !!`.

# Using the test network

Before trying to use the peer command you shoud run this (in the root folder):

```bash
. ./env.vars ; \
export PATH=${PWD}/fabric-samples/bin:$PATH ; \
export FABRIC_CFG_PATH=${PWD}/fabric-samples/config/ ; \
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt ; \
export CORE_PEER_MSPCONFIGPATH=${PWD}/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ; \
```

When using the peer command,
Always use those parameters:
- `-C $CHANNEL_NAME`
- `-n logContract`

## Querry All

```bash
peer chaincode query -C $CHANNEL_NAME -n logContract -c '{"Args":["GetAllAssets"]}'
```
