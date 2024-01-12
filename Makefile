init: init-git init-project

init-git:
	git config --global core.autocrlf false
	git config --global core.longpaths true
	git submodule init
	cd fabric-samples && git config --global core.autocrlf false && git config --global core.longpaths true

init-project:
	./install-fabric.sh

run-test-network:
	{ \
		set -e ; \
		. ./env.vars ; \
		export PATH=$${PWD}/fabric-samples/bin:$$PATH ; \
		cd ./fabric-samples/test-network ; \
		echo "Initialisation du network ..." ; \
		./network.sh up ; \
		echo "Création du Channel ..." ; \
		./network.sh createChannel -c $$CHANNEL_NAME ; \
		echo "Déploiement du smartcontract. Cette opération prend du temps ..." ; \
		./network.sh deployCC -ccn logContract -ccp $$CHAINCODE_PATH -ccl go -c $$CHANNEL_NAME ; \
		echo "Initialisation de la chaine ..." ; \
		export FABRIC_CFG_PATH=$${PWD}/../config/ ; \
		export CORE_PEER_TLS_ROOTCERT_FILE=$${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt ; \
		export CORE_PEER_MSPCONFIGPATH=$${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ; \
		peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C $$CHANNEL_NAME -n logContract --peerAddresses localhost:7051 --tlsRootCertFiles "$${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}' ; \
		echo "Network Ready !!" ; \
	}

stop-test-network:
	{ \
		set -e ; \
		. ./env.vars ; \
		cd ./fabric-samples/test-network ; \
		./network.sh down ; \
	}
