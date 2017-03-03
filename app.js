var hfc = require('hfc');
var fs = require('fs');
var util = require('util');
var tutil = require('./test-util.js');

setup();

function setup() {
  chain = hfc.newChain("rewardsChain");

  process.on('exit', function (){
    chain.eventHubDisconnect();
  });
}

function enrollAndRegisterUsers() {
  chain.enroll("jim", "6avZQLwcUe9b", function(err, user) {
      if (err) throw Error("\nERROR: failed to enroll user : " + err);

      console.log("\nEnrolled user successfully");

      // Set this user as the chain's registrar which is authorized to register other users.
      chain.setRegistrar(admin);

      //creating a new user
      var registrationRequest = {
          enrollmentID: newUserName,
          affiliation: config.user.affiliation
      };
      chain.registerAndEnroll(registrationRequest, function(err, user) {
          if (err) throw Error(" Failed to register and enroll " + newUserName + ": " + err);

          console.log("\nEnrolled and registered " + newUserName + " successfully");
          userObj = user;
          //setting timers for fabric waits
          chain.setDeployWaitTime(config.deployWaitTime);
          console.log("\nDeploying chaincode ...");
          deployChaincode();
      });
  });
}

function deploy() {

}

function assign() {

}

function redeem() {

 }

function query() {

}
