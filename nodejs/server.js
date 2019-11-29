var exec = require('child_process').exec, child;

/*
function _query(func, name) {
    let command = "peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\""+ func +"\",\"" + name + "\"]}'";
    child = exec(command,  function (error, stdout, stderr) {

    // child = exec("peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\"readSecret\",\"Secret@1\"]}'", function (error, stdout, stderr) {
        console.log('stdout: ' + stdout);

        console.log('stderr: ' + stderr);
        if (error !== null) {
            console.log('exec error: ' + error);
        }
    });
}
*/

/*
function _export(env, val) {
    let command = "export " + env + "=" + val;
    child = exec(command, function (error, stdout, stderr) {
        console.log('stdout: ' + stdout);
        console.log('stderr: ' + stderr);
        if (error !== null){
            console.log('exec error: ' + error);
        }
    });
}
*/

// main

// query
// _query("readSecret", "Secret@1");

// export
// _export("TEST", "1234");

const express = require("express");
const bodyParser = require("body-parser");
const cors = require("cors");

const http_port = 30301;

function initHttpServer() {
    const app = express();
    app.use(cors());
    app.use(bodyParser.json());

    app.get("/readSecret", function (req, res) {
        // const args = req.body.args;
        const args = "Secret@1";
        console.log("args: ", args);

        let command = "peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\""+ "readSecret" +"\",\"" + args + "\"]}'";
        child = exec(command,  function (error, stdout, stderr) {

    // child = exec("peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\"readSecret\",\"Secret@1\"]}'", function (error, stdout, stderr) {
            console.log('stdout: ' + stdout);

            console.log('stderr: ' + stderr);
            if (error !== null) {
                console.log('exec error: ' + error);
            }
    
            // res.send(JSON.stringify([JSON.parse(stdout)], null, 2));
            res.send([JSON.parse(stdout)]);
        });
    });

    app.get("/readOriginalDetails", function (req, res) {
        const args = "Secret@1";
        console.log("args: ", args);

        let command = "peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\""+ "readOriginalDetails" +"\",\"" + args + "\"]}'";
        child = exec(command,  function (error, stdout, stderr) {

    // child = exec("peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\"readSecret\",\"Secret@1\"]}'", function (error, stdout, stderr) {
            console.log('stdout: ' + stdout);

            console.log('stderr: ' + stderr);
            if (error !== null) {
                console.log('exec error: ' + error);
            }

            res.send(stdout);
        });
    });

    app.get("/readOriginal", function (req, res) {
        const args = "Secret@1";
        console.log("args: ", args);

        let command = "peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\""+ "readOriginal" +"\",\"" + args + "\"]}'";
        child = exec(command,  function (error, stdout, stderr) {

    // child = exec("peer chaincode query -C mychannel -n marblesp -c '{\"Args\":[\"readSecret\",\"Secret@1\"]}'", function (error, stdout, stderr) {
            console.log('stdout: ' + stdout);

            console.log('stderr: ' + stderr);
            if (error !== null) {
                console.log('exec error: ' + error);
            }
            
            res.send(stdout);
        });
    });

    /*
    app.get("/blocks", function (req, res) {
        res.send(bc.getBlockchain());
    });
    app.get('/block/:number', function (req, res) {
        const targetBlock = bc.getBlockchain().find(function (block) {
            return block.header.index == req.params.number;
        });
        res.send(targetBlock);
    });
    app.post("/mineBlock", function (req, res) {
        const data = req.body.data || [];
        const newBlock = bc.mineBlock(data);
        if (newBlock === null) {
            res.status(400).send('Bad Request');
        }
        else {
            res.send(newBlock);
        }
    });
    */

    app.listen(http_port, function () { console.log("Listening http port on: " + http_port) });
}

// main
initHttpServer();
