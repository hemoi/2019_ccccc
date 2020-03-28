BLACKCHAIN IDEATON 2019
=======================

## 기술임치제도 2.0

기술임치제도란 중소기업의 기술자료를 신뢰성 있는 전문기관에 보관해 중소기업의 기술유출을 방지하는 것으로, 대기업 입장에서는 중소기업이 폐업 또는 파산하더라도 기술사용을 보장받을 수 있으며, 중소기업 입장에서는 부득이하게 기술탈취가 발생했을 때 특정시기부터 기술을 보유하고 있었다는 증거로 기술탈취시 해당 제도를 유용하게 사용할 수 있다. 


#### 기존 제도의 문제점을 해결하기 위한 탈중앙화 된 기술 임치제도 2.0을 설계 했다.

hyperledger의 privateDB를 활용, 권한이 있는 사람만이 원본데이터를 열람하고 권한이 없을 경우에는 hash값만 열람할 수 있게 함

------
## How to test

*ref. hyperledger fabric privateDB*  
<https://hyperledger-fabric.readthedocs.io/en/release-1.4/private-data/private-data.html>

하이퍼레저 패브릭 first-network 환경설정이 되었다고 가정  
해당 폴더를 first-network fabric-samples/chaincode에 넣는다.

### first network up을 위한 기존 network down
```
cd fabric-samples/first-network
./byfn.sh down

docker rm -f $(docker ps -a | awk '($2 ~ /dev-peer.*.docCC.*/) {print $1}')
docker rmi -f $(docker images | awk '($1 ~ /dev-peer.*.docCC.*/) {print $3}')

./byfn.sh up -c mychannel -s couchdb
```

```
docker exec -it cli bash
```

```
peer chaincode install -n docCC -v 1.0 -p github.com/chaincode/manageSecret/go/

export CORE_PEER_ADDRESS=peer1.org1.example.com:8051
peer chaincode install -n docCC -v 1.0 -p github.com/chaincode/manageSecret/go/

export CORE_PEER_LOCALMSPID=Org2MSP
export PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

export ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile $ORDERER_CA -C mychannel -n docCC -v 1.0 -c '{"Args":["init"]}' -P "OR('Org1MSP.member','Org2MSP.member')" --collections-config  $GOPATH/src/github.com/chaincode/manageSecret/collections_config.json


export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

### initSecret
```
export Secret=$(echo -n "{\"name\":\"Secret@1\",\"owner\":\"choi\",\"orignal\":\"this is the original Files\"}" | base64 | tr -d \\n)
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n docCC -c '{"Args":["initSecret"]}'  --transient "{\"secret\":\"$Secret\"}"
```
name : 해당 정보 이름\
owner : 해당 정보의 소유자\
original : 원본 데이터 값

### readSecret
```
peer chaincode query -C mychannel -n docCC -c '{"Args":["readSecret","Secret@1"]}'
```


### readOriginal
```
peer chaincode query -C mychannel -n docCC -c '{"Args":["readOriginalDetails","Secret@1"]}'
```
권한이 있기 때문에 문제가 없다.

----

*권한이 없을경우* 

 Org2로 변경
```
export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
export CORE_PEER_LOCALMSPID=Org2MSP
export PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
```
```
peer chaincode query -C mychannel -n docCC -c '{"Args":["readSecret","Secret@1"]}'
```
문제가 없는 것 확인
```
peer chaincode query -C mychannel -n docCC -c '{"Args":["readOriginalDetails","Secret@1"]}'
```
권한이 없기 때문에 Error

-----
#### 새로운 터미널에서
```
docker logs peer0.org1.example.com 2>&1 | grep -i -a -E 'private|pvt|privdata'
```
현재 블록 상황을 볼 수 있다.