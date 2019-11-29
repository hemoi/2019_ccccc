# How to start

## Enter the docker

```bash
docker exec -it cli bash
```

# Run back-end (node.js)
in docker.
<!--
### At the below repo:
```bash
> pwd
/opt/gopath/src/github.com/hyperledger/fabric/peer
```
### Run
-->
```bash
node server.js
```

# Run front-end (react.js)
in docker

## test
```bash
cd ~
cd /
cd test/test
npm start
```

## test2
```bash
cd ~
cd /
cd test/test2
npm start
```

tmp. mapping:
* main --> readSecret
* viewer --> readOriginalDetails
* tx --> readOriginal
