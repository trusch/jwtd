jwtd - a JSON Web Token Daemon
==============================

jwtd is a service responsible for authenticating users

It provides cryptographically signed JSON Web Tokens based on a role based access control system. These tokens can then be validated on arbitary services to get rid of implementing user authentication over and over again.

## Concept

jwtd is a HTTP server. It only provides one default route which implements token creation. You simply send your username + password and a specification of what you want to access and jwtd sends a token back to you if your user is allowed to do this.

To declare users, roles and access rights jwtd has a companion program called jwtd-ctl. This is a commandline tool which allows you to:

* set username + password combinations
* associate users with groups
* declare access rights for groups

Access rights are tuples of service and subject. The service is typically the URI of the some service you (as a system architect) want to protect. The subject is used to implement fine grained control. In fact, the subject is just a string which is passed into the token. It's up to the usecase which subjects exists.

## Example

Lets image four entities / entity groups:

* A fleet of sensor devices
* A webserver serving a frontend
* A datastore for sensordata
* A set of human system users

Some auth related requierments:

* The sensor devices should be able to write to the datastore, but should not read from it.
* The webserver should only serve its content to authenticated users.
* Authenticated users should be able to read from the datastore, but are not allowed to write to it.

To accomplish this, we define two roles with the following access rights:

* Human
  * Service: webserver, Subject: read
  * Service: datastore, Subject: read
* Sensor
  * Service: datastore, Subject: write

Next we create some users (we will only create one shared user for all sensor devices)

* Name: temp-sensor, Role: Sensor
* Name: alice, Role: Human
* Name: bob, Role: Human
* Name: carl, Role: Human

### Sensor write

If a sensor wants to send data to the datastore, it must fetch a token from jwtd. Therefore it requests a token for the service "datastore" and the subject "write" using its username and password.
jwtd lookups the groups of the given user and asserts that the requested token is authorized by checking the accessrights for the given group. Because our "temp-sensor" is in the "Sensor" group, access is granted and a token containing the requested service and subject gets created, signed and send to our sensor. These token can now be embedded into the request for the datastore.
The datastore extracts the token and verifies the signature with jwtd's public key. If its valid it checks if the token contains its own service identifier ("datastore") and the subject "write". If all this is fullfilled, the requests payload is written to the datastore.

### User read

If a alice wants to read data, she needs access to two services: The webserver and the datastore. So alice needs to request a token which grants her read-access to both services. So she sends her username and passwords with the two claims (webserver:read, datastore:read) to jwtd. Because alice has the role "Human", access is granted. The returned token contains both claims: webserver:read and datastore:read.

@TODO: hier weitermachen!
@TODO: tagesmutter anrufen!

(The datastore validates the token and recognizes the "read" subject and therefore presents data to alice.)
