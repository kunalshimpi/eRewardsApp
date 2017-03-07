var hfc = require('hfc');
var fs = require('fs');
var util = require('util');
var tutil = require('./test-util.js');

setup();

function setup() {
  chain = tutil.getTestChain("rewardsChain");

  process.on('exit', function (){
    chain.eventHubDisconnect();
  });
}

var chain,admin,webAdmin,webUser;

function enrollAdmin() {
  console.log("Enrolling Admin");

  chain.enroll("admin", "Xurw3yU9zI0l", function (err, user) {
     if (err) return err;
     admin = user;
     chain.setRegistrar(admin);
  });
}

function enrollWebAdmin() {
  console.log("Enrolling Web Admin");
  registerAndEnroll("webAdmin", "client", null, {roles:['client']}, chain, function(err,user) {
     if (err) return err;
     webAdmin = user;
     chain.setRegistrar(webAdmin);
  });
}

function enrollWebUser() {
  console.log("Enrolling Web User");
  registerThenEnroll("webUser", "client", null, webAdmin, chain, function(err, user) {
     if (err) return err;
     webUser = user;
  });
}

function registerAndEnroll(name, role, attributes, registrar, chain, cb) {
    console.log("registerAndEnroll %s",name);
    // User is not enrolled yet, so perform both registration and enrollment
    var registrationRequest = {
         roles: [ role ],
         enrollmentID: name,
         affiliation: "bank_a",
         attributes: attributes,
         registrar: registrar
    };
    chain.registerAndEnroll(registrationRequest,cb);
}

function registerThenEnroll(name, role, attributes, registrar, chain, cb) {
    console.log("registerThenEnroll %s",name);
    // User is not enrolled yet, so perform both registration and enrollment
    var registrationRequest = {
         roles: [ role ],
         enrollmentID: name,
         affiliation: "bank_a",
         attributes: attributes,
         registrar: registrar
    };
    // Test chain.register()
    chain.register(registrationRequest, function(err, enrollmentPassword) {
        if (err) {
            console.log("registerThenEnroll: couldn't register name ", err)
            return cb(err);
        }
        // Fetch name's member so we can set the Registrar
        chain.getMember(registrar, function(err, member) {
            if (err) {
                console.log("could not get member for ", name, err);
                return cb(err);
            }
            //console.log("I did find this member", member)
            chain.setRegistrar(member);
            });

        // Test chain.enroll using password returned by chain.register()
        chain.enroll(name, enrollmentPassword, function(err, member) {
        //console.log("am I defined?", cb)  /* yes, defined */
            if (err) {
                console.log("registerThenEnroll: enroll failed", err);
                return cb(err);
            }
            //console.log("registerThenEnroll: enroll succeeded for registration request =", registrationRequest);
            return cb(err, member);
        });
   });   /* end chain.register */
   chain.setDeployWaitTime(200);
   console.log("\nDeploying chaincode ...");
   deploy();
}

function deploy() {
  var args = null;

  var deployRequest = {
      fcn: init,
      args: args,
      chaincodePath: "chaincode",
      certificatePath: "/certs/peer/cert.pem"
  };

  var deployTx = webUser.deploy(deployRequest);

  deployTx.on('complete', function(results) {

      chaincodeID = results.chaincodeID;
      console.log("\nChaincode ID : " + chaincodeID);
      console.log(util.format("\nSuccessfully deployed chaincode: request=%j, response=%j", deployRequest, results));

      fs.writeFileSync(chaincodeIDPath, chaincodeID);
      assign();
  });

  deployTx.on('error', function(err) {
      console.log(util.format("\nFailed to deploy chaincode: request=%j, error=%j", deployRequest, err));
      process.exit(1);
  });
}

function assign() {
  var args = ["1000","kunal","moto360"];
  var eh = chain.getEventHub();
  // Construct the invoke request
  var invokeRequest = {
      // Name (hash) required for invoke
      chaincodeID: chaincodeID,
      // Function to trigger
      fcn: "assign",
      // Parameters for the invoke function
      args: args
  };

  // Trigger the invoke transaction
  var invokeTx = webUser.invoke(invokeRequest);

  // Print the invoke results
  invokeTx.on('submitted', function(results) {
      // Invoke transaction submitted successfully
      console.log(util.format("\nSuccessfully submitted chaincode invoke transaction: request=%j, response=%j", invokeRequest, results));
  });
  invokeTx.on('complete', function(results) {
      // Invoke transaction completed successfully
      console.log(util.format("\nSuccessfully completed chaincode invoke transaction: request=%j, response=%j", invokeRequest, results));
      redeem();
  });
  invokeTx.on('error', function(err) {
      // Invoke transaction submission failed
      console.log(util.format("\nFailed to submit chaincode invoke transaction: request=%j, error=%j", invokeRequest, err));
      process.exit(1);
  });

  //Listen to custom events
  var regid = eh.registerChaincodeEvent(chaincodeID, "evtsender", function(event) {
      console.log(util.format("Custom event received, payload: %j\n", event.payload.toString()));
      eh.unregisterChaincodeEvent(regid);
  });
}

function redeem() {
  var args = ["400","kunal","apolloHospital"];
  var eh = chain.getEventHub();
  // Construct the invoke request
  var invokeRequest = {
      // Name (hash) required for invoke
      chaincodeID: chaincodeID,
      // Function to trigger
      fcn: "redeem",
      // Parameters for the invoke function
      args: args
  };

  // Trigger the invoke transaction
  var invokeTx = webUser.invoke(invokeRequest);

  // Print the invoke results
  invokeTx.on('submitted', function(results) {
      // Invoke transaction submitted successfully
      console.log(util.format("\nSuccessfully submitted chaincode invoke transaction: request=%j, response=%j", invokeRequest, results));
  });
  invokeTx.on('complete', function(results) {
      // Invoke transaction completed successfully
      console.log(util.format("\nSuccessfully completed chaincode invoke transaction: request=%j, response=%j", invokeRequest, results));
      query();
  });
  invokeTx.on('error', function(err) {
      // Invoke transaction submission failed
      console.log(util.format("\nFailed to submit chaincode invoke transaction: request=%j, error=%j", invokeRequest, err));
      process.exit(1);
  });

  //Listen to custom events
  var regid = eh.registerChaincodeEvent(chaincodeID, "evtsender", function(event) {
      console.log(util.format("Custom event received, payload: %j\n", event.payload.toString()));
      eh.unregisterChaincodeEvent(regid);
  });
 }

function query() {
  var args = ["kunal"];
  // Construct the query request
  var queryRequest = {
      // Name (hash) required for query
      chaincodeID: chaincodeID,
      // Function to trigger
      fcn: "read",
      // Existing state variable to retrieve
      args: args
  };

  // Trigger the query transaction
  var queryTx = webUser.query(queryRequest);

  // Print the query results
  queryTx.on('complete', function(results) {
      // Query completed successfully
      console.log("\nSuccessfully queried  chaincode function: request=%j, value=%s", queryRequest, results.result.toString());
      process.exit(0);
  });
  queryTx.on('error', function(err) {
      // Query failed
      console.log("\nFailed to query chaincode, function: request=%j, error=%j", queryRequest, err);
      process.exit(1);
  });
}
