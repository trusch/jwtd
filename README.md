jwtd - a JSON Web Token Daemon
==============================

jwtd is a service responsible for authenticating users

It provides cryptographically signed JSON Web Tokens based on a role based access control system. These tokens can then be validated using jwtd-proxy, providing the possibility to secure arbitary http based services.

## Concept

jwtd is a HTTP server. It provides an endpoint which implements token creation. You simply send your username + password and a specification of what you want to access and jwtd sends a token back to you if your user is allowed to do this.

To declare users, roles and access rights jwtd has a companion program called jwtd-ctl. This is a commandline tool which allows you to:

* set username + password combinations
* associate users with groups
* declare access rights for groups

You can also use a restfull API provided by jwtd.

Access rights are tuples of a service id and a list of labels. The service is typically the URI of the some service you (as a system architect) want to protect. The labels are string-string tuples and used to implement fine grained control. It's up to the usecase which labels exists.

## Features
* JWTD Core
  * User Management
  * Group Management
  * Multi Project Support
  * Commandline Interface
  * REST API Interface
  * Wildcard Rights
* JWTD Proxy
  * Token validation for multiple services
  * TLS Termination
  * Pattern based Requierements
  * Variables in Route->Requirement mapping
  * File based Configuration

## Example
The `example` directory contains a docker based example setup. It only assumes that the hostnames "jwtd" and "http-echo" are locally resolvable to 127.0.0.1. You can achieve that by manipulation your /etc/hosts for example. You can start the example via `bash example.sh`.

A more complex setup is shown in the `test` directory. It contains the integration test setup for this project. Therefore it starts by running an example setup with docker and then fires HTTP requests against jwtd-proxy for testing purpuses. I think about 90% of the functionality of jwtd is covered here. The tests are written as golang tests and can therefore be run via `go test -v`.
